#!/usr/bin/env python3
"""Verify en/bn message catalogs have matching keys (broken-key CI check)."""
import json
import sys


def flat(d, prefix=""):
    out = {}
    for k, v in d.items():
        key = f"{prefix}.{k}" if prefix else k
        if isinstance(v, dict):
            out.update(flat(v, key))
        else:
            out[key] = v
    return out


def main():
    with open("apps/console-web/messages/en.json") as f:
        en = json.load(f)
    with open("apps/console-web/messages/bn.json") as f:
        bn = json.load(f)

    fe, fb = flat(en), flat(bn)
    missing = set(fe) - set(fb)
    extra = set(fb) - set(fe)

    if missing:
        print(f"ERROR: bn missing keys: {sorted(missing)}")
        sys.exit(1)
    if extra:
        print(f"warning: bn has extra keys: {sorted(extra)}")
    print(f"i18n OK: {len(fe)} keys en/bn in sync")


if __name__ == "__main__":
    main()
