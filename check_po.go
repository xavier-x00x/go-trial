package main

import (
	"fmt"
	"go-trial/internal/domain/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/erp_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var pos []entity.PurchaseOrder
	db.Preload("Items").Order("created_at desc").Limit(5).Find(&pos)

	for _, po := range pos {
		fmt.Printf("PO Number: %s, Items Count: %d\n", po.PONumber, len(po.Items))
	}
}
