#!/usr/bin/env bash
# Demo: REST API — all three image-supply strategies.
#
# Starts the server, runs all three strategies, then shuts down.
#
# Usage (from repo root):
#   ./scripts/demo-rest.sh
set -euo pipefail
cd "$(dirname "$0")/.."

PORT=8080
HOST="http://localhost:$PORT"
SIGN=".demo/demo-signature-1280x720-rectangle.png"
SEAL=".demo/demo-logo-1024x1024-square.png"

# ── Shared tax info JSON (used by all three strategies) ───────────────────────
read -r -d '' TAX_INFO_JSON <<'JSON' || true
{
  "documentDetails": { "bookNumber": "001", "documentNumber": "WHT-001-2568" },
  "payer": {
    "taxId": "1234567890123",
    "taxId10Digit": "1234567890",
    "name": "บริษัท ตัวอย่าง จำกัด",
    "address": "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110"
  },
  "payee": {
    "taxId": "3210987654321",
    "taxId10Digit": "1234567890",
    "name": "นางสาวสมชาย นามสกุลยาวมาก",
    "address": "555 ต.ทุ่งนา อ.ทุ่งนา จ.ชลบุรี 20000",
    "sequenceNumber": "1",
    "pnd_3": true
  },
  "income40_1": { "datePaid": "01 มกราคม 2568", "amountPaid": "100,000.00", "taxWithheld": "3,000.00" },
  "totals": {
    "totalAmountPaid": "100,000.00",
    "totalTaxWithheld": "3,000.00",
    "totalTaxWithheldInWords": "สามพันบาทถ้วน"
  },
  "withholdingType": { "withholdingTax": true },
  "certification": { "dateOfIssuance": { "day": "1", "month": "มกราคม", "year": "2568" } }
}
JSON

# ── Start the server ──────────────────────────────────────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " Starting REST server on port $PORT ..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

go run ./cmd/rest &>/tmp/pdf50tawi-rest.log &
SERVER_PID=$!
trap 'echo ""; echo "Stopping server (PID $SERVER_PID)..."; kill "$SERVER_PID" 2>/dev/null; wait "$SERVER_PID" 2>/dev/null; exit' EXIT

# Wait until the server accepts connections
MAX_WAIT=15
ELAPSED=0
until curl -sf "$HOST/" &>/dev/null || [[ $ELAPSED -ge $MAX_WAIT ]]; do
  sleep 0.3
  ELAPSED=$(( ELAPSED + 1 ))
done
# Give it one more moment after the port opens
sleep 0.5
echo " Server ready."
echo ""

# ── Strategy A: multipart/form-data ──────────────────────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " Strategy A — multipart/form-data"
echo " POST $HOST/api/v1/taxes/multipart"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

curl -s -X POST "$HOST/api/v1/taxes/multipart" \
  -F "taxInfo=$TAX_INFO_JSON" \
  -F "signature=@$SIGN;type=image/png" \
  -F "seal=@$SEAL;type=image/png" \
  -o certificate-multipart.pdf \
  -w "HTTP %{http_code} — %{size_download} bytes\n"

echo " Output: certificate-multipart.pdf"
echo ""

# ── Strategy B: JSON body with base64-encoded images ─────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " Strategy B — base64 images in JSON"
echo " POST $HOST/api/v1/taxes/base64"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

SIGN_B64=$(base64 < "$SIGN" | tr -d '\n')
SEAL_B64=$(base64 < "$SEAL" | tr -d '\n')

# Build JSON with base64 images using Python to avoid escaping issues
BODY=$(python3 - <<PYEOF
import json, sys

with open("$SIGN", "rb") as f:
    import base64
    sign_b64 = base64.b64encode(f.read()).decode()

with open("$SEAL", "rb") as f:
    seal_b64 = base64.b64encode(f.read()).decode()

tax_info = json.loads("""$TAX_INFO_JSON""")
payload = {
    "taxInfo": tax_info,
    "signatureBase64": sign_b64,
    "sealBase64": seal_b64,
}
print(json.dumps(payload))
PYEOF
)

curl -s -X POST "$HOST/api/v1/taxes/base64" \
  -H "Content-Type: application/json" \
  -d "$BODY" \
  -o certificate-base64.pdf \
  -w "HTTP %{http_code} — %{size_download} bytes\n"

echo " Output: certificate-base64.pdf"
echo ""

# ── Strategy C: JSON body with image URLs ─────────────────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " Strategy C — image URLs (server fetches)"
echo " POST $HOST/api/v1/taxes/url"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Spin up a temporary local file server to serve the demo images.
# In production, replace these with real CDN/storage URLs.
FILE_SERVER_PORT=9090
FILE_SERVER_ROOT=".demo"
python3 -m http.server "$FILE_SERVER_PORT" --directory "$FILE_SERVER_ROOT" &>/dev/null &
FILE_SERVER_PID=$!
trap 'kill "$SERVER_PID" "$FILE_SERVER_PID" 2>/dev/null; wait 2>/dev/null' EXIT
sleep 0.5  # let the file server start

SIGN_FILENAME=$(basename "$SIGN")
SEAL_FILENAME=$(basename "$SEAL")
SIGN_URL="http://localhost:$FILE_SERVER_PORT/$SIGN_FILENAME"
SEAL_URL="http://localhost:$FILE_SERVER_PORT/$SEAL_FILENAME"

BODY=$(python3 - <<PYEOF
import json
tax_info = json.loads("""$TAX_INFO_JSON""")
payload = {
    "taxInfo": tax_info,
    "signatureURL": "$SIGN_URL",
    "sealURL": "$SEAL_URL",
}
print(json.dumps(payload))
PYEOF
)

curl -s -X POST "$HOST/api/v1/taxes/url" \
  -H "Content-Type: application/json" \
  -d "$BODY" \
  -o certificate-url.pdf \
  -w "HTTP %{http_code} — %{size_download} bytes\n"

kill "$FILE_SERVER_PID" 2>/dev/null || true

echo " Output: certificate-url.pdf"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " All three strategies completed successfully."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
