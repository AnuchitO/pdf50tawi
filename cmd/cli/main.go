package main

// Strategy: CLI — images supplied as local file paths via flags.
//
// Usage:
//
//	go run ./cmd/cli \
//	  --signature path/to/signature.png \
//	  --seal      path/to/seal.png \
//	  --output    certificate.pdf

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AnuchitO/pdf50tawi"
)

func main() {
	outputPath := flag.String("output", "certificate.pdf", "Output PDF file path")
	signPath := flag.String("signature", "", "Signature image file path (PNG)")
	sealPath := flag.String("seal", "", "Company seal image file path (PNG)")
	flag.Parse()

	taxInfo := demoTaxInfo()
	if err := pdf50tawi.ValidateTaxInfo(taxInfo); err != nil {
		log.Fatalf("validation error: %v", err)
	}

	sign := loadOptional(*signPath, "signature")
	seal := loadOptional(*sealPath, "seal")

	out, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("create output: %v", err)
	}
	defer out.Close()

	if err := pdf50tawi.IssueWHTCertificatePDF(out, taxInfo, sign, seal); err != nil {
		log.Fatalf("generate certificate: %v", err)
	}

	fmt.Printf("Certificate written to %s\n", *outputPath)
}

// loadOptional opens a file and returns its reader, or nil if the path is empty.
// Nil is safe — IssueWHTCertificatePDF renders the certificate without the image.
func loadOptional(path, label string) io.Reader {
	if path == "" {
		return nil
	}
	r, err := pdf50tawi.LoadImageFromFile(path)
	if err != nil {
		log.Fatalf("load %s (%s): %v", label, path, err)
	}
	return r
}
