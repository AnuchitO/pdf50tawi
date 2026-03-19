# pdf50tawi — Thai WHT Certificate PDF Generator

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**Go library สำหรับกรอกข้อมูลและประทับตราลงในแบบฟอร์มหนังสือรับรองการหักภาษี ณ ที่จ่าย (แบบ 50 ทวิ) แบบอัตโนมัติ**

A Go library for automatically filling Thai withholding tax certificate forms (Form 50 Tawi / แบบ 50 ทวิ) with full Thai language support.

---

## ตัวอย่างผลลัพธ์ / Preview

<img src=".demo/tax50tawi-certificated-demo.png" alt="ตัวอย่างแบบฟอร์ม 50 ทวิ" width="700"/>

---

## คุณสมบัติ / Features

| คุณสมบัติ | รายละเอียด |
|-----------|------------|
| รองรับภาษาไทย 100% | ใช้ฟอนต์ TH Sarabun New แบบ subsetted ไม่มีปัญหาภาษาต่างดาว |
| ขนาดไฟล์เล็ก | output ประมาณ 150–500 KB (จากเดิม 2.5 MB) |
| กรอกข้อมูลครบทุกช่อง | รองรับทุก field ในแบบฟอร์ม 50 ทวิ |
| ประทับลายเซ็น & โลโก้ | รองรับรูปภาพ PNG/JPEG |
| ใช้งานผ่าน CLI และ Go API | ยืดหยุ่น เหมาะทั้ง standalone และ integration |

---

## การติดตั้ง / Installation

```bash
go get github.com/anuchito/pdf50tawi
```

---

## การใช้งาน / Usage

### Go API

```go
import "github.com/anuchito/pdf50tawi"

// กรอกข้อมูลภาษี
taxInfo := pdf50tawi.TaxInfo{
    Payer: pdf50tawi.Payer{
        Name:  "บริษัท ตัวอย่าง จำกัด",
        TaxID: "1234567890123",
    },
    Payee: pdf50tawi.Payee{
        Name:  "นาย ผู้รับเงิน",
        TaxID: "9876543210987",
    },
    // ... ข้อมูลอื่น ๆ
}

// โหลดลายเซ็นและโลโก้
sign, _ := os.Open("signature.png")
seal, _ := os.Open("logo.png")

// สร้างไฟล์ PDF
out, _ := os.Create("tax50tawi.pdf")
pdf50tawi.IssueWHTCertificatePDF(out, taxInfo, sign, seal)
```

### CLI

```bash
# build
go build ./cmd/demo-cli/

# รันพร้อมรูปลายเซ็นและโลโก้
go run ./cmd/demo-cli/ \
  -signature path/to/signature.png \
  -seal path/to/logo.png \
  -output tax50tawi.pdf
```

### Makefile shortcuts

```bash
make run              # รันด้วยข้อมูลตัวอย่าง (ไม่มีรูปภาพ)
make run-with-args    # รันพร้อมลายเซ็นและโลโก้ตัวอย่าง
```

---

## ข้อกำหนดขนาดรูปภาพ / Image Requirements

รูปภาพต้องเป็น **PNG พื้นหลังโปร่งใส** และมีขนาดตามนี้:

| ประเภท | ขนาด |
|--------|------|
| ลายเซ็น (Signature) | 1280 × 720 px |
| โลโก้ / ตราประทับ สี่เหลี่ยมจัตุรัสหรือวงกลม | 1024 × 1024 px |
| โลโก้ / ตราประทับ สี่เหลี่ยมผืนผ้า | 1280 × 720 px |

---

## ขนาดไฟล์ผลลัพธ์ / Output Size

| สถานการณ์ | ขนาดไฟล์ |
|-----------|----------|
| ข้อมูลข้อความ (ไม่มีรูป) | ~150 KB |
| พร้อมลายเซ็น + โลโก้ | ~400–500 KB |
| เวอร์ชันเก่า (pdfcpu watermarks) | ~2,500 KB |

---

## ไลบรารีแนะนำ / Related Libraries

- [bahttext](https://github.com/anuchito/bahttext) — แปลงตัวเลขเป็นตัวอักษรภาษาไทย (บาท)
- [currency-formatter](https://github.com/anuchito/currency-formatter) — จัดรูปแบบตัวเลขเงินบาท
- [date-thai-formatter](https://github.com/anuchito/date-thai-formatter) — แปลงวันที่เป็นภาษาไทย (พ.ศ., ชื่อเดือนเต็ม/ย่อ)

---

## ปัญหาที่พบ / Known Issues

- ตัวอักษร `%` ในช่อง `Income40_4B_1_4_Rate` อาจแสดงผลไม่ถูกต้อง

---

## License

MIT
