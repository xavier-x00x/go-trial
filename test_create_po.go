package main

import (
	"context"
	"fmt"
	"go-trial/internal/infrastructure/repository"
	"go-trial/internal/infrastructure/uow"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func main() {
	dsn := "root:rootpassword@tcp(mysql-db:3306)/noto?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	repo := repository.NewPurchaseOrderRepository(db)
	uowObj := uow.NewUnitOfWork(db)
	ctx := context.Background()

	txCtx, err := uowObj.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer uowObj.Rollback(txCtx)

	po := &entity.PurchaseOrder{
		PONumber:      "TEST-PO-001",
		SupplierID:    uuid.New(),
		StoreID:       uuid.New(),
		WarehouseID:   uuid.New(),
		OrderDate:     time.Now(),
		Status:        "DRAFT",
		CreatedByID:   uuid.New(),
	}
	po.GenerateID()

	items := []entity.PurchaseOrderItem{
		{
			PurchaseOrderID: po.ID,
			SeqNo:           1,
			ProductID:       uuid.New(),
			UOMID:           uuid.New(),
			QtyOrdered:      decimal.NewFromInt(10),
			UnitPrice:       decimal.NewFromInt(100),
			Subtotal:        decimal.NewFromInt(1000),
		},
	}
	po.Items = items

	err = repo.Create(txCtx, po)
	if err != nil {
		fmt.Println("Error creating PO:", err)
		return
	}
	uowObj.Commit(txCtx)

	// Fetch it back
	savedPO, err := repo.FindByID(ctx, po.ID.String())
	if err != nil {
		fmt.Println("Error finding PO:", err)
		return
	}
	fmt.Printf("Saved PO items count: %d\n", len(savedPO.Items))
}
