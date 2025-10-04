package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anuchito/pdf50tawi"
)

func main() {
	outputPDF := flag.String("output", "tax50tawi-stamped.pdf", "Output PDF file")
	signature := flag.String("signature", "cmd/demo-cli/demo-signature-1280x720-rectangle.png", "Signature png image file")
	logo := flag.String("logo", "cmd/demo-cli/demo-logo-1024x1024-square.png", "Logo png image file")
	// logo := flag.String("logo", "cmd/demo-cli/demo-signature-1280x720-rectangle.png", "Logo png image file")
	flag.Parse()

	// create output file "tax50tawi-stamped.pdf" as io.Writer
	outputFile, err := os.Create(*outputPDF)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	// load signature image as io.Reader
	signatureReader, err := os.Open(*signature)
	if err != nil {
		log.Fatalf("Error opening signature image: %v", err)
	}
	defer signatureReader.Close()

	// load logo image as io.Reader
	logoReader, err := os.Open(*logo)
	if err != nil {
		log.Fatalf("Error opening logo image: %v", err)
	}
	defer logoReader.Close()

	taxInfo := DemoTaxInfo()
	// demo logo url
	taxInfo.Certification.CompanySealImage = pdf50tawi.Image{
		SourceType: pdf50tawi.SourceTypeURL,
		Value:      "https://raw.githubusercontent.com/AnuchitO/pdf50tawi/main/cmd/demo-cli/demo-logo-1024x1024-square.png",
	}
	if err := pdf50tawi.ValidateTaxInfo(taxInfo); err != nil {
		log.Fatalf("Error validating tax info: %v", err)
	}

	logoReaderx, err := pdf50tawi.LoadImage(taxInfo.Certification.CompanySealImage, nil)
	if err != nil {
		log.Fatalf("Error loading logo image: %v", err)
	}
	defer logoReaderx.Close()

	if err := pdf50tawi.IssueWHTCertificatePDF(outputFile, taxInfo, signatureReader, logoReaderx); err != nil {
		log.Fatalf("Error adding image stamp: %v", err)
	}

	fmt.Println("Successfully processed PDF with Thai text stamp")
}
