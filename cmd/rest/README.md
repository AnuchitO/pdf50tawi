# cmd/rest — REST API server ตัวอย่าง / Example REST server

server ตัวอย่างที่แสดงวิธีใช้ `pdf50tawi` ผ่าน HTTP API โดยรองรับ 3 วิธีในการส่งรูปภาพลายเซ็นและตราประทับ

An example HTTP server showing how to integrate `pdf50tawi` into a REST API. Supports three strategies for supplying signature and seal images.

---

## รันเซิร์ฟเวอร์ / Start the server

```bash
go run ./cmd/rest
```

เปลี่ยน port ได้ผ่าน environment variable (default: `8080`):

```bash
PORT=9000 go run ./cmd/rest
```

---

## รัน demo ทั้งหมดในคราวเดียว / Run full demo

```bash
./scripts/demo-rest.sh
```

สคริปต์จะ start server, เรียก 3 routes, แล้ว shutdown ให้อัตโนมัติ

The script starts the server, calls all three routes, and shuts down automatically.

---

## ข้อมูลตัวอย่าง / Demo tax info JSON

JSON ด้านล่างนี้ใช้ร่วมกันในทุก strategy — ครอบคลุมทุก field ในแบบฟอร์ม 50 ทวิ

The JSON below is shared across all three strategy examples — it covers every field in the form.

<details>
<summary>TAX_INFO_JSON (คลิกเพื่อขยาย / click to expand)</summary>

```json
{
  "documentDetails": { "bookNumber": "001", "documentNumber": "001" },
  "payer": {
    "taxId": "1234567890123",
    "taxId10Digit": "1234567890",
    "name": "บริษัท ตัวอย่าง จำกัด",
    "address": "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110"
  },
  "payee": {
    "taxId": "3210987654321",
    "taxId10Digit": "1234567890",
    "name": "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน",
    "address": "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345",
    "sequenceNumber": "321",
    "pnd_1a": true,
    "pnd_1aSpecial": true,
    "pnd_2": true,
    "pnd_2a": true,
    "pnd_3": true,
    "pnd_3a": true,
    "pnd_53": true
  },
  "income40_1":      { "datePaid": "01 มกราคม 2568", "amountPaid": "401,010.01", "taxWithheld": "12,030.30" },
  "income40_2":      { "datePaid": "02 ก.พ. 2568",   "amountPaid": "402,020.02", "taxWithheld": "12,060.60" },
  "income40_3":      { "datePaid": "03 มี.ค. 2568",  "amountPaid": "403,030.03", "taxWithheld": "12,090.90" },
  "income40_4A":     { "datePaid": "04 เม.ย. 2568",  "amountPaid": "404,040.04", "taxWithheld": "12,121.20" },
  "income40_4B_1_1": { "datePaid": "05 พ.ค. 2568",   "amountPaid": "411,010.01", "taxWithheld": "12,330.30" },
  "income40_4B_1_2": { "datePaid": "06 มิ.ย. 2568",  "amountPaid": "412,020.02", "taxWithheld": "12,360.60" },
  "income40_4B_1_3": { "datePaid": "07 ก.ค. 2568",   "amountPaid": "413,030.03", "taxWithheld": "12,390.90" },
  "income40_4B_1_4_rate": "ร้อยละ 7",
  "income40_4B_1_4": { "datePaid": "08 ส.ค. 2568",   "amountPaid": "414,040.04", "taxWithheld": "12,421.20" },
  "income40_4B_2_1": { "datePaid": "09 ก.ย. 2568",   "amountPaid": "421,010.01", "taxWithheld": "12,630.30" },
  "income40_4B_2_2": { "datePaid": "10 ต.ค. 2568",   "amountPaid": "422,020.02", "taxWithheld": "12,660.60" },
  "income40_4B_2_3": { "datePaid": "11 พ.ย. 2568",   "amountPaid": "423,030.03", "taxWithheld": "12,690.90" },
  "income40_4B_2_4": { "datePaid": "12 ธ.ค. 2568",   "amountPaid": "424,040.04", "taxWithheld": "12,721.20" },
  "income40_4B_2_5_note": "กำไรอื่นๆ",
  "income40_4B_2_5": { "datePaid": "13 ม.ค. 2568",   "amountPaid": "425,050.05", "taxWithheld": "12,751.50" },
  "income5":         { "datePaid": "14 ก.พ. 2568",   "amountPaid": "500,010.01", "taxWithheld": "15,000.30" },
  "income6_note":    "รายได้อื่นๆ",
  "income6":         { "datePaid": "15 มี.ค. 2568",  "amountPaid": "600,060.06", "taxWithheld": "18,001.80" },
  "totals": {
    "totalAmountPaid": "5,741,320.36",
    "totalTaxWithheld": "172,239.60",
    "totalTaxWithheldInWords": "หนึ่งแสนเจ็ดหมื่นสองพันสองร้อยสามสิบเก้าบาทหกสิบสตางค์"
  },
  "otherPayments": {
    "governmentPensionFund": "5,000.00",
    "socialSecurityFund": "750.00",
    "providentFund": "3,000.00"
  },
  "withholdingType": {
    "withholdingTax": true,
    "forever": true,
    "oneTime": true,
    "other": true,
    "otherDetails": "อื่นๆ อื่นๆ อื่นๆ อื่นๆ"
  },
  "certification": { "dateOfIssuance": { "day": "22", "month": "ธันวาคม", "year": "2568" } }
}
```

