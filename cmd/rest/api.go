package main

// REST API — three strategies for supplying signature and seal images.
//
// Strategy A  POST /api/v1/taxes/multipart  multipart/form-data upload
// Strategy B  POST /api/v1/taxes/base64     JSON body with base64-encoded images
// Strategy C  POST /api/v1/taxes/url        JSON body with image URLs (server fetches)

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	e.POST("/api/v1/taxes/multipart", handleMultipart)
	e.POST("/api/v1/taxes/base64", handleBase64)
	e.POST("/api/v1/taxes/url", handleURL)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(e.Start(":" + port))
}

// ── Strategy A: multipart/form-data ──────────────────────────────────────────
//
// curl -X POST http://localhost:8080/api/v1/taxes/multipart \
//   -F 'taxInfo={"payer":{"taxId":"1234567890123","name":"บริษัท ตัวอย่าง จำกัด","address":"123 ถนน"},...}' \
//   -F 'signature=@signature.png' \
//   -F 'seal=@seal.png' \
//   -o certificate.pdf
func handleMultipart(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("parse multipart form: "+err.Error()))
	}
	defer form.RemoveAll()

	taxInfo, err := parseTaxInfoFromForm(form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp(err.Error()))
	}
	if err := pdf50tawi.ValidateTaxInfo(taxInfo); err != nil {
		return c.JSON(http.StatusBadRequest, errResp(err.Error()))
	}

	sign, err := readFormFile(form, "signature")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("signature: "+err.Error()))
	}
	seal, err := readFormFile(form, "seal")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("seal: "+err.Error()))
	}

	return streamCertificate(c, taxInfo, sign, seal)
}

// ── Strategy B: JSON with base64-encoded images ───────────────────────────────
//
// curl -X POST http://localhost:8080/api/v1/taxes/base64 \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "taxInfo": {"payer": {...}, ...},
//     "signatureBase64": "<base64-encoded PNG>",
//     "sealBase64": "<base64-encoded PNG>"
//   }' \
//   -o certificate.pdf
type base64Request struct {
	TaxInfo         pdf50tawi.TaxInfo `json:"taxInfo"`
	SignatureBase64 string             `json:"signatureBase64"`
	SealBase64      string             `json:"sealBase64"`
}

func handleBase64(c echo.Context) error {
	var req base64Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errResp(err.Error()))
	}
	if err := pdf50tawi.ValidateTaxInfo(req.TaxInfo); err != nil {
		return c.JSON(http.StatusBadRequest, errResp(err.Error()))
	}

	signData, err := base64.StdEncoding.DecodeString(req.SignatureBase64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("invalid signatureBase64: "+err.Error()))
	}
	sealData, err := base64.StdEncoding.DecodeString(req.SealBase64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("invalid sealBase64: "+err.Error()))
	}

	return streamCertificate(c, req.TaxInfo, bytes.NewReader(signData), bytes.NewReader(sealData))
}

// ── Strategy C: JSON with image URLs (server fetches) ────────────────────────
//
// curl -X POST http://localhost:8080/api/v1/taxes/url \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "taxInfo": {"payer": {...}, ...},
//     "signatureURL": "https://storage.example.com/signatures/company-sign.png",
//     "sealURL": "https://storage.example.com/seals/company-seal.png"
//   }' \
//   -o certificate.pdf
type urlRequest struct {
	TaxInfo      pdf50tawi.TaxInfo `json:"taxInfo"`
	SignatureURL string             `json:"signatureURL"`
	SealURL      string             `json:"sealURL"`
}

func handleURL(c echo.Context) error {
	var req urlRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errResp(err.Error()))
	}
	if err := pdf50tawi.ValidateTaxInfo(req.TaxInfo); err != nil {
		return c.JSON(http.StatusBadRequest, errResp(err.Error()))
	}

	sign, err := pdf50tawi.LoadImageFromURL(req.SignatureURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("signatureURL: "+err.Error()))
	}
	seal, err := pdf50tawi.LoadImageFromURL(req.SealURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errResp("sealURL: "+err.Error()))
	}

	return streamCertificate(c, req.TaxInfo, sign, seal)
}

// ── Shared helpers ────────────────────────────────────────────────────────────

func streamCertificate(c echo.Context, taxInfo pdf50tawi.TaxInfo, sign, seal io.Reader) error {
	var buf bytes.Buffer
	if err := pdf50tawi.IssueWHTCertificatePDF(&buf, taxInfo, sign, seal); err != nil {
		return c.JSON(http.StatusInternalServerError, errResp("generate certificate: "+err.Error()))
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=certificate.pdf")
	return c.Stream(http.StatusOK, "application/pdf", &buf)
}

func parseTaxInfoFromForm(form *multipart.Form) (pdf50tawi.TaxInfo, error) {
	values, ok := form.Value["taxInfo"]
	if !ok || len(values) == 0 {
		return pdf50tawi.TaxInfo{}, errors.New("missing 'taxInfo' form field")
	}
	var taxInfo pdf50tawi.TaxInfo
	if err := json.Unmarshal([]byte(values[0]), &taxInfo); err != nil {
		return pdf50tawi.TaxInfo{}, fmt.Errorf("invalid taxInfo JSON: %w", err)
	}
	return taxInfo, nil
}

func readFormFile(form *multipart.Form, field string) (io.Reader, error) {
	files, ok := form.File[field]
	if !ok || len(files) == 0 {
		return nil, fmt.Errorf("missing file field '%s'", field)
	}
	f, err := files[0].Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, err
	}
	return &buf, nil
}

func errResp(msg string) map[string]string {
	return map[string]string{"error": msg}
}
