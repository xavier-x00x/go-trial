package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductData struct {
	SKU        string
	Name       string
	UOM        string
	Barcode    string
	Markup    string
	BuyPrice   string
	Supplier  string
	MaxStock  string
	AvgSales  string
	Length    string
	Width     string
	Height    string
}

type CreateProductRequest struct {
	SKU         string          `json:"sku"`
	Barcode    *string         `json:"barcode,omitempty"`
	Name       string          `json:"name"`
	CategoryID *uuid.UUID      `json:"category_id,omitempty"`
	BaseUOMID  uuid.UUID       `json:"base_uom_id"`
	IsStockable bool           `json:"is_stockable"`
	Length     decimal.Decimal `json:"length"`
	Width      decimal.Decimal `json:"width"`
	Height     decimal.Decimal `json:"height"`
	Weight     decimal.Decimal `json:"weight"`
	IsStackable bool          `json:"is_stackable"`
}

type ProposalRequest struct {
	EntityType  string `json:"entity_type"`
	ActionType  string `json:"action_type"`
	PayloadJSON string `json:"payload_json"`
	Reason      string `json:"reason"`
}

var baseUOMID = uuid.MustParse("019dd000-c945-7669-a461-1250bb5fc786") // PCS

func main() {
	// Login first
	token := login()
	if token == "" {
		log.Fatal("Failed to login")
	}
	log.Printf("Logged in with token: %s...", token[:50])

	// Parse CSV file
	products, err := parseCSV("MASTER_PRODUCT.txt")
	if err != nil {
		log.Fatalf("Error parsing CSV: %v", err)
	}

	log.Printf("Parsed %d products", len(products))

	// Create proposals
	createProducts(token, products)
}

func login() string {
	loginReq := map[string]string{
		"identity": "dwi@example.com",
		"password": "12345678",
	}

	jsonData, _ := json.Marshal(loginReq)
	resp, err := http.Post("http://localhost:3050/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		log.Printf("Login error: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if data, ok := result["data"].(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok {
			return token
		}
	}
	return ""
}

func parseCSV(filename string) ([]ProductData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var products []ProductData
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if len(record) < 12 {
			continue
		}

		product := ProductData{
			SKU:       strings.Trim(record[0], "\""),
			Name:      strings.Trim(record[1], "\""),
			UOM:       strings.Trim(record[2], "\""),
			Barcode:   strings.Trim(record[3], "\""),
			Markup:    strings.Trim(record[4], "\""),
			BuyPrice:  strings.Trim(record[5], "\""),
			Supplier: strings.Trim(record[6], "\""),
			MaxStock:  strings.Trim(record[7], "\""),
			AvgSales:  strings.Trim(record[8], "\""),
			Length:    strings.Trim(record[9], "\""),
			Width:     strings.Trim(record[10], "\""),
			Height:   strings.Trim(record[11], "\""),
		}
		products = append(products, product)
	}

	return products, nil
}

func createProducts(token string, products []ProductData) {
	client := &http.Client{Timeout: 30 * time.Second}

	count := 0
	for i, p := range products {
		if i >= 10 { // Limit to 10 for testing
			break
		}

		// Prepare the product request
		req := CreateProductRequest{
			SKU:        p.SKU,
			Name:       p.Name,
			BaseUOMID:  baseUOMID,
			IsStockable: true,
		}

		if p.Barcode != "" {
			req.Barcode = &p.Barcode
		}

		length, _ := strconv.ParseFloat(p.Length, 64)
		width, _ := strconv.ParseFloat(p.Width, 64)
		height, _ := strconv.ParseFloat(p.Height, 64)

		if length > 0 {
			req.Length = decimal.NewFromFloat(length)
		}
		if width > 0 {
			req.Width = decimal.NewFromFloat(width)
		}
		if height > 0 {
			req.Height = decimal.NewFromFloat(height)
		}

		payloadJSON, err := json.Marshal(req)
		if err != nil {
			log.Printf("Error marshalling product %s: %v", p.SKU, err)
			continue
		}

		proposalReq := ProposalRequest{
			EntityType:  "PRODUCT",
			ActionType:  "CREATE",
			PayloadJSON: string(payloadJSON),
			Reason:     fmt.Sprintf("Bulk import: %s - %s", p.SKU, p.Name),
		}

		jsonData, _ := json.Marshal(proposalReq)

		httpReq, _ := http.NewRequest("POST", "http://localhost:3050/api/master-data/proposals", strings.NewReader(string(jsonData)))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(httpReq)
		if err != nil {
			log.Printf("Error creating proposal for %s: %v", p.SKU, err)
			continue
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			count++
			log.Printf("Created proposal %d: %s - %s", i+1, p.SKU, p.Name)
		} else {
			log.Printf("Failed to create proposal for %s: status=%d, result=%v", p.SKU, resp.StatusCode, result)
		}
	}

	log.Printf("Total created: %d/%d", count, len(products))
}