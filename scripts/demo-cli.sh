#!/usr/bin/env bash
# Demo: CLI strategy
# Passes signature and seal as local file paths via --signature / --seal flags.
#
# Usage (from repo root):
#   ./scripts/demo-cli.sh                                    # uses bundled demo images
#   ./scripts/demo-cli.sh path/to/sign.png path/to/seal.png # your own images
set -euo pipefail
cd "$(dirname "$0")/.."

SIGN="${1:-cmd/demo-cli/demo-signature-1280x720-rectangle.png}"
SEAL="${2:-cmd/demo-cli/demo-logo-1024x1024-square.png}"
OUTPUT="certificate-cli.pdf"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " CLI Demo — local file paths"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " signature : $SIGN"
echo " seal      : $SEAL"
echo " output    : $OUTPUT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

go run ./cmd/cli \
  --signature "$SIGN" \
  --seal      "$SEAL" \
  --output    "$OUTPUT"
