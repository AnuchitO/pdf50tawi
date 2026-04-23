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
  -F 'taxInfo={
    "payer": {"taxId":"1234567890123","name":"บริษัท ตัวอย่าง จำกัด","address":"123 ถนนสุขุมวิท กรุงเทพฯ"},
    "payee": {"taxId":"3210987654321","name":"นาย ผู้รับเงิน","pnd_3":true},
    "income40_1": {"datePaid":"01 มกราคม 2568","amountPaid":"100,000.00","taxWithheld":"3,000.00"},
    "totals": {"totalAmountPaid":"100,000.00","totalTaxWithheld":"3,000.00","totalTaxWithheldInWords":"สามพันบาทถ้วน"},
    "withholdingType": {"withholdingTax":true},
    "certification": {"dateOfIssuance":{"day":"1","month":"มกราคม","year":"2568"}}
  }' \
  -F "signature=@.demo/demo-signature-1280x720-rectangle.png" \
  -F "seal=@.demo/demo-logo-1024x1024-square.png" \
  -o certificate.pdf
```

---

## Strategy B — base64 images in JSON body

เหมาะสำหรับ API client ที่รับส่งแค่ JSON (ไม่ต้องการ multipart)

Best for JSON-only API clients that cannot send multipart.

**Endpoint:** `POST /api/v1/taxes/base64`

**Request body:**

```json
{
  "taxInfo": { ... },
  "signatureBase64": "<base64-encoded PNG>",
  "sealBase64": "<base64-encoded PNG>"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/taxes/base64 \
  -H "Content-Type: application/json" \
  -d '{
    "taxInfo": {
      "payer": {"taxId":"1234567890123","name":"บริษัท ตัวอย่าง จำกัด","address":"123 ถนนสุขุมวิท กรุงเทพฯ"},
      "payee": {"taxId":"3210987654321","name":"นาย ผู้รับเงิน","pnd_3":true},
      "income40_1": {"datePaid":"01 มกราคม 2568","amountPaid":"100,000.00","taxWithheld":"3,000.00"},
      "totals": {"totalAmountPaid":"100,000.00","totalTaxWithheld":"3,000.00","totalTaxWithheldInWords":"สามพันบาทถ้วน"},
      "withholdingType": {"withholdingTax":true},
      "certification": {"dateOfIssuance":{"day":"1","month":"มกราคม","year":"2568"}}
    },
    "signatureBase64": "'$(base64 < .demo/demo-signature-1280x720-rectangle.png | tr -d '\n')'",
    "sealBase64": "'$(base64 < .demo/demo-logo-1024x1024-square.png | tr -d '\n')'"
  }' \
  -o certificate.pdf
```

---

## Strategy C — image URLs (server fetches)

เหมาะสำหรับระบบที่รูปภาพอยู่บน CDN, S3, หรือ storage service อยู่แล้ว — server จะดึงรูปเองจาก URL ที่ส่งมา

Best when images are already hosted on a CDN, S3, or storage service — the server fetches them by URL.

**Endpoint:** `POST /api/v1/taxes/url`

**Request body:**

```json
{
  "taxInfo": { ... },
  "signatureURL": "https://cdn.example.com/signature.png",
  "sealURL": "https://cdn.example.com/logo.png"
}
```

```bash
curl -X POST http://localhost:8080/api/v1/taxes/url \
  -H "Content-Type: application/json" \
  -d '{
    "taxInfo": {
      "payer": {"taxId":"1234567890123","name":"บริษัท ตัวอย่าง จำกัด","address":"123 ถนนสุขุมวิท กรุงเทพฯ"},
      "payee": {"taxId":"3210987654321","name":"นาย ผู้รับเงิน","pnd_3":true},
      "income40_1": {"datePaid":"01 มกราคม 2568","amountPaid":"100,000.00","taxWithheld":"3,000.00"},
      "totals": {"totalAmountPaid":"100,000.00","totalTaxWithheld":"3,000.00","totalTaxWithheldInWords":"สามพันบาทถ้วน"},
      "withholdingType": {"withholdingTax":true},
      "certification": {"dateOfIssuance":{"day":"1","month":"มกราคม","year":"2568"}}
    },
    "signatureURL": "https://cdn.example.com/signature.png",
    "sealURL": "https://cdn.example.com/logo.png"
  }' \
  -o certificate.pdf
```

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
