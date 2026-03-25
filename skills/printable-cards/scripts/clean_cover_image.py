#!/usr/bin/env python3
"""
Clean an AI-generated card cover image by removing background artifacts
(table surfaces, shadows, grey borders) and cropping to the card content only.

Usage:
  python clean_cover_image.py <input_image> <output_image>

Dependencies: Pillow, numpy
"""

import sys
import numpy as np
from PIL import Image


def find_card_boundary(arr, h, w):
    """Scan from each edge at multiple positions to find the card content area."""

    def is_content(r, g, b, bg_brightness=200):
        """Heuristic: content pixels differ significantly from a white/grey/brown background."""
        brightness = (int(r) + int(g) + int(b)) / 3
        spread = max(int(r), int(g), int(b)) - min(int(r), int(g), int(b))
        # Dark content (charcoal, navy, etc.)
        if brightness < 120 and spread < 40:
            return True
        # Colored content (not grey/brown background)
        if brightness < 200 and spread > 50:
            return True
        return False

    def scan_edge(positions, scan_range, get_pixel):
        results = []
        for pos in positions:
            for i in scan_range:
                r, g, b = get_pixel(pos, i)
                if is_content(r, g, b):
                    results.append(i)
                    break
        return results

    y_positions = [int(h * p) for p in [0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8]]
    x_positions = [int(w * p) for p in [0.3, 0.4, 0.5, 0.6, 0.7]]

    lefts = scan_edge(y_positions, range(w // 2), lambda y, x: arr[y, x])
    rights = scan_edge(y_positions, range(w - 1, w // 2, -1), lambda y, x: arr[y, x])
    tops = scan_edge(x_positions, range(h // 2), lambda x, y: arr[y, x])
    bottoms = scan_edge(x_positions, range(h - 1, h // 2, -1), lambda x, y: arr[y, x])

    margin = 20
    left = max(lefts) + margin if lefts else int(w * 0.1)
    right = min(rights) - margin if rights else int(w * 0.9)
    top = max(tops) + margin if tops else int(h * 0.05)
    bottom = min(bottoms) - margin if bottoms else int(h * 0.95)

    return left, top, right, bottom


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <input_image> <output_image>")
        sys.exit(1)

    input_path, output_path = sys.argv[1], sys.argv[2]
    img = Image.open(input_path)
    arr = np.array(img)
    h, w = arr.shape[:2]

    left, top, right, bottom = find_card_boundary(arr, h, w)
    cropped = img.crop((left, top, right, bottom))
    cropped.save(output_path)

    cw, ch = cropped.size
    print(f"Cleaned: {w}x{h} -> {cw}x{ch}")
    print(f"Saved: {output_path}")


if __name__ == "__main__":
    main()
