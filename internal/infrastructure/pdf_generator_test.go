package infrastructure

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-trial/internal/domain/entity"
	"github.com/shopspring/decimal"
)

func TestPurchaseOrderPDFGenerator_GeneratePO(t *testing.T) {
	generator := NewPurchaseOrderPDFGenerator()

	supplierName := "PT Maju Jaya"
	storeName := "Toko Pusat"
	approvedByName := "Admin User"
	uomCode := "PCS"
	productName := "Barang A"
	address := "Jl. Merdeka No. 123"
	phone := "+628112345678"

	po := &entity.PurchaseOrder{
		PONumber:        "PO/2505/00001",
		OrderDate:       time.Now(),
		Supplier:       entity.Supplier{Name: supplierName},
		Store:          entity.Store{Name: storeName},
		PaymentTermDays: 30,
		PaymentMode:    "TRANSFER",
		TotalAmount:     decimal.NewFromInt(1110000),
		Status:        entity.POStatusApproved,
		ApprovedBy:    &entity.User{Name: approvedByName},
		Items: []entity.PurchaseOrderItem{
			{
				Product:    entity.Product{Name: productName},
				UOM:      entity.UOM{Code: uomCode},
				QtyOrdered: decimal.NewFromInt(100),
				UnitPrice: decimal.NewFromInt(10000),
				Subtotal:  decimal.NewFromInt(1000000),
			},
		},
	}

	po.Supplier.Address = &address
	po.Supplier.PhoneNumber = &phone

	pdfBytes, err := generator.GeneratePO(po)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF bytes")
	}

	if string(pdfBytes[:4]) != "%PDF" {
		t.Error("expected PDF magic bytes %PDF")
	}

	tmpDir := "./test-output"
	os.MkdirAll(tmpDir, 0755)
	tmpFile := filepath.Join(tmpDir, "test_po.pdf")

	if err := os.WriteFile(tmpFile, pdfBytes, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	t.Logf("Generated PDF: %s (%d bytes)", tmpFile, len(pdfBytes))
}

func TestPurchaseOrderPDFGenerator_GeneratePO_WithNotes(t *testing.T) {
	generator := NewPurchaseOrderPDFGenerator()

	notes := "Harap dikirim pada jam kerja"
	po := &entity.PurchaseOrder{
		PONumber:        "PO/2505/00002",
		OrderDate:      time.Now(),
		Supplier:       entity.Supplier{Name: "PT ABC"},
		Store:          entity.Store{Name: "Toko ABC"},
		PaymentTermDays: 30,
		PaymentMode:    "TRANSFER",
		TotalAmount:     decimal.NewFromInt(555000),
		Status:        entity.POStatusApproved,
		Notes:          &notes,
		Items:         []entity.PurchaseOrderItem{},
	}

	pdfBytes, err := generator.GeneratePO(po)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF bytes")
	}
}

func TestPurchaseOrderPDFGenerator_GeneratePO_EmptyItems(t *testing.T) {
	generator := NewPurchaseOrderPDFGenerator()

	po := &entity.PurchaseOrder{
		PONumber:        "PO/2505/00003",
		OrderDate:      time.Now(),
		Supplier:       entity.Supplier{Name: "PT XYZ"},
		Store:          entity.Store{Name: "Toko XYZ"},
		PaymentTermDays: 30,
		PaymentMode:    "TRANSFER",
		TotalAmount:     decimal.NewFromInt(0),
		Status:        entity.POStatusApproved,
		Items:         []entity.PurchaseOrderItem{},
	}

	pdfBytes, err := generator.GeneratePO(po)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF bytes")
	}
}

func TestPurchaseOrderPDFGenerator_GeneratePO_WithDiscount(t *testing.T) {
	generator := NewPurchaseOrderPDFGenerator()

	po := &entity.PurchaseOrder{
		PONumber:        "PO/2505/00004",
		OrderDate:      time.Now(),
		Supplier:       entity.Supplier{Name: "PT Discount"},
		Store:          entity.Store{Name: "Toko Discount"},
		PaymentTermDays: 30,
		PaymentMode:    "TRANSFER",
		TotalAmount:     decimal.NewFromInt(999000),
		Status:        entity.POStatusApproved,
		Items:         []entity.PurchaseOrderItem{},
	}

	pdfBytes, err := generator.GeneratePO(po)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF bytes")
	}
}