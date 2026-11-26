"""Reconcile two CSV exports, fill a Word template and check its exception table."""

import csv
from collections import defaultdict
from copy import deepcopy
from decimal import Decimal
import argparse
from pathlib import Path
import subprocess
import sys
from xml.etree import ElementTree as ET
from zipfile import ZipFile

from router import check, evidence, record, review, write_json

W = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
NS = {"w": W}
ET.register_namespace("w", W)

SCHEMA = "scripts/office/schemas/ecma/fouth-edition/opc-relationships.xsd"
PROMPT = ("Reconcile the invoice and payment CSV exports, account for duplicate payment references, "
          "and fill the existing Word template with the exception rows. Validate the report and its package.")
STAGES = [
    {"id": "reconcile", "kind": "deliverable", "text": "Clean duplicate payment records and join CSV exports by invoice ID",
     "sources": [("data-analyst-engine", "Phase 2: Data Cleaning (80% of the work)", "Check for duplicates"),
                 ("data-analyst-engine", "Pandas Power Patterns", "Merge/Join")],
     "facets": [{"id": "deduplicate", "kind": "action", "anyOf": ["Check for duplicates"]}],
     "checks": [check("dedup", "Explain which repeated payment reference was excluded", "decision"),
                check("exceptions", "Save unmatched invoices and payments with exact currency deltas", "artifact")]},
    {"id": "report", "kind": "deliverable", "text": "Fill the supplied DOCX template with labelled exception rows",
     "dependsOn": ["reconcile"], "sources": [("docx", "DOCX creation, editing, and analysis", "edit `word/document.xml`")],
     "facets": [{"id": "word-xml", "kind": "format", "anyOf": ["word/document.xml"]}],
     "checks": [check("report-file", "Save the populated Word report", "artifact")]},
    {"id": "audit", "kind": "verification", "text": "Independently check the produced DOCX table cells and retained template parts",
     "dependsOn": ["report"], "sources": [("verification-before-completion", "The Gate Function", "check exit code"),
                                         ("docx", "CT_Relationship", 'name="Target"', SCHEMA)],
     "checks": [check("report-rows", "Read each exception row and compare ID, status and Decimal amount", "test"),
                check("template-parts", "Compare template parts including the unchanged styles", "test"),
                check("relationships", "Use the bundled OPC schema to check relationship attributes and targets", "test")]},
]

INVOICES = "invoice,amount\nA10,160.25\nA11,80.00\nA12,25.40\n"
PAYMENTS = "reference,invoice,amount\np1,A10,160.25\np2,A11,30.00\np2,A11,30.00\np3,A11,40.00\np4,A99,4.25\n"
EXPECTED = [["A11", "underpaid", "10.00"], ["A12", "unpaid", "25.40"], ["A99", "unmatched payment", "-4.25"]]
HEADERS = ["Invoice", "Exception", "Net outstanding (USD)"]
REL = "http://schemas.openxmlformats.org/package/2006/relationships"
OFFICE = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"


def reconcile(out):
    invoices = {}
    with (out / "invoices.csv").open(newline="") as f:
        for row in csv.DictReader(f):
            assert row["invoice"] not in invoices
            invoices[row["invoice"]] = Decimal(row["amount"])
    transfers, duplicates = {}, []
    paid = defaultdict(Decimal)
    with (out / "payments.csv").open(newline="") as f:
        for row in csv.DictReader(f):
            key = row["reference"]
            if key in transfers:
                assert transfers[key] == row, "Same reference with conflicting payment details"
                duplicates.append(key)
                continue
            transfers[key] = row
            paid[row["invoice"]] += Decimal(row["amount"])
    rows = []
    for key in sorted(invoices.keys() | paid.keys()):
        balance = invoices.get(key, Decimal(0)) - paid[key]
        if balance:
            status = ("unmatched payment" if key not in invoices else "unpaid" if not paid[key]
                      else "underpaid" if balance > 0 else "overpaid")
            rows.append([key, status, f"{balance:.2f}"])
    return duplicates, rows


