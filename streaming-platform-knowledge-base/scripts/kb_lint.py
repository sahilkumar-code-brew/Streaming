#!/usr/bin/env python3
import re
import subprocess
import sys
from datetime import date, timedelta
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
TODAY = date.today()
STALE_AFTER = timedelta(days=30)
OPEN_LOG_STALE_AFTER = timedelta(days=14)
MAX_ENTRY_LINES = 60
MAX_DOC_LINES = 300
MAX_LINE_CHARS = 500
DOC_TYPES = {"registry", "architecture", "instruction", "feature", "log"}
VERIFIED_BY = {"code", "docs", "live"}
SURFACES = {"backend", "mobile", "admin", "cross"}
ID_PATTERN = r"(?:LL|CC|C|D|R)-\d{3}"
LOG_NAME = re.compile(r"^\d{4}-\d{2}-\d{2}-(backend|mobile|admin|cross)-[a-z0-9-]+\.md$")
PATH_REF = re.compile(r"`([\w./-]+/[\w.-]+\.(?:md|sh|py))`|\]\(([^)#]+\.md)\)")
REGISTRY_ROW = re.compile(rf"^\|\s*({ID_PATTERN})\s*\|", re.M)

failures = []
warnings = []


def fail(path, message):
    failures.append(f"{path.relative_to(ROOT)}: {message}")


def warn(path, message):
    warnings.append(f"{path.relative_to(ROOT)}: {message}")


def parse_frontmatter(text):
    match = re.match(r"^---\n(.*?)\n---\n", text, re.S)
    if not match:
        return None
    pairs = (line.split(":", 1) for line in match.group(1).splitlines() if ":" in line)
    return {key.strip(): value.strip() for key, value in pairs}


def parse_date(value):
    try:
        return date.fromisoformat(value)
    except (TypeError, ValueError):
        return None


def kb_docs():
    excluded = {"templates", ".git"}
    return sorted(p for p in ROOT.rglob("*.md") if not excluded & set(p.relative_to(ROOT).parts))


def owner_roles():
    text = (ROOT / "registries/owners.md").read_text()
    return set(re.findall(r"^\| ([a-z-]+) \|", text, re.M)) - {"role"}


def defined_ids():
    ids = set()
    for path in (ROOT / "registries").glob("*.md"):
        ids.update(REGISTRY_ROW.findall(path.read_text()))
    return ids


def check_sizes(path, lines):
    limit = MAX_ENTRY_LINES if path.name == "CLAUDE.md" else MAX_DOC_LINES
    if len(lines) > limit:
        fail(path, f"{len(lines)} lines, cap is {limit} (LL-005)")
    for number, line in enumerate(lines, 1):
        if len(line) > MAX_LINE_CHARS:
            fail(path, f"line {number} is {len(line)} chars, cap is {MAX_LINE_CHARS}")


def check_frontmatter(path, meta, roles):
    if meta is None:
        fail(path, "missing frontmatter")
        return
    doc_type = meta.get("type")
    if doc_type not in DOC_TYPES:
        fail(path, f"type '{doc_type}' not in {sorted(DOC_TYPES)}")
    if meta.get("owner") not in roles:
        fail(path, f"owner '{meta.get('owner')}' not in registries/owners.md")
    if doc_type != "registry" and meta.get("verified-by") not in VERIFIED_BY:
        fail(path, f"verified-by must be one of {sorted(VERIFIED_BY)}")
    surfaces = set(re.findall(r"[a-z]+", meta.get("surfaces", "")))
    if doc_type != "registry" and not surfaces <= SURFACES:
        fail(path, f"unknown surfaces {sorted(surfaces - SURFACES)}")
    if doc_type not in {"registry", "log"}:
        check_freshness(path, meta.get("last-verified"), "last-verified")
    if doc_type == "registry":
        check_freshness(path, meta.get("last-updated"), "last-updated")


