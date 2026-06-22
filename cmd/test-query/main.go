package main

import (
	"context"
	"fmt"

	"go-trial/internal/query/params"
	"go-trial/internal/query/service"
	"go-trial/internal/config"
	"go-trial/internal/registry"
	"go-trial/internal/infrastructure/database"
)

func main() {
	cfg := config.Load()
	db := database.NewMySQL(&cfg.Database)
	rdb := database.NewRedis(&cfg.Redis)
	reg := registry.New(db, rdb, cfg)

	qService := service.NewProductQueryService(reg.DB.Debug())
	
	meta := &params.MetaRequest{
		Page: 1,
		Limit: 10,
		OrderColumn: "id",
		OrderDir: "asc",
	}

	data, metaResp, err := qService.GetListPagination(context.Background(), meta)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("DATA LENGTH: %d\n", len(data))
		fmt.Printf("META: %+v\n", metaResp)
		if len(data) > 0 {
			fmt.Printf("FIRST ROW: %+v\n", data[0])
		}
	}
}
