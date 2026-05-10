#!/usr/bin/env python3
"""
Manus Skill Deployer — Deploy skills to Manus projects via gRPC-web API.

Usage:
    python3 deploy_skills.py --token TOKEN --project-uid UID [options]
    python3 deploy_skills.py --token TOKEN --plan plan.json [options]

Commands:
    --list                  List skills in a project
    --upload FILE.zip       Upload a single skill zip
    --upload-dir DIR        Upload all zips from a directory
    --delete-name NAME      Delete a skill by name
    --sync NAME [NAME ...]  Full sync: delete unwanted + upload missing
    --plan PLAN.json        Deploy skills to multiple projects per plan file
    --package               Package all skills as zips (no upload)

Required:
    --token TOKEN           JWT session token from Manus browser session
    --project-uid UID       Target project UID (not needed with --plan/--package)

Options:
    --skills-dir DIR        Source skills directory (default: /home/ubuntu/skills)
    --zip-dir DIR           Directory for zip files (default: /tmp/skill_zips)
    --rate-limit SECS       Delay between API calls (default: 0.15)
    --dry-run               Preview changes without executing
    --max-retries N         Max retries on 429 errors (default: 3)
"""
import argparse
import base64
import json
import os
import subprocess
import sys
import time
import zipfile

try:
    import requests
except ImportError:
    ret = subprocess.call([sys.executable, "-m", "pip", "install", "requests", "-q"])
    if ret != 0:
        print("ERROR: Failed to install 'requests'. Run: pip3 install requests", file=sys.stderr)
        sys.exit(1)
    import requests

API_BASE = "https://api.manus.im"
HEADERS_TEMPLATE = {
    "Content-Type": "application/json",
    "Connect-Protocol-Version": "1",
}
MAX_SKILLS_PER_PROJECT = 500


def grpc_call(method, data, token):
    """Make a gRPC-web call to the Manus API."""
    headers = {**HEADERS_TEMPLATE, "Authorization": f"Bearer {token}"}
    url = f"{API_BASE}/skill.v1.ProjectSkillService/{method}"
    try:
        resp = requests.post(url, json=data, headers=headers, timeout=30)
    except requests.RequestException as e:
        return 0, f"Network error: {e}"
    if resp.status_code == 401:
        print("  ERROR: Unauthorized — token expired or invalid. Re-extract from browser.", file=sys.stderr)
    return resp.status_code, resp.text


def list_skills(project_uid, token):
    """List all skills in a project. Returns {name: skill_id} dict."""
    status, text = grpc_call("ListProjectSkills", {"project_uid": project_uid}, token)
    if status == 200:
        try:
            data = json.loads(text)
        except json.JSONDecodeError:
            print(f"  ERROR: Invalid JSON response: {text[:200]}")
            return {}
        return {
            s["skill"]["name"]: s["skill"]["id"]
            for s in data.get("projectSkills", [])
            if "skill" in s
        }
    print(f"  ERROR listing skills: {status} {text[:200]}")
    return {}


def upload_skill(project_uid, zip_path, token, max_retries=3):
    """Upload a skill zip to a project. Returns (success, message)."""
    with open(zip_path, "rb") as f:
        content = base64.b64encode(f.read()).decode()

    for attempt in range(max_retries + 1):
        status, text = grpc_call(
            "UploadProjectSkill",
            {"project_uid": project_uid, "content": content},
            token,
        )
        if status == 200:
            return True, "uploaded"
        if status == 401:
            return False, "unauthorized — token expired, re-extract from browser"
        if "already" in text.lower() or "exist" in text.lower():
            return True, "already exists"
        if status == 429:
            if attempt < max_retries:
                wait = 2 ** (attempt + 1)
                print(f"    Rate limited, waiting {wait}s (attempt {attempt + 1}/{max_retries})...")
                time.sleep(wait)
                continue
            return False, f"rate limited after {max_retries} retries"
        if "limit" in text.lower() or "maximum" in text.lower() or "exceed" in text.lower():
            return False, "project skill limit reached (500 max)"
        if "SKILL_MD_INVALID_FORMAT" in text:
            return False, "invalid SKILL.md format"
        return False, f"{status}: {text[:200]}"

    return False, "max retries exceeded"


