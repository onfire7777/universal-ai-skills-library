#!/usr/bin/env python3
"""Extract base RGB color from sampled colors JSON."""
import json
import sys

def get_base_rgb(colors_path):
    """Read sampled colors and return base RGB as comma-separated string."""
    with open(colors_path, 'r') as f:
        colors = json.load(f)
    
    if not colors:
        return "0,0,0"
    
    base_rgb = colors[0]["rgb"]
    return ",".join(map(str, base_rgb))

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: get_base_rgb.py <sampled.json>", file=sys.stderr)
        sys.exit(1)
    
    rgb = get_base_rgb(sys.argv[1])
    print(rgb)
