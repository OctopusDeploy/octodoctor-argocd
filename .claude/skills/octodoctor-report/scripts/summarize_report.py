#!/usr/bin/env python3
"""Summarize an octodoctor report.json into a compact, readable digest.

Usage: summarize_report.py <path-to-report.json>
"""
import json
import sys

SEVERITY_ORDER = {"high": 0, "medium": 1, "low": 2}


def main() -> None:
    if len(sys.argv) != 2:
        print("usage: summarize_report.py <path-to-report.json>", file=sys.stderr)
        sys.exit(2)

    with open(sys.argv[1]) as f:
        report = json.load(f)

    meta = report["metadata"]
    summary = report["summary"]
    rules = summary["rules"]
    resources = summary["resources"]

    print(f"octodoctor {meta['version']} — scanned {meta['scannedDirectory']}")
    print(
        f"Files: {summary['totalAnalyzedFiles']}/{summary['totalYamlFiles']} analyzed "
        f"(Applications={resources['applications']}, ApplicationSets={resources['applicationSets']}, "
        f"AppProjects={resources['appProjects']}, ConfigMaps={resources['configmaps']}, "
        f"Secrets={resources['secrets']})"
    )
    print(
        f"Rules: {rules['total']} total — "
        f"{rules['failed']} failed, {rules['passed']} passed, {rules['inactive']} inactive"
    )

    failed = sorted(
        report["results"]["failed"],
        key=lambda r: SEVERITY_ORDER.get(r["severity"], 99),
    )
    if failed:
        print("\nFailed rules (most severe first):")
        for r in failed:
            files = ", ".join(r["affectedFiles"][:5])
            more = f" (+{len(r['affectedFiles']) - 5} more)" if len(r["affectedFiles"]) > 5 else ""
            print(f"  [{r['severity'].upper()}] {r['id']} — {r['name']}")
            print(f"      {r['description']}")
            print(f"      affected: {files}{more}")
            print(f"      docs: {r['documentationUrl']}")
    else:
        print("\nNo failed rules.")


if __name__ == "__main__":
    main()
