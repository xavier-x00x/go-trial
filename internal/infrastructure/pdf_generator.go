package infrastructure

import (
	"go-trial/internal/domain/entity"
	"go-trial/pkg/go-jasperxml"
)

type PDFGenerator interface {
	GeneratePO(po *entity.PurchaseOrder) ([]byte, error)
}

type PurchaseOrderPDFGenerator struct {
	templatePath string
}

func NewPurchaseOrderPDFGenerator() *PurchaseOrderPDFGenerator {
	return &PurchaseOrderPDFGenerator{
		templatePath: "templates/PO_Template.jrxml",
	}
}

func (g *PurchaseOrderPDFGenerator) GeneratePO(po *entity.PurchaseOrder) ([]byte, error) {
	data := convertPurchaseOrderToMap(po)

	j, err := gojasperxml.NewFromJRXMLWithData(g.templatePath, "")
	if err != nil {
		return nil, err
	}

	j.SetData(data)

	return j.GeneratePDF()
}

func convertPurchaseOrderToMap(po *entity.PurchaseOrder) []map[string]interface{} {
	var result []map[string]interface{}

	subtotal := po.TotalAmount.String()
	diskon := "0"
	ppn := "0"
	grandtotal := po.TotalAmount.String()

	supplierName := po.Supplier.Name
	supplierAddr := ""
	if po.Supplier.Address != nil {
		supplierAddr = *po.Supplier.Address
	}
	supplierPhone := ""
	if po.Supplier.PhoneNumber != nil {
		supplierPhone = *po.Supplier.PhoneNumber
	}

	email := ""
	if po.Supplier.Email != nil {
		email = *po.Supplier.Email
	}

	storeName := po.Store.Name
	storeAddress := ""
	if po.Store.Address != nil {
		storeAddress = *po.Store.Address
	} else if po.StoreAddress != nil {
		storeAddress = *po.StoreAddress
	}

	approvedBy := ""
	if po.ApprovedBy != nil {
		approvedBy = po.ApprovedBy.Name
	} else if po.ApprovedByName != nil {
		approvedBy = *po.ApprovedByName
	}

	notes := ""
	if po.Notes != nil {
		notes = *po.Notes
	}
	
	expectedDelivery := ""
	if po.ExpectedDelivery != nil {
		expectedDelivery = po.ExpectedDelivery.Format("02 January 2006")
	}

	for i, item := range po.Items {
		row := map[string]interface{}{
			"no":            po.PONumber,
			"tanggal":       po.OrderDate.Format("02 January 2006"),
			"supplier":      supplierName,
			"alamat":        supplierAddr,
			"telp":          supplierPhone,
			"email":         email,
			"store":         storeName,
			"store_address": storeAddress,
			"expected_delivery": expectedDelivery,
			"namaItem":      item.Product.Name,
			"qty":           item.QtyOrdered.IntPart(),
			"harga":         item.UnitPrice.IntPart(),
			"total":         item.Subtotal.IntPart(),
			"sat":           item.UOM.Code,
			"subtotal":      subtotal,
			"diskon":        diskon,
			"ppn":           ppn,
			"grandtotal":    grandtotal,
			"paymentTerms":  po.PaymentTermDays,
			"paymentMode":   po.PaymentMode,
			"approvedBy":    approvedBy,
			"notes":         notes,
			"rowNum":        i + 1,
		}
		result = append(result, row)
	}

	if len(result) == 0 {
		result = append(result, map[string]interface{}{
			"no":            po.PONumber,
			"tanggal":       po.OrderDate.Format("02 January 2006"),
			"supplier":      supplierName,
			"alamat":        supplierAddr,
			"telp":          supplierPhone,
			"email":         email,
			"store":         storeName,
			"store_address": storeAddress,
			"expected_delivery": expectedDelivery,
			"namaItem":      "",
			"qty":           0,
			"harga":         0,
			"total":         0,
			"sat":           "",
			"subtotal":      subtotal,
			"diskon":        diskon,
			"ppn":           ppn,
			"grandtotal":    grandtotal,
			"paymentTerms":  po.PaymentTermDays,
			"paymentMode":   po.PaymentMode,
			"approvedBy":    approvedBy,
			"notes":         notes,
			"rowNum":        1,
		})
	}

	return result
}
