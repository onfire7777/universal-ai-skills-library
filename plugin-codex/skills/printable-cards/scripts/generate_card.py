#!/usr/bin/env python3
"""
Generate a printable foldable greeting card as a 2-page PDF.

Layout: US Letter landscape (11x8.5in).
  Page 1 (outside): Left = back cover, Right = front cover image
  Page 2 (inside):  Left = open space for handwriting, Right = printed message

Usage:
  python generate_card.py \
    --name "Chloe" \
    --front-image /path/to/front_cover.png \
    --message-file /path/to/message.html \
    --closing "All my love, your little brother" \
    --accent "#C4889B" \
    --accent-dark "#9E6B7E" \
    --light-bg "#FDF2F5" \
    --dark-text "#5A2D45" \
    --output /path/to/output.pdf \
    [--masculine]

Dependencies: weasyprint (pip install weasyprint)
"""

import argparse
import os
import sys

def build_html(name, front_img, accent, accent_dark, light_bg, dark_text,
               message_html, closing, is_masculine=False):

    if is_masculine:
        title_css = f"font-family: 'Noto Sans', Arial, sans-serif; font-weight: 700; font-size: 22pt; color: {dark_text}; letter-spacing: 1px;"
        name_css = f"font-family: 'Noto Sans', Arial, sans-serif; font-weight: 700; font-size: 30pt; color: {accent}; letter-spacing: 2px;"
        body_css = f"font-family: 'Noto Serif', Georgia, serif; font-size: 10.5pt; line-height: 1.75; color: {dark_text}; text-align: center;"
        closing_css = f"font-family: 'Noto Sans', Arial, sans-serif; font-style: italic; font-size: 10pt; color: {accent};"
        back_label_css = f"font-family: 'Noto Sans', Arial, sans-serif; font-style: italic; font-size: 7.5pt; color: {accent};"
        back_symbol = "&#10022;"
        corner_tl = f'<div style="position:absolute;top:18px;left:18px;width:40px;height:40px;border-top:2px solid {accent}40;border-left:2px solid {accent}40;"></div>'
        corner_tr = f'<div style="position:absolute;top:18px;right:18px;width:40px;height:40px;border-top:2px solid {accent}40;border-right:2px solid {accent}40;"></div>'
        corner_bl = f'<div style="position:absolute;bottom:18px;left:18px;width:40px;height:40px;border-bottom:2px solid {accent}40;border-left:2px solid {accent}40;"></div>'
        corner_br = f'<div style="position:absolute;bottom:18px;right:18px;width:40px;height:40px;border-bottom:2px solid {accent}40;border-right:2px solid {accent}40;"></div>'
    else:
        title_css = f"font-family: 'Noto Serif', Georgia, serif; font-style: italic; font-size: 22pt; color: {dark_text};"
        name_css = f"font-family: 'Noto Serif', Georgia, serif; font-style: italic; font-size: 30pt; color: {accent};"
        body_css = f"font-family: 'Noto Serif', Georgia, serif; font-style: italic; font-size: 10.5pt; line-height: 1.75; color: {dark_text}; text-align: center;"
        closing_css = f"font-family: 'Noto Serif', Georgia, serif; font-style: italic; font-size: 10pt; color: {accent};"
        back_label_css = f"font-family: 'Noto Serif', Georgia, serif; font-style: italic; font-size: 7.5pt; color: {accent};"
        back_symbol = "&#9829;"
        corner_tl = f'<div style="position:absolute;top:16px;left:16px;width:35px;height:35px;border-top:1.5px solid {accent}50;border-left:1.5px solid {accent}50;"></div>'
        corner_tr = f'<div style="position:absolute;top:16px;right:16px;width:35px;height:35px;border-top:1.5px solid {accent}50;border-right:1.5px solid {accent}50;"></div>'
        corner_bl = f'<div style="position:absolute;bottom:16px;left:16px;width:35px;height:35px;border-bottom:1.5px solid {accent}50;border-left:1.5px solid {accent}50;"></div>'
        corner_br = f'<div style="position:absolute;bottom:16px;right:16px;width:35px;height:35px;border-bottom:1.5px solid {accent}50;border-right:1.5px solid {accent}50;"></div>'

    front_img_path = os.path.abspath(front_img)
    corners_html = corner_tl + corner_tr + corner_bl + corner_br

    return f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  @page {{
    size: 11in 8.5in;
    margin: 0;
  }}
  * {{ margin: 0; padding: 0; box-sizing: border-box; }}
  body {{ margin: 0; padding: 0; font-size: 0; }}

  .page {{
    width: 11in;
    height: 8.5in;
    position: relative;
    page-break-after: always;
    overflow: hidden;
  }}
  .page:last-child {{ page-break-after: auto; }}

  .half-left {{
    position: absolute;
    top: 0; left: 0;
    width: 5.5in;
    height: 8.5in;
    overflow: hidden;
  }}
  .half-right {{
    position: absolute;
    top: 0; left: 5.5in;
    width: 5.5in;
    height: 8.5in;
    overflow: hidden;
  }}

  .fold {{
    position: absolute;
    top: 8px; left: 5.5in;
    width: 0; height: calc(8.5in - 16px);
    border-left: 0.5px dashed #d0d0d0;
    z-index: 100;
  }}

  .print-note {{
    position: absolute;
    bottom: 1px; left: 0; right: 0;
    text-align: center;
    font-family: Arial, sans-serif;
    font-size: 5pt;
    color: #cccccc;
    z-index: 100;
  }}

  .back-cover {{
    background: {light_bg};
    width: 100%; height: 100%;
    position: relative;
  }}
  .back-border-1 {{
    position: absolute;
    top: 14px; left: 14px; right: 14px; bottom: 14px;
    border: 1px solid {accent}60;
    border-radius: 4px;
  }}
  .back-border-2 {{
    position: absolute;
    top: 20px; left: 20px; right: 20px; bottom: 20px;
    border: 0.5px solid {accent}40;
    border-radius: 3px;
  }}
  .back-symbol {{
    position: absolute;
    top: 50%; left: 50%;
    transform: translate(-50%, -50%);
    font-size: 16pt;
    color: {accent};
  }}
  .back-label {{
    position: absolute;
    bottom: 28px; left: 0; right: 0;
    text-align: center;
    {back_label_css}
  }}

  .front-cover {{
    width: 100%; height: 100%;
    background-image: url('file://{front_img_path}');
    background-size: cover;
    background-position: center center;
    background-repeat: no-repeat;
  }}

  .inside-left {{
    background: {light_bg};
    width: 100%; height: 100%;
    position: relative;
  }}

  .inside-right {{
    background: {light_bg};
    width: 100%; height: 100%;
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 35px;
  }}
  .inside-right .header {{
    text-align: center;
    margin-bottom: 6px;
  }}
  .inside-right .title {{
    {title_css}
    margin-bottom: 2px;
  }}
  .inside-right .name-text {{
    {name_css}
    margin-bottom: 8px;
  }}
  .inside-right .divider {{
    width: 50px;
    height: 0;
    border-top: 1.5px solid {accent};
    margin: 0 auto 18px auto;
  }}
  .inside-right .message {{
    {body_css}
    padding: 0 10px;
  }}
  .inside-right .message p {{
    margin-bottom: 12px;
  }}
  .inside-right .closing-div {{
    width: 30px;
    height: 0;
    border-top: 1px solid {accent};
    margin: 14px auto 12px auto;
  }}
  .inside-right .closing {{
    text-align: center;
    {closing_css}
  }}