def delete_skill(project_uid, skill_id, token, max_retries=3):
    """Delete a skill from a project. Returns (success, message)."""
    for attempt in range(max_retries + 1):
        status, text = grpc_call(
            "DeleteProjectSkill",
            {"project_uid": project_uid, "skill_id": skill_id},
            token,
        )
        if status == 200:
            return True, "deleted"
        if status == 401:
            return False, "unauthorized — token expired, re-extract from browser"
        if status == 429:
            if attempt < max_retries:
                wait = 2 ** (attempt + 1)
                time.sleep(wait)
                continue
            return False, f"rate limited after {max_retries} retries"
        return False, f"{status}: {text[:200]}"

    return False, "max retries exceeded"


def package_skill(skill_path, zip_dir):
    """Package a skill directory into a zip file. Returns zip path or None."""
    skill_name = os.path.basename(skill_path)
    skill_md = os.path.join(skill_path, "SKILL.md")
    if not os.path.isfile(skill_md):
        return None

    os.makedirs(zip_dir, exist_ok=True)
    zip_path = os.path.join(zip_dir, f"{skill_name}.zip")

    try:
        with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zf:
            for root, dirs, files in os.walk(skill_path, followlinks=False):
                # Skip hidden dirs and __pycache__
                dirs[:] = [d for d in dirs if not d.startswith(".") and d != "__pycache__"]
                for file in files:
                    if file.startswith("."):
                        continue
                    file_path = os.path.join(root, file)
                    # Skip broken symlinks and unreadable files
                    if os.path.islink(file_path) and not os.path.exists(file_path):
                        continue
                    if not os.path.isfile(file_path):
                        continue
                    try:
                        arcname = os.path.relpath(file_path, skill_path)
                        zf.write(file_path, arcname)
                    except (OSError, PermissionError):
                        continue
    except Exception:
        if os.path.exists(zip_path):
            os.remove(zip_path)
        return None

    return zip_path


def package_all_skills(skills_dir, zip_dir):
    """Package all skills in a directory. Returns {name: zip_path} dict."""
    result = {}
    for name in sorted(os.listdir(skills_dir)):
        skill_path = os.path.join(skills_dir, name)
        if not os.path.isdir(skill_path):
            continue
        zip_path = package_skill(skill_path, zip_dir)
        if zip_path:
            result[name] = zip_path
    print(f"Packaged {len(result)} skills into {zip_dir}")
    return result


def find_zip(name, zip_dir):
    """Find a zip file for a skill name, handling case differences."""
    path = os.path.join(zip_dir, f"{name}.zip")
    if os.path.exists(path):
        return path
    name_lower = name.lower()
    for f in os.listdir(zip_dir):
        if f.lower() == f"{name_lower}.zip":
            return os.path.join(zip_dir, f)
    return None


def _has_zips(zip_dir):
    """Check if a directory contains any .zip files."""
    return os.path.isdir(zip_dir) and any(f.endswith(".zip") for f in os.listdir(zip_dir))


