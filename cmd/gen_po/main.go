package main

import (
	"fmt"
	"log"

	gojasperxml "go-trial/pkg/go-jasperxml"
)

func main() {
	templatePath := "templates/PO_Template.jrxml"
	dataPath := "test_output/po_test_data.json"

	j, err := gojasperxml.NewFromJRXMLWithData(templatePath, dataPath)
	if err != nil {
		log.Fatalf("Failed to create JasperXML: %v", err)
	}

	err = j.SavePDF("test_output/purchase_order_from_jrxml.pdf")
	if err != nil {
		log.Fatalf("Failed to generate PDF: %v", err)
	}

	fmt.Println("PDF generated from JRXML template: test_output/purchase_order_from_jrxml.pdf")
}