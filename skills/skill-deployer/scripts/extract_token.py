#!/usr/bin/env python3
"""
Extract JWT session token from a provider browser session.

This script provides the JavaScript snippet to run in the browser console
to extract the authentication token needed for API calls.

Usage:
    python3 extract_token.py           # Print the JS snippet
    python3 extract_token.py --verify TOKEN  # Verify a token works
"""
import argparse
import json
import os
import sys
from urllib.parse import urlparse


def normalize_api_base(value):
    """Return a safe API base URL for token-bearing provider calls."""
    base = (value or "").strip().rstrip("/")
    parsed = urlparse(base)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        raise ValueError("--api-base must be an absolute http(s) URL")
    if parsed.scheme == "http" and parsed.hostname not in {"localhost", "127.0.0.1", "::1"}:
        raise ValueError("--api-base must use https unless targeting localhost")
    return base


DEFAULT_API_BASE = os.environ.get("SKILL_DEPLOYER_API_BASE", "https://api.manus.im")

JS_SNIPPET = """
// Run this in the provider web app browser console while logged in:
// 1. Open DevTools (F12)
// 2. Go to Console tab
// 3. Paste and run:

(function() {
    const cookies = document.cookie.split(';');
    for (const c of cookies) {
        const [name, ...val] = c.trim().split('=');
        if (name === 'session_id' || name === 'token') {
            console.log('TOKEN:', val.join('='));
            return val.join('=');
        }
    }
    // Fallback: check localStorage
    for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key.includes('token') || key.includes('session') || key.includes('auth')) {
            const val = localStorage.getItem(key);
            if (val && val.length > 50) {
                console.log('TOKEN:', val);
                return val;
            }
        }
    }
    console.log('No token found. Make sure you are logged in to the provider web app');
})();
""".strip()


def verify_token(token, api_base=None):
    """Verify a token works by listing skills for a test call."""
    api_base = normalize_api_base(api_base or DEFAULT_API_BASE)
    try:
        import requests
    except ImportError:
        print("requests not installed")
        return False

    resp = requests.post(
        f"{api_base}/skill.v1.ProjectSkillService/ListProjectSkills",
        json={"project_uid": "test"},
        headers={
            "Content-Type": "application/json",
            "Connect-Protocol-Version": "1",
            "Authorization": f"Bearer {token}",
        },
        timeout=10,
    )
    # A 400 (bad project uid) means auth worked; 401/403 means token is bad
    if resp.status_code in (200, 400):
        print("Token is VALID (authentication successful)")
        return True
    elif resp.status_code in (401, 403):
        print("Token is INVALID (authentication failed)")
        return False
    else:
        print(f"Unexpected status: {resp.status_code} — {resp.text[:200]}")
        return False


def main():
    parser = argparse.ArgumentParser(description="Extract/verify provider JWT token")
    parser.add_argument("--verify", help="Verify a token")
    parser.add_argument("--api-base", default=DEFAULT_API_BASE, help="Provider API base URL")
    args = parser.parse_args()

    if args.verify:
        try:
            api_base = normalize_api_base(args.api_base)
        except ValueError as e:
            parser.error(str(e))
        verify_token(args.verify, api_base)
    else:
        print("=== Provider JWT Token Extraction ===\n")
        print("Run this JavaScript in your provider web app browser console:\n")
        print(JS_SNIPPET)
        print("\n\nAlternatively, browser automation can extract it via browser_console_exec:")
        print("  document.cookie.split(';').find(c => c.trim().startsWith('session_id=')).split('=').slice(1).join('=')")


if __name__ == "__main__":
    main()