def sync_project(project_uid, target_skills, zip_dir, token, rate_limit=0.15,
                 dry_run=False, max_retries=3):
    """Sync a project to match the target skill set. Returns results dict."""
    existing = list_skills(project_uid, token)
    existing_names = set(existing.keys())
    target_set = set(target_skills)

    to_delete = existing_names - target_set
    to_add = target_set - existing_names
    to_keep = existing_names & target_set

    print(f"  Existing: {len(existing)}, Target: {len(target_set)}")
    print(f"  Keep: {len(to_keep)}, Delete: {len(to_delete)}, Add: {len(to_add)}")

    results = {"kept": len(to_keep), "deleted": 0, "uploaded": 0,
               "failed": 0, "skipped": 0, "errors": []}

    if dry_run:
        print("  [DRY RUN] Would delete:", sorted(to_delete)[:10], "..." if len(to_delete) > 10 else "")
        print("  [DRY RUN] Would add:", sorted(to_add)[:10], "..." if len(to_add) > 10 else "")
        return results

    # Delete unwanted skills
    if to_delete:
        print(f"  Deleting {len(to_delete)} skills...")
        for name in sorted(to_delete):
            skill_id = existing[name]
            ok, msg = delete_skill(project_uid, skill_id, token, max_retries)
            if ok:
                results["deleted"] += 1
            else:
                results["errors"].append(f"delete {name}: {msg}")
                if "unauthorized" in msg:
                    print("    ABORTING: Token expired")
                    break
            if results["deleted"] % 50 == 0 and results["deleted"] > 0:
                print(f"    Deleted {results['deleted']}/{len(to_delete)}...")
            time.sleep(rate_limit)
        print(f"    Deleted: {results['deleted']}")

    # Upload new skills — capacity = 500 minus what we're keeping (deleted are gone)
    if to_add:
        remaining_capacity = MAX_SKILLS_PER_PROJECT - len(to_keep)
        if len(to_add) > remaining_capacity:
            print(f"  WARNING: {len(to_add)} to add but only {remaining_capacity} slots available")

        print(f"  Uploading {len(to_add)} skills...")
        upload_count = 0
        for name in sorted(to_add):
            zip_path = find_zip(name, zip_dir)
            if not zip_path:
                results["skipped"] += 1
                continue

            ok, msg = upload_skill(project_uid, zip_path, token, max_retries)
            if ok:
                results["uploaded"] += 1
            else:
                results["failed"] += 1
                results["errors"].append(f"upload {name}: {msg}")
                if "limit reached" in msg:
                    print(f"    HIT PROJECT LIMIT at {results['uploaded']} uploads")
                    break
                if "unauthorized" in msg:
                    print("    ABORTING: Token expired")
                    break

            upload_count += 1
            if upload_count % 50 == 0:
                print(f"    Progress: {results['uploaded']} uploaded, {results['failed']} failed, {results['skipped']} skipped")

            time.sleep(rate_limit)

        print(f"    Uploaded: {results['uploaded']}, Failed: {results['failed']}, Skipped: {results['skipped']}")

    # Verify final count
    final = list_skills(project_uid, token)
    results["final_count"] = len(final)
    print(f"  FINAL: {len(final)} skills")

    return results


def deploy_plan(plan_path, zip_dir, token, rate_limit=0.15, dry_run=False, max_retries=3):
    """Deploy skills to multiple projects per a plan file."""
    with open(plan_path) as f:
        plan = json.load(f)

    all_results = {}
    for proj_name, proj_data in plan.items():
        if proj_name.startswith("_"):
            continue
        uid = proj_data["uid"]
        skills = proj_data["skills"]

        print(f"\n{'=' * 60}")
        print(f"PROJECT: {proj_name}")
        print(f"UID: {uid}")
        print(f"Target: {len(skills)} skills")
        print(f"{'=' * 60}")

        results = sync_project(uid, skills, zip_dir, token, rate_limit, dry_run, max_retries)
        all_results[proj_name] = {**results, "uid": uid, "target": len(skills)}

    # Summary
    print(f"\n{'=' * 60}")
    print("DEPLOYMENT SUMMARY")
    print(f"{'=' * 60}")
    for name, r in all_results.items():
        status = "OK" if r["failed"] == 0 else f"{r['failed']} FAILED"
        print(f"  {name:30s} | Target: {r['target']:3d} | Final: {r.get('final_count', '?'):>3} | {status}")

    return all_results