</style>
</head>
<body>

<div class="page">
  <div class="fold"></div>
  <div class="half-left">
    <div class="back-cover">
      <div class="back-border-1"></div>
      <div class="back-border-2"></div>
      <div class="back-symbol">{back_symbol}</div>
      <div class="back-label">Made with love for {name}</div>
    </div>
  </div>
  <div class="half-right">
    <div class="front-cover"></div>
  </div>
  <div class="print-note">PAGE 1 of 2 &bull; Print double-sided (flip on short edge). Fold in half along center line.</div>
</div>

<div class="page">
  <div class="fold"></div>
  <div class="half-left">
    <div class="inside-left">
      {corners_html}
    </div>
  </div>
  <div class="half-right">
    <div class="inside-right">
      {corners_html}
      <div class="header">
        <div class="title">Happy Birthday,</div>
        <div class="name-text">{name}</div>
      </div>
      <div class="divider"></div>
      <div class="message">
        {message_html}
      </div>
      <div class="closing-div"></div>
      <div class="closing">{closing}</div>
    </div>
  </div>
  <div class="print-note">PAGE 2 of 2 &bull; Print on the back of page 1 (double-sided, flip on short edge).</div>
</div>

</body>
</html>"""


def main():
    parser = argparse.ArgumentParser(description="Generate a printable foldable greeting card PDF.")
    parser.add_argument("--name", required=True, help="Recipient name")
    parser.add_argument("--front-image", required=True, help="Path to front cover image")
    parser.add_argument("--message-file", required=True, help="Path to HTML file with message paragraphs (use <p> tags)")
    parser.add_argument("--closing", required=True, help="Closing line (e.g. 'All my love, your brother')")
    parser.add_argument("--accent", default="#C4889B", help="Accent color hex (e.g. #C4889B)")
    parser.add_argument("--accent-dark", default="#9E6B7E", help="Darker accent color hex")
    parser.add_argument("--light-bg", default="#FDF2F5", help="Light background color hex")
    parser.add_argument("--dark-text", default="#5A2D45", help="Dark text color hex")
    parser.add_argument("--output", required=True, help="Output PDF path")
    parser.add_argument("--masculine", action="store_true", help="Use masculine styling (bold sans-serif)")
    args = parser.parse_args()

    with open(args.message_file, "r") as f:
        message_html = f.read()

    from weasyprint import HTML
    html = build_html(
        name=args.name,
        front_img=args.front_image,
        accent=args.accent,
        accent_dark=args.accent_dark,
        light_bg=args.light_bg,
        dark_text=args.dark_text,
        message_html=message_html,
        closing=args.closing,
        is_masculine=args.masculine,
    )
    base_url = os.path.dirname(os.path.abspath(args.front_image))
    HTML(string=html, base_url=base_url).write_pdf(args.output)
    print(f"Created: {args.output}")


if __name__ == "__main__":
    main()
