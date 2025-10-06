package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/anuchito/pdf50tawi"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// POST /api/v1/taxes - Generate tax certificate with uploaded files
	e.POST("/api/v1/taxes", handleTaxCertificate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(e.Start(":" + port))
}

// handleTaxCertificate handles POST /api/v1/taxes with JSON payload
func handleTaxCertificate(c echo.Context) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse multipart form: " + err.Error(),
		})
	}
	defer form.RemoveAll()

	// Parse tax information from form fields (for now)
	taxInfo, err := parseTaxInfoFromForm(form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse tax information: " + err.Error(),
		})
	}

	// Validate tax information
	if err := pdf50tawi.ValidateTaxInfo(taxInfo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid tax information: " + err.Error(),
		})
	}

	// read signatureImage into buffer
	files := form.File["signatureImage"] // or "companySeal"
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No signature image uploaded",
		})
	}

	file, err := files[0].Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open signature image: " + err.Error(),
		})
	}
	defer file.Close()

	// Read file into buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read file",
		})
	}

	// read sealImage into buffer
	files = form.File["companySeal"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No company seal image uploaded",
		})
	}

	file, err = files[0].Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open company seal image: " + err.Error(),
		})
	}
	defer file.Close()

	// Read file into buffer
	var buf2 bytes.Buffer
	if _, err := io.Copy(&buf2, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read file",
		})
	}

	// TODO: resolve images from multipart form upload

	sign := io.NopCloser(&buf)
	seal := io.NopCloser(&buf2)

	out := c.Response().Writer
	if err := pdf50tawi.IssueWHTCertificatePDF(out, taxInfo, sign, seal); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate PDF: " + err.Error(),
		})
	}

	// Set response headers for PDF
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=certificate.pdf")

	return nil
}

// parseTaxInfoFromForm parses complete tax information from multipart form data
func parseTaxInfoFromForm(form *multipart.Form) (pdf50tawi.TaxInfo, error) {
	data := form.Value["taxInfo"]
	if len(data) == 0 {
		return pdf50tawi.TaxInfo{}, errors.New("missing 'taxInfo' field")
	}

	taxInfo := pdf50tawi.TaxInfo{}
	if err := json.Unmarshal([]byte(data[0]), &taxInfo); err != nil {
		return pdf50tawi.TaxInfo{}, err
	}

	return taxInfo, nil
}

// processUploadedFiles processes uploaded signature and seal files
func processUploadedFiles(c echo.Context, form *multipart.Form, taxInfo *pdf50tawi.TaxInfo) error {
	files := form.File

	// Process signature file
	if signatureFiles, ok := files["signature"]; ok && len(signatureFiles) > 0 {
		signatureFile := signatureFiles[0]

		// For demo purposes, we'll use the uploaded file directly
		// In a real implementation, you might want to save it temporarily or validate it
		taxInfo.Certification.PayerSignatureImage = pdf50tawi.Image{
			SourceType: pdf50tawi.Upload,
			Value:      signatureFile.Filename, // This would need to be handled differently in real implementation
		}
	}

	// Process seal file
	if sealFiles, ok := files["seal"]; ok && len(sealFiles) > 0 {
		sealFile := sealFiles[0]

		taxInfo.Certification.CompanySealImage = pdf50tawi.Image{
			SourceType: pdf50tawi.Upload,
			Value:      sealFile.Filename, // This would need to be handled differently in real implementation
		}
	}

	return nil
}

/*
Usage example with JSON payload:

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
    "address": "456 ถนนตัวอย่าง2 แขวงตัวอย่าง2 เขตตัวอย่าง2 กรุงเทพฯ 10220"
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
  -F "signatureImage=@cmd/demo-cli/demo-logo-1280x720-rectangle.png" \
  -F "companySeal=@cmd/demo-cli/demo-logo-1280x720-rectangle.png" \
  -o taxes.pdf

The API accepts:
- data: JSON string with complete TaxInfo structure (including all payer, payee, income, etc. fields)
- signatureImage: Signature image file
- companySeal: Company seal image file

It returns a PDF tax certificate as response.
*/
