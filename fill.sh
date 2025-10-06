#!/bin/bash

curl -X POST http://localhost:8080/api/v1/taxes \
  -F 'taxInfo={
  "documentDetails": {
      "bookNumber": "001",
      "documentNumber": "2568-001"
  },
  "payer": {
      "taxId": "1234567890123",
      "taxId10Digit": "0987654321",
      "name": "บริษัท ตัวอย่าง จำกัด",
      "address": "123 ถนนสุขุมวิท แขวงคลองตัน เขตวัฒนา กรุงเทพฯ 10110"
  },
  "payee": {
      "taxId": "3210987654321",
      "taxId10Digit": "1234567890",
      "name": "นางสาวสมชาย นามสกุลยาวมากไหมนะก็ไม่รู้เหมือนกัน",
      "address": "555 ต.ทุ่งนา  อ.ทุ่งนา  จ.ชลบุรี  12345",
      "sequenceNumber": "321",
      "pnd_1a": false,
      "pnd_1aSpecial": false,
      "pnd_2": false,
      "pnd_3": true,
      "pnd_2a": false,
      "pnd_3a": false,
      "pnd_53": false
  },
  "income40_1": {
      "datePaid": "01 มกราคม 2568",
      "amountPaid": "401,010.01",
      "taxWithheld": "12,030.30"
  },
  "income40_2": {
      "datePaid": "02 ก.พ. 2568",
      "amountPaid": "402,020.02",
      "taxWithheld": "12,060.60"
  },
  "income40_3": {
      "datePaid": "03/03/2568",
      "amountPaid": "403,030.03",
      "taxWithheld": "12,090.90"
  },
  "income40_4A": {
      "datePaid": "04/04/2568",
      "amountPaid": "404,040.04",
      "taxWithheld": "12,121.20"
  },
  "income40_4B_1_1": {
      "datePaid": "11/04/2568",
      "amountPaid": "404,101.05",
      "taxWithheld": "12,123.03"
  },
  "income40_4B_1_2": {
      "datePaid": "12/04/2568",
      "amountPaid": "404,102.06",
      "taxWithheld": "12,123.06"
  },
  "income40_4B_1_3": {
      "datePaid": "13/04/2568",
      "amountPaid": "404,103.07",
      "taxWithheld": "12,123.09"
  },
  "income40_4B_1_4_rate": "ร้อยละ 10",
  "income40_4B_1_4": {
      "datePaid": "14/04/2568",
      "amountPaid": "404,104.08",
      "taxWithheld": "12,123.12"
  },
  "income40_4B_2_1": {
      "datePaid": "21/04/2568",
      "amountPaid": "404,201.09",
      "taxWithheld": "12,126.03"
  },
  "income40_4B_2_2": {
      "datePaid": "22/04/2568",
      "amountPaid": "404,202.10",
      "taxWithheld": "12,126.06"
  },
  "income40_4B_2_3": {
      "datePaid": "23/04/2568",
      "amountPaid": "404,203.11",
      "taxWithheld": "12,126.09"
  },
  "income40_4B_2_4": {
      "datePaid": "24/04/2568",
      "amountPaid": "404,204.12",
      "taxWithheld": "12,126.12"
  },
  "income40_4B_2_5_note": "ระบุ ใดๆๆๆๆๆ ",
  "income40_4B_2_5": {
      "datePaid": "25/04/2568",
      "amountPaid": "404,205.13",
      "taxWithheld": "12,126.15"
  },
  "income5": {
      "datePaid": "05/05/2568",
      "amountPaid": "500,555.55",
      "taxWithheld": "15,555.55"
  },
  "income6": {
      "datePaid": "06/06/2568",
      "amountPaid": "600,666.66",
      "taxWithheld": "16,666.66"
  },
  "income6_note": "อื่นๆ (อื่นๆ) อื่นๆ อื่นๆ อื่นๆ",
  "totals": {
      "totalAmountPaid": "5,847,066.90",
      "totalTaxWithheld": "175,411.22",
      "totalTaxWithheldInWords": "หนึ่งแสนเจ็ดหมื่นห้าพันสี่ร้อยสิบเอ็ดบาทยี่สิบสองสตางค์"
  },
  "otherPayments": {
      "governmentPensionFund": "22,222.22",
      "socialSecurityFund": "33,333.33",
      "providentFund": "44,444.44"
  },
  "withholdingType": {
      "withholdingTax": true,
      "forever": false,
      "oneTime": false,
      "other": false,
      "otherDetails": "อื่นๆ (อื่นๆ) อื่นๆ อื่นๆ อื่นๆ"
  },
  "certification": {
    "payerSignatureImage": {
      "sourceType": "upload",
      "value": "signatureImage"
    },
    "companySealImage": {
      "sourceType": "upload",
      "value": "companySeal"
    },
    "dateOfIssuance": {
      "day": "26",
      "month": "09",
      "year": "2025"
    }
  }
}' \
  -F "signatureImage=@cmd/demo-cli/demo-signature-1280x720-rectangle.png" \
  -F "companySeal=@cmd/demo-cli/demo-logo-1024x1024-square.png" \
  --output output.pdf