def template(path):
    types = ('<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">'
             '<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>'
             '<Default Extension="xml" ContentType="application/xml"/>'
             '<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>'
             '<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/></Types>')
    with ZipFile(path, "w") as f:
        f.writestr("[Content_Types].xml", types)
        f.writestr("_rels/.rels", f'<Relationships xmlns="{REL}"><Relationship Id="rId1" Type="{OFFICE}/officeDocument" Target="word/document.xml"/></Relationships>')
        f.writestr("word/_rels/document.xml.rels", f'<Relationships xmlns="{REL}"><Relationship Id="rId1" Type="{OFFICE}/styles" Target="styles.xml"/></Relationships>')
        f.writestr("word/styles.xml", f'<w:styles xmlns:w="{W}"><w:style w:type="paragraph" w:styleId="Heading1"><w:name w:val="Heading 1"/></w:style></w:styles>')
        f.writestr("word/document.xml", f'<w:document xmlns:w="{W}"><w:body><w:p><w:pPr><w:pStyle w:val="Heading1"/></w:pPr><w:r><w:t>Payment exceptions</w:t></w:r></w:p><w:tbl/><w:sectPr/></w:body></w:document>')


def populate(source, dest, rows):
    with ZipFile(source) as original, ZipFile(dest, "w") as result:
        for name in original.namelist():
            content = original.read(name)
            if name == "word/document.xml":
                root = ET.fromstring(content)
                table = root.find(".//w:tbl", NS)
                grid = ET.SubElement(table, f"{{{W}}}tblGrid")
                for width in ("1700", "3200", "2800"):
                    ET.SubElement(grid, f"{{{W}}}gridCol", {f"{{{W}}}w": width})
                for values in [HEADERS, *rows]:
                    tr = ET.SubElement(table, f"{{{W}}}tr")
                    for value in values:
                        tc = ET.SubElement(tr, f"{{{W}}}tc")
                        p = ET.SubElement(tc, f"{{{W}}}p")
                        r = ET.SubElement(p, f"{{{W}}}r")
                        ET.SubElement(r, f"{{{W}}}t").text = value
                content = ET.tostring(root, encoding="utf-8", xml_declaration=True)
            result.writestr(name, content)


def check_rows(doc):
    with ZipFile(doc) as f:
        root = ET.fromstring(f.read("word/document.xml"))
    rows = [[cell.findtext("w:p/w:r/w:t", namespaces=NS) for cell in tr.findall("w:tc", NS)]
            for tr in root.findall(".//w:tbl/w:tr", NS)]
    assert rows == [HEADERS, *EXPECTED], f"Actual Word table rows: {rows}"
    total = sum(Decimal(row[2]) for row in rows[1:])
    assert total == Decimal("31.15")
    return f"{len(rows)-1} exception rows matched by ID/status/amount; net outstanding USD {total}."


def check_parts(source, doc):
    with ZipFile(source) as original, ZipFile(doc) as opened:
        assert opened.testzip() is None
        assert set(original.namelist()) == set(opened.namelist())
        for name in original.namelist():
            if name != "word/document.xml":
                assert opened.read(name) == original.read(name), name
        rel = ET.fromstring(opened.read("_rels/.rels")).find(f"{{{REL}}}Relationship")
        assert rel.attrib["Type"] == OFFICE + "/officeDocument" and rel.attrib["Target"] in opened.namelist()
        assert ET.fromstring(opened.read("word/document.xml")).find(".//w:pStyle", NS).attrib[f"{{{W}}}val"] == "Heading1"
    return "DOCX ZIP integrity passed; catalog, relationships and style parts match the supplied template."


def check_relationships(doc, schema):
    xsd = {"xsd": "http://www.w3.org/2001/XMLSchema"}
    root = ET.parse(schema).getroot()
    namespace = root.attrib["targetNamespace"]
    attrs = root.findall(".//xsd:complexType[@name='CT_Relationship']//xsd:attribute[@use='required']", xsd)
    required = {attr.attrib["name"] for attr in attrs}
    assert required == {"Id", "Target", "Type"}
    checked = []
    with ZipFile(doc) as opened:
        for part in ("_rels/.rels", "word/_rels/document.xml.rels"):
            rels = ET.fromstring(opened.read(part))
            assert rels.tag == f"{{{namespace}}}Relationships"
            for rel in rels:
                assert rel.tag == f"{{{namespace}}}Relationship" and required <= set(rel.attrib)
                target = Path(part).parent.parent / rel.attrib["Target"]
                assert target.as_posix() in opened.namelist()
                checked.append(target.as_posix())
    assert checked == ["word/document.xml", "word/styles.xml"]
    return "Read the loaded OPC schema. Both relationships contain Id, Type and Target and resolve to the preserved parts."


