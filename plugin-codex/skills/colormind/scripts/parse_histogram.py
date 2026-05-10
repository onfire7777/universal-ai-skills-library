#!/usr/bin/env python3
"""Parse ImageMagick histogram output into JSON."""
import json
import re
import sys

def parse_histogram(text):
    """Parse histogram lines into color records."""
    pattern = re.compile(r'^(\d+):\s+\((\d+),(\d+),(\d+)\)\s+(#[0-9A-Fa-f]{6})')
    colors = []
    
    for line in text.splitlines():
        match = pattern.search(line.strip())
        if not match:
            continue
        
        count = int(match.group(1))
        r, g, b = int(match.group(2)), int(match.group(3)), int(match.group(4))
        hex_color = match.group(5).upper()
        
        colors.append({
            "count": count,
            "rgb": [r, g, b],
            "hex": hex_color
        })
    
    # Sort by count descending
    colors.sort(key=lambda x: x["count"], reverse=True)
    return colors

if __name__ == "__main__":
    text = sys.stdin.read()
    colors = parse_histogram(text)
    print(json.dumps(colors))