def check_freshness(path, value, field):
    stamped = parse_date(value)
    if stamped is None:
        fail(path, f"{field} missing or not YYYY-MM-DD")
    elif TODAY - stamped > STALE_AFTER:
        warn(path, f"{field} {stamped} is older than {STALE_AFTER.days} days; re-verify against code")


def check_log_entry(path, meta):
    if meta is None or meta.get("type") != "log":
        return
    if not LOG_NAME.match(path.name):
        fail(path, "log name must be YYYY-MM-DD-<backend|mobile|admin|cross>-<topic>.md")
    status = meta.get("status")
    if status not in {"open", "closed"}:
        fail(path, "status must be open or closed")
    opened = parse_date(meta.get("opened"))
    if opened is None:
        fail(path, "opened missing or not YYYY-MM-DD")
    if status == "open" and opened and TODAY - opened > OPEN_LOG_STALE_AFTER:
        warn(path, f"open since {opened}; close it or update it")
    if status == "closed":
        check_rolled_up(path, meta)


def check_rolled_up(path, meta):
    closed = parse_date(meta.get("closed"))
    if closed is None:
        fail(path, "closed entry needs a closed: date")
        return
    feature = meta.get("feature", "")
    if feature == "none":
        return
    feature_path = ROOT / feature
    if not feature_path.is_file():
        fail(path, f"feature '{feature}' does not exist")
        return
    feature_meta = parse_frontmatter(feature_path.read_text()) or {}
    verified = parse_date(feature_meta.get("last-verified"))
    if verified is None or verified < closed:
        fail(path, f"closed {closed} but {feature} last-verified {verified}; roll the outcome up (LL-007)")


def check_registry_copies(path, text):
    if path.parent.name == "registries":
        return
    for row_id in REGISTRY_ROW.findall(text):
        fail(path, f"registry row {row_id} copied outside registries/ (LL-006)")


def check_id_references(path, text, known_ids):
    for line in text.splitlines():
        if "Next ID" in line:
            continue
        for ref in set(re.findall(rf"\b{ID_PATTERN}\b", line)) - known_ids:
            fail(path, f"references {ref}, which no registry defines")


def check_path_references(path, text):
    for match in PATH_REF.finditer(text):
        target = match.group(1) or match.group(2)
        if any(token in target for token in ("<", "*", "{", "YYYY")):
            continue
        if not ((path.parent / target).exists() or (ROOT / target).exists()):
            fail(path, f"broken path reference {target}")


def check_next_ids():
    for path in (ROOT / "registries").glob("*.md"):
        text = path.read_text()
        stated = re.search(r"Next ID: ([A-Z]+)-(\d{3})", text)
        if not stated:
            continue
        prefix = stated.group(1)
        numbers = [int(n) for n in re.findall(rf"^\|\s*{prefix}-(\d{{3}})\s*\|", text, re.M)]
        expected = max(numbers, default=0) + 1
        if int(stated.group(2)) != expected:
            fail(path, f"Next ID {prefix}-{stated.group(2)}, expected {prefix}-{expected:03d}")


def run_contract_sync():
    result = subprocess.run([str(ROOT / "scripts/check_kb_sync.sh")], capture_output=True, text=True)
    for line in result.stdout.splitlines():
        if line.startswith("FAIL"):
            failures.append(f"check_kb_sync: {line}")
        elif line.startswith("WARN"):
            warnings.append(f"check_kb_sync: {line}")


def main():
    roles = owner_roles()
    known_ids = defined_ids()
    for path in kb_docs():
        text = path.read_text()
        check_sizes(path, text.splitlines())
        check_registry_copies(path, text)
        check_id_references(path, text, known_ids)
        check_path_references(path, text)
        if path.name == "CLAUDE.md":
            continue
        meta = parse_frontmatter(text)
        check_frontmatter(path, meta, roles)
        check_log_entry(path, meta)
    check_next_ids()
    run_contract_sync()
    for line in warnings:
        print(f"WARN  {line}")
    for line in failures:
        print(f"FAIL  {line}")
    print(f"{len(failures)} failures, {len(warnings)} warnings")
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