def main():
    parser = argparse.ArgumentParser(description="Deploy skills to Manus projects via API")
    parser.add_argument("--token", required=True, help="JWT session token")
    parser.add_argument("--project-uid", help="Target project UID")
    parser.add_argument("--skills-dir", default="/home/ubuntu/skills", help="Source skills directory")
    parser.add_argument("--zip-dir", default="/tmp/skill_zips", help="Directory for zip files")
    parser.add_argument("--rate-limit", type=float, default=0.15, help="Delay between API calls")
    parser.add_argument("--max-retries", type=int, default=3, help="Max retries on 429")
    parser.add_argument("--dry-run", action="store_true", help="Preview without executing")

    # Commands
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--list", action="store_true", help="List skills in project")
    group.add_argument("--upload", help="Upload a single skill zip")
    group.add_argument("--upload-dir", help="Upload all zips from directory")
    group.add_argument("--upload-skill", help="Package and upload a skill by name from skills-dir")
    group.add_argument("--delete-name", help="Delete skill by name")
    group.add_argument("--sync", nargs="+", help="Sync project to match skill names list")
    group.add_argument("--plan", help="Deploy per plan JSON file")
    group.add_argument("--package", action="store_true", help="Package all skills as zips (no upload)")

    args = parser.parse_args()

    if args.package:
        zips = package_all_skills(args.skills_dir, args.zip_dir)
        print(f"Done. {len(zips)} zips in {args.zip_dir}")
        return

    if args.plan:
        # Ensure zips exist — check for actual .zip files, not just any files
        if not _has_zips(args.zip_dir):
            print("Packaging skills first...")
            package_all_skills(args.skills_dir, args.zip_dir)
        results = deploy_plan(args.plan, args.zip_dir, args.token, args.rate_limit, args.dry_run, args.max_retries)
        base, _ = os.path.splitext(args.plan)
        out_path = f"{base}_results.json"
        with open(out_path, "w") as f:
            json.dump(results, f, indent=2)
        print(f"\nResults saved to {out_path}")
        return

    # All other commands require --project-uid
    if not args.project_uid:
        parser.error("--project-uid is required for this command")

    if args.list:
        skills = list_skills(args.project_uid, args.token)
        print(f"\n{len(skills)} skills in project {args.project_uid}:")
        for name in sorted(skills.keys()):
            print(f"  {name} (id: {skills[name]})")

    elif args.upload:
        ok, msg = upload_skill(args.project_uid, args.upload, args.token, args.max_retries)
        print(f"{'SUCCESS' if ok else 'FAILED'}: {msg}")

    elif args.upload_dir:
        if not os.path.isdir(args.upload_dir):
            print(f"Upload directory does not exist: {args.upload_dir}")
            sys.exit(1)
        zips = sorted([f for f in os.listdir(args.upload_dir) if f.endswith(".zip")])
        if not zips:
            print(f"No .zip files found in: {args.upload_dir}")
            sys.exit(1)
        print(f"Uploading {len(zips)} skills from {args.upload_dir}...")
        uploaded, failed = 0, 0
        for z in zips:
            path = os.path.join(args.upload_dir, z)
            ok, msg = upload_skill(args.project_uid, path, args.token, args.max_retries)
            if ok:
                uploaded += 1
            else:
                failed += 1
                print(f"  FAILED {z}: {msg}")
                if "unauthorized" in msg:
                    print("  ABORTING: Token expired")
                    break
            time.sleep(args.rate_limit)
        print(f"Done. Uploaded: {uploaded}, Failed: {failed}")

    elif args.upload_skill:
        skill_path = os.path.join(args.skills_dir, args.upload_skill)
        if not os.path.isdir(skill_path):
            print(f"Skill not found: {skill_path}")
            sys.exit(1)
        zip_path = package_skill(skill_path, args.zip_dir)
        if not zip_path:
            print(f"No SKILL.md found in {skill_path}")
            sys.exit(1)
        ok, msg = upload_skill(args.project_uid, zip_path, args.token, args.max_retries)
        print(f"{'SUCCESS' if ok else 'FAILED'}: {msg}")

    elif args.delete_name:
        skills = list_skills(args.project_uid, args.token)
        if args.delete_name in skills:
            ok, msg = delete_skill(args.project_uid, skills[args.delete_name], args.token, args.max_retries)
            print(f"{'SUCCESS' if ok else 'FAILED'}: {msg}")
        else:
            print(f"Skill '{args.delete_name}' not found in project")

    elif args.sync:
        # Ensure zips exist
        if not _has_zips(args.zip_dir):
            print("Packaging skills first...")
            package_all_skills(args.skills_dir, args.zip_dir)
        results = sync_project(
            args.project_uid, args.sync, args.zip_dir, args.token,
            args.rate_limit, args.dry_run, args.max_retries
        )
        print(json.dumps(results, indent=2))


if __name__ == "__main__":
    main()
