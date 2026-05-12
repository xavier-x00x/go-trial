//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	gojasperxml "go-trial/pkg/go-jasperxml"
)

func main() {
	jx, err := gojasperxml.NewFromJRXML("templates/barcode_test.jrxml")
	if err != nil {
		fmt.Printf("Error loading JRXML: %v\n", err)
		os.Exit(1)
	}

	jx.SetData([]map[string]interface{}{
		{
			"barcode_data": "ABC-123-XYZ",
			"qrcode_data":  "https://github.com/jung-kurt/gofpdf",
			"ean13_data":   "123456789012",
		},
	})

	pdfBytes, err := jx.GeneratePDF()
	if err != nil {
		fmt.Printf("Error generating PDF: %v\n", err)
		os.Exit(1)
	}

	err = os.MkdirAll("test_output", 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile("test_output/barcode_test.pdf", pdfBytes, 0644)
	if err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Barcode PDF generated: test_output/barcode_test.pdf")
}
