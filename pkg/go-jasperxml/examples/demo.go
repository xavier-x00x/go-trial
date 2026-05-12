//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	gojasperxml "go-trial/pkg/go-jasperxml"
)

func main() {
	jsonData := []byte(`[
		{
			"no": "PO/2505/00001",
			"tanggal": "01 May 2025",
			"supplier": "PT Maju Jaya",
			"alamat": "Jl. Merdeka No. 123",
			"telp": "+628112345678",
			"store": "Toko Pusat",
			"namaItem": "Barang A",
			"qty": 100,
			"harga": 10000,
			"total": 1000000,
			"sat": "PCS",
			"subtotal": 1000000,
			"diskon": 0,
			"ppn": 110000,
			"grandtotal": 1110000,
			"paymentTerms": 30,
			"paymentMode": "TRANSFER",
			"approvedBy": "Admin User"
		}
	]`)

	jx, err := gojasperxml.NewFromJSON(jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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

	err = os.WriteFile("test_output/jasper_po.pdf", pdfBytes, 0644)
	if err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("PDF generated: test_output/jasper_po.pdf")
}