</details>

---

## Strategy A — multipart/form-data

เหมาะสำหรับ client ที่ upload ไฟล์โดยตรง เช่น web form หรือ mobile app

Best when the client uploads files directly (web form, mobile app, etc.)

**Endpoint:** `POST /api/v1/taxes/multipart`

**Request fields:**

| Field | Type | คำอธิบาย |
|-------|------|-----------|
| `taxInfo` | JSON string | ข้อมูล TaxInfo ใน JSON |
| `signature` | file (PNG) | รูปลายเซ็น |
| `seal` | file (PNG) | รูปตราประทับ |

```bash
curl -X POST http://localhost:8080/api/v1/taxes/multipart \
  -F "taxInfo=$(cat <<'JSON'
{ ...TAX_INFO_JSON... }
JSON
)" \
  -F "signature=@.demo/demo-signature-1280x720-rectangle.png" \
  -F "seal=@.demo/demo-logo-1024x1024-square.png" \
  -o certificate.pdf
```

> ดู TAX_INFO_JSON แบบเต็มด้านบน / See the full TAX_INFO_JSON above.

---

## Strategy B — base64 images in JSON body

เหมาะสำหรับ API client ที่รับส่งแค่ JSON (ไม่ต้องการ multipart)

Best for JSON-only API clients that cannot send multipart.

**Endpoint:** `POST /api/v1/taxes/base64`

**Request body:**

```json
{
  "taxInfo": { ...TAX_INFO_JSON... },
  "signatureBase64": "<base64-encoded PNG>",
  "sealBase64": "<base64-encoded PNG>"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/taxes/base64 \
  -H "Content-Type: application/json" \
  -d "$(python3 - <<'PYEOF'
import json, base64
with open(".demo/demo-signature-1280x720-rectangle.png","rb") as f: sign = base64.b64encode(f.read()).decode()
with open(".demo/demo-logo-1024x1024-square.png","rb") as f: seal = base64.b64encode(f.read()).decode()
tax_info = { ...TAX_INFO_JSON... }
print(json.dumps({"taxInfo": tax_info, "signatureBase64": sign, "sealBase64": seal}))
PYEOF
)" \
  -o certificate.pdf
```

> ดู TAX_INFO_JSON แบบเต็มด้านบน / See the full TAX_INFO_JSON above.
> สคริปต์ `./scripts/demo-rest.sh` ทำทุกขั้นตอนนี้ให้อัตโนมัติ

---

## Strategy C — image URLs (server fetches)

เหมาะสำหรับระบบที่รูปภาพอยู่บน CDN, S3, หรือ storage service อยู่แล้ว — server จะดึงรูปเองจาก URL ที่ส่งมา

Best when images are already hosted on a CDN, S3, or storage service — the server fetches them by URL.

**Endpoint:** `POST /api/v1/taxes/url`

**Request body:**

```json
{
  "taxInfo": { ...TAX_INFO_JSON... },
  "signatureURL": "https://cdn.example.com/signature.png",
  "sealURL": "https://cdn.example.com/logo.png"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/taxes/url \
  -H "Content-Type: application/json" \
  -d '{
    "taxInfo": { ...TAX_INFO_JSON... },
    "signatureURL": "https://cdn.example.com/signature.png",
    "sealURL": "https://cdn.example.com/logo.png"
  }' \
  -o certificate.pdf
```

> ดู TAX_INFO_JSON แบบเต็มด้านบน / See the full TAX_INFO_JSON above.

> ถ้า URL ต้องการ authentication ให้ใช้ `pdf50tawi.LoadImageFromRequest` แทน `LoadImageFromURL` และ set header เองใน handler
>
> For authenticated URLs, use `pdf50tawi.LoadImageFromRequest` and set headers on the request in your own handler.

---

## Response

ทุก endpoint คืน `application/pdf` เมื่อสำเร็จ หรือ JSON error เมื่อเกิดปัญหา

All endpoints return `application/pdf` on success, or a JSON error body on failure.

**Success:** `HTTP 200` + PDF binary stream
```
Content-Disposition: attachment; filename=certificate.pdf
Content-Type: application/pdf
```

**Error:** `HTTP 400` or `HTTP 500`
```json
{ "error": "description of the problem" }
```
