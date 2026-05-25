package main

import (
	"fmt"
	"log"
	"time"

	"go-trial/internal/config"
	"go-trial/internal/domain/entity"
	"go-trial/internal/infrastructure/database"

	"github.com/google/uuid"
)

type uomSeed struct {
	Code string
	Name string
}

var masterUOMs = []uomSeed{
	{Code: "PCS", Name: "Pieces / Buah"},
	{Code: "LUSIN", Name: "Lusin"},
	{Code: "KODI", Name: "Kodi"},
	{Code: "GROSS", Name: "Gros"},
	{Code: "PACK", Name: "Pack / Bungkus"},
	{Code: "BOX", Name: "Box / Kotak"},
	{Code: "KARTON", Name: "Karton / Dus"},
	{Code: "SAK", Name: "Sak / Kantong"},
	{Code: "BOTOL", Name: "Botol"},
	{Code: "GRAM", Name: "Gram"},
	{Code: "ONS", Name: "Ons"},
	{Code: "KG", Name: "Kilogram"},
	{Code: "KUINTAL", Name: "Kuintal"},
	{Code: "TON", Name: "Ton"},
	{Code: "ML", Name: "Mililiter"},
	{Code: "CL", Name: "Centiliter"},
	{Code: "DL", Name: "Desiliter"},
	{Code: "LITER", Name: "Liter"},
	{Code: "M3", Name: "Meter Kubik"},
}

func main() {
	cfg := config.Load()
	db := database.NewMySQL(&cfg.Database)

	now := time.Now().UTC()
	var created, skipped int

	for _, u := range masterUOMs {
		var existing entity.UOM
		err := db.Where("code = ?", u.Code).First(&existing).Error
		if err == nil {
			skipped++
			fmt.Printf("  SKIP %s - already exists (ID: %s)\n", u.Code, existing.ID)
			continue
		}

		uom := entity.UOM{
			BaseModel: entity.BaseModel{
				ID:        uuid.Must(uuid.NewV7()),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code: u.Code,
			Name: u.Name,
		}

		if err := db.Create(&uom).Error; err != nil {
			log.Fatalf("Failed to create UOM %s: %v", u.Code, err)
		}
		created++
		fmt.Printf("  OK   %s - %s (ID: %s)\n", u.Code, u.Name, uom.ID)
	}

	fmt.Printf("\nDone! %d created, %d skipped\n", created, skipped)
}