def verify_docx(out, label, operation, document, source=None):
    args = [sys.executable, "-B", str(Path(__file__).resolve()), operation, str(document)]
    if source is not None:
        args.append(str(source))
    run = subprocess.run(args, cwd=out, capture_output=True, text=True, timeout=20)
    write_json(out / f"{label}.json", {"args": args, "exitCode": run.returncode, "stdout": run.stdout, "stderr": run.stderr})
    return run


def run(router, out):
    out.mkdir()
    (out / "invoices.csv").write_text(INVOICES)
    (out / "payments.csv").write_text(PAYMENTS)
    source, dest = out / "template.docx", out / "reconciliation.docx"
    template(source)
    plan = router.plan(out, PROMPT, STAGES)
    inputs = router.load(out, plan, "reconcile")
    duplicate, rows = reconcile(out)
    assert duplicate == ["p2"] and rows == EXPECTED
    write_json(out / "reconciled.json", {"excludedDuplicateReferences": duplicate, "exceptions": rows})
    records = record(inputs, "applied", [
        evidence("dedup", "decision", "payments.csv + reconciled.json", f"Removed the second identical payment with reference {duplicate[0]}."),
        evidence("exceptions", "artifact", "reconciled.json", f"Retained {len(rows)} underpaid, missing and unmatched invoice/payment rows."),
    ])
    inputs = router.load(out, plan, "report")
    populate(source, dest, rows)
    records += record(inputs, "applied", [evidence("report-file", "artifact", dest.name, "Filled the existing exception table in the supplied Word template.")])
    audit = router.load(out, plan, "audit")
    checks_by_skill = {c["skill"]: c["id"] for c in next(s for s in plan["stages"] if s["requirement"]["id"] == "audit")["contributions"]}
    check_receipts = {c["contribution"]: c for c in audit}
    gate = [check_receipts[checks_by_skill["verification-before-completion"]]]
    relationship = [check_receipts[checks_by_skill["docx"]]]
    schema = out / "sources" / "docx" / SCHEMA
    shape = verify_docx(out, "relationships-run", "relationships", dest, schema)
    assert shape.returncode == 0, shape.stderr
    records += record(relationship, "verified", [
        evidence("relationships", "test", "relationships-run.json", shape.stdout.strip(), "passed")
    ])
    incorrect = deepcopy(rows)
    incorrect[0][2] = "0.01"
    corrupt = out / "incorrect-amount.docx"
    populate(source, corrupt, incorrect)
    rejected = verify_docx(out, "incorrect-amount-run", "rows", corrupt)
    assert rejected.returncode == 1 and "0.01" in rejected.stderr
    negative = record(gate, "applied", [evidence("report-rows", "test", "incorrect-amount-run.json", rejected.stderr.strip(), "failed")])
    review(router, out, plan, records + negative, "incorrect-amount", ["report-rows", "template-parts"])
    rows_run = verify_docx(out, "rows-run", "rows", dest)
    parts_run = verify_docx(out, "parts-run", "parts", dest, source)
    assert rows_run.returncode == parts_run.returncode == 0, (rows_run.stderr, parts_run.stderr)
    first = evidence("report-rows", "test", "rows-run.json", rows_run.stdout.strip(), "passed")
    second = evidence("template-parts", "test", "parts-run.json", parts_run.stdout.strip(), "passed")
    review(router, out, plan, records + record(gate, "verified", [first]), "missing-package-check", ["template-parts"])
    records += record(gate, "verified", [first, second])
    write_json(out / "records.json", records)
    review(router, out, plan, records, "complete")
    print("CSV → Word:", rows_run.stdout.strip(), parts_run.stdout.strip())


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Independently reopen and verify the produced Word file")
    parser.add_argument("operation", choices=("rows", "parts", "relationships"))
    parser.add_argument("document", type=Path)
    parser.add_argument("source", type=Path, nargs="?")
    args = parser.parse_args()
    if args.operation != "rows" and args.source is None:
        parser.error("parts requires the original template and relationships requires the OPC schema")
    try:
        if args.operation == "rows":
            result = check_rows(args.document)
        elif args.operation == "parts":
            result = check_parts(args.source, args.document)
        else:
            result = check_relationships(args.document, args.source)
    except AssertionError as error:
        print(str(error), file=sys.stderr)
        sys.exit(1)
    print(result)
