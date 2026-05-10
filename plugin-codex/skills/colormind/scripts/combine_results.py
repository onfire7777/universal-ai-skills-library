#!/usr/bin/env python3
"""Combine sampled colors and Colormind palette into final JSON."""
import json
import sys

def combine_results(sampled_path, colormind_path):
    """Read two JSON files and combine them."""
    with open(sampled_path, 'r') as f:
        sampled = json.load(f)
    
    with open(colormind_path, 'r') as f:
        colormind = json.load(f)
    
    output = {
        "sampled": {
            "colors": sampled,
            "base": sampled[0] if sampled else None,
        },
        "colormind": colormind,
    }
    
    return output

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: combine_results.py <sampled.json> <colormind.json>", file=sys.stderr)
        sys.exit(1)
    
    result = combine_results(sys.argv[1], sys.argv[2])
    print(json.dumps(result, indent=2))
