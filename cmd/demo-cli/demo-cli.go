package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anuchito/pdf50tawi"
)

func main() {
	demoUsingTaxInfo()
	// demoPassSignatureAndLogoUsingFlag()
}

func demoUsingTaxInfo() {
	taxInfo := demoTaxInfo()

	if err := pdf50tawi.ValidateTaxInfo(taxInfo); err != nil {
		log.Fatalf("Error validating tax info: %v\n", err)
	}

	sign, err := pdf50tawi.LoadImage(taxInfo.Certification.PayerSignatureImage)
	if err != nil {
		log.Fatalf("load signature image fail: %s\n", err)
	}
	defer sign.Close()

	// demo logo url
	// taxInfo.Certification.CompanySealImage = pdf50tawi.Image{
	// 	SourceType: pdf50tawi.URL,
	// 	Value:      "https://raw.githubusercontent.com/AnuchitO/pdf50tawi/main/cmd/demo-cli/demo-logo-1024x1024-square.png",
	// }

	seal, err := pdf50tawi.LoadImage(taxInfo.Certification.CompanySealImage)
	if err != nil {
		log.Fatalf("load company seal image fail: %s\n", err)
	}
	defer seal.Close()

	outputFile, err := os.Create("tax50tawi-stamped.pdf")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	if err := pdf50tawi.IssueWHTCertificatePDF(outputFile, taxInfo, sign, seal); err != nil {
		log.Fatalf("Error adding image stamp: %v", err)
	}

	fmt.Println("Successfully processed PDF with Thai text stamp look at tax50tawi-stamped.pdf")
}

func demoPassSignatureAndLogoUsingFlag() {
	outputPDF := flag.String("output", "tax50tawi-stamped.pdf", "Output PDF file")
	signPath := flag.String("signature", "cmd/demo-cli/demo-signature-1280x720-rectangle.png", "Signature png image file")
	sealPath := flag.String("seal", "cmd/demo-cli/demo-logo-1024x1024-square.png", "Logo png image file")
	// logo := flag.String("logo", "cmd/demo-cli/demo-signature-1280x720-rectangle.png", "Logo png image file")
	flag.Parse()

	// create output file "tax50tawi-stamped.pdf" as io.Writer
	outputFile, err := os.Create(*outputPDF)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	// load signature image as io.Reader
	sign, err := os.Open(*signPath)
	if err != nil {
		log.Fatalf("Error opening signature image: %v", err)
	}
	defer sign.Close()

	// load logo image as io.Reader
	companySeal, err := os.Open(*sealPath)
	if err != nil {
		log.Fatalf("Error opening logo image: %v", err)
	}
	defer companySeal.Close()

	taxInfo := demoTaxInfo()

	if err := pdf50tawi.ValidateTaxInfo(taxInfo); err != nil {
		log.Fatalf("Error validating tax info: %v", err)
	}

	if err := pdf50tawi.IssueWHTCertificatePDF(outputFile, taxInfo, sign, companySeal); err != nil {
		log.Fatalf("Error adding image stamp: %v", err)
	}

	fmt.Println("Successfully processed PDF with Thai text stamp look at " + *outputPDF)
}
