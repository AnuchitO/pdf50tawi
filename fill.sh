#!/bin/bash

curl -X POST http://localhost:8080/api/v1/taxes \
  -F 'taxInfo={
  "documentDetails": {
    "bookNumber": "001",
    "documentNumber": "001"
  },
  "payer": {
    "taxId": "1234567890123",
    "name": "บริษัท ตัวอย่าง จำกัด",
    "address": "123 ถนนตัวอย่าง แขวงตัวอย่าง เขตตัวอย่าง กรุงเทพฯ 10110"
  },
  "payee": {
    "taxId": "9876543210987",
    "name": "นายตัวอย่าง ตัวอย่าง",
    "address": "456 ถนนตัวอย่าง2 แขวงตัวอย่าง2 เขตตัวอย่าง2 กรุงเทพฯ 10220", 
    "sequenceNumber": "1",
    "pnd_1a": false,
    "pnd_1aSpecial": false,
    "pnd_2": false,
    "pnd_3": true,
    "pnd_2a": false,
    "pnd_3a": false,
    "pnd_53": false
  },
  "income40_1": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_2": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_3": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4A": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_1_1": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_1_2": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_1_3": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_1_4": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_2_1": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_2_2": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_2_3": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_2_4": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income40_4B_2_5": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income5": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income6": {
    "datePaid": "2024-01-15",
    "amountPaid": "100000",
    "taxWithheld": "5000"
  },
  "income6_note": "note income6",
  "totals": {
    "totalAmountPaid": "100000",
    "totalTaxWithheld": "5000",
    "totalTaxWithheldInWords": "5000"
  },
  "totalsInWords": "",
  "otherPayments": {
    "governmentPensionFund": "100000",
    "socialSecurityFund": "5000",
    "providentFund": "5000"
  },
  "withholdingType": {
    "withholdingTax": true,
    "forever": false,
    "oneTime": false,
    "other": false,
    "otherDetails": ""
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