# API Documentation - Go Trial

## Overview

Sistem manajemen pembelian dan inventori yang dibangun dengan Go menggunakan framework Fiber. API ini menyediakan endpoint untuk mengelola produk, supplier, pembelian, gudang, dan akun keuangan.

**Base URL:** `http://localhost:8080/api`

**Framework:** Fiber v2  
**Database:** PostgreSQL  
**Authentication:** JWT (Bearer Token)

---

## Table of Contents

1. [Authentication](#authentication)
2. [Stores](#stores)
3. [Products](#products)
4. [Suppliers](#suppliers)
5. [Finance](#finance)
6. [Warehouses](#warehouses)
7. [Master Data](#master-data)
8. [Purchase Orders](#purchase-orders)
9. [Goods Receipts](#goods-receipts)
10. [Purchase Payments](#purchase-payments)
11. [Planning](#planning)
12. [Roles & Permissions](#roles--permissions)

---

## Grouping

### Data Master

- [Master Data Proposals](#master-data)
- [Stores](#stores)
- [Products](#products)
- [Suppliers](#suppliers)
- [Warehouses](#warehouses)

### Finance

- [Chart of Accounts](#finance)
- [Customers](#finance)
- [Payment Methods](#finance)
- [Price Lists](#finance)
- [Taxes](#finance)

### Transactions

- [Purchase Orders](#purchase-orders)
- [Goods Receipts](#goods-receipts)
- [Purchase Payments](#purchase-payments)
- [Planning](#planning)

---

## Authentication

### Register

Create a new user account.

```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2026-05-13T10:00:00Z"
  }
}
```

---

### Login

Authenticate user and get JWT token.

```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

---

### Refresh Token

Get a new access token using refresh token.

```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

---

### Logout

Invalidate the current session.

```http
POST /auth/logout
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Logout successful"
}
```

---

### Get Current User

Get authenticated user information.

```http
GET /auth/me
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2026-05-13T10:00:00Z"
  }
}
```

---

## Stores

### Create Store

Create a new store.

```http
POST /stores
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Store 001",
  "code": "STR001",
  "address": "Jl. Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "phone": "021123456",
  "email": "store@example.com",
  "is_active": true
}
```

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Store created successfully",
  "data": {
    "id": "uuid",
    "name": "Store 001",
    "code": "STR001",
    "address": "Jl. Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "phone": "021123456",
    "email": "store@example.com",
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

---

### Get All Stores

Get all stores.

```http
GET /stores
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Stores retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "Store 001",
      "code": "STR001",
      "is_active": true,
      "created_at": "2026-05-13T10:00:00Z"
    }
  ]
}
```

---

### Get Stores with Pagination

Get stores with pagination.

```http
GET /stores/pagination?page=1&limit=10&search=&order_column=id&order_dir=asc
Authorization: Bearer {access_token}
```

**Query Parameters:**

- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 10) - Items per page
- `search` (string) - Search term
- `order_column` (string, default: id) - Column to sort by
- `order_dir` (string, default: asc) - Sort direction (asc/desc)

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Stores retrieved successfully",
  "data": [...],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Store by ID

Get a store by ID.

```http
GET /stores/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

---

### Update Store

Update a store.

```http
PUT /stores/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Store 001 Updated",
  "is_active": true
}
```

**Response:** `200 OK`

---

### Delete Store

Delete a store.

```http
DELETE /stores/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

---

## Products

### Create Product

Create a new product.

```http
POST /products
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "sku": "PROD001",
  "barcode": "1234567890123",
  "name": "Product Name",
  "category_id": "uuid-or-null",
  "base_uom_id": "uuid",
  "is_stockable": true,
  "length": 10.5,
  "width": 20.0,
  "height": 5.0,
  "weight": 2.5,
  "is_stackable": true,
  "max_stack_layer": 5
}
```

**Validation:**

- `sku`: Required, max 50 characters, must be unique
- `barcode`: Optional, max 50 characters, must be unique
- `name`: Required, max 200 characters
- `base_uom_id`: Required (UUID)

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Product created successfully",
  "data": {
    "id": "uuid",
    "sku": "PROD001",
    "barcode": "1234567890123",
    "name": "Product Name",
    "category_id": "uuid",
    "base_uom_id": "uuid",
    "is_stockable": true,
    "length": 10.5,
    "width": 20.0,
    "height": 5.0,
    "weight": 2.5,
    "is_stackable": true,
    "max_stack_layer": 5,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

---

### Get All Products

Get all products.

```http
GET /products
Authorization: Bearer {access_token}
```

---

### Get Products with Pagination

Get products with pagination.

```http
GET /products/pagination?page=1&limit=10&search=&order_column=id&order_dir=asc
Authorization: Bearer {access_token}
```

---

### Get Product by ID

Get a product by ID.

```http
GET /products/{id}
Authorization: Bearer {access_token}
```

---

### Update Product

Update a product.

```http
PUT /products/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Product Name",
  "is_stockable": true
}
```

---

### Delete Product

Delete a product.

```http
DELETE /products/{id}
Authorization: Bearer {access_token}
```

---

## Product Categories

### Create Category

```http
POST /categories
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Category Name",
  "code": "CAT001",
  "description": "Category description"
}
```

---

### Get All Categories

```http
GET /categories
Authorization: Bearer {access_token}
```

---

### Get Categories with Pagination

```http
GET /categories/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

### Get Category by ID

```http
GET /categories/{id}
Authorization: Bearer {access_token}
```

---

### Update Category

```http
PUT /categories/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Category Name"
}
```

---

### Delete Category

```http
DELETE /categories/{id}
Authorization: Bearer {access_token}
```

---

## Units of Measure (UOM)

### Create UOM

```http
POST /uoms
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Piece",
  "code": "PC",
  "symbol": "pc",
  "is_active": true
}
```

---

### Get All UOMs

```http
GET /uoms
Authorization: Bearer {access_token}
```

---

### Get UOMs with Pagination

```http
GET /uoms/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

### Get UOM by ID

```http
GET /uoms/{id}
Authorization: Bearer {access_token}
```

---

### Update UOM

```http
PUT /uoms/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated UOM",
  "is_active": true
}
```

---

## Product Prices

### Create Product Price

```http
POST /product-prices
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "product_id": "uuid",
  "effective_date": "2026-05-13",
  "cost_price": 10000.00,
  "selling_price": 15000.00,
  "minimum_quantity": 1,
  "is_active": true
}
```

---

### Get All Product Prices

```http
GET /product-prices
Authorization: Bearer {access_token}
```

---

### Get Product Price by ID

```http
GET /product-prices/{id}
Authorization: Bearer {access_token}
```

---

### Get Prices by Product ID

```http
GET /product-prices/product/{productId}
Authorization: Bearer {access_token}
```

---

### Update Product Price

```http
PUT /product-prices/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "selling_price": 16000.00
}
```

---

### Delete Product Price

```http
DELETE /product-prices/{id}
Authorization: Bearer {access_token}
```

---

## Product UOM Conversions

### Create UOM Conversion

```http
POST /product-uoms
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "product_id": "uuid",
  "from_uom_id": "uuid",
  "to_uom_id": "uuid",
  "conversion_factor": 12.0,
  "is_active": true
}
```

---

### Get All UOM Conversions

```http
GET /product-uoms
Authorization: Bearer {access_token}
```

---

### Get UOM Conversion by ID

```http
GET /product-uoms/{id}
Authorization: Bearer {access_token}
```

---

### Get Conversions by Product ID

```http
GET /product-uoms/product/{productId}
Authorization: Bearer {access_token}
```

---

### Update UOM Conversion

```http
PUT /product-uoms/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "conversion_factor": 13.0
}
```

---

### Delete UOM Conversion

```http
DELETE /product-uoms/{id}
Authorization: Bearer {access_token}
```

---

## Suppliers

### Create Supplier

```http
POST /suppliers
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "code": "SUP001",
  "name": "Supplier Name",
  "contact_person": "John Doe",
  "email": "supplier@example.com",
  "phone": "021123456",
  "address": "Jl. Supplier Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "is_active": true
}
```

---

### Get All Suppliers

```http
GET /suppliers
Authorization: Bearer {access_token}
```

---

### Get Suppliers with Pagination

```http
GET /suppliers/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

### Get Supplier by ID

```http
GET /suppliers/{id}
Authorization: Bearer {access_token}
```

---

### Update Supplier

```http
PUT /suppliers/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Supplier Name"
}
```

---

### Delete Supplier

```http
DELETE /suppliers/{id}
Authorization: Bearer {access_token}
```

---

## Finance

### Chart of Accounts

#### Create Account

```http
POST /accounts
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "account_number": "1000",
  "account_name": "Cash",
  "account_type": "ASSET",
  "is_active": true
}
```

---

#### Get All Accounts

```http
GET /accounts
Authorization: Bearer {access_token}
```

---

#### Get Accounts with Pagination

```http
GET /accounts/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Account by ID

```http
GET /accounts/{id}
Authorization: Bearer {access_token}
```

---

#### Get Accounts by Type

```http
GET /accounts/type/{type}
Authorization: Bearer {access_token}
```

Account Types: `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`

---

#### Update Account

```http
PUT /accounts/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "account_name": "Updated Account Name"
}
```

---

### Customers

#### Create Customer

```http
POST /customers
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Customer Name",
  "code": "CUST001",
  "email": "customer@example.com",
  "phone": "021123456",
  "address": "Jl. Customer Street",
  "is_active": true
}
```

---

#### Get All Customers

```http
GET /customers
Authorization: Bearer {access_token}
```

---

#### Get Customers with Pagination

```http
GET /customers/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Customer by ID

```http
GET /customers/{id}
Authorization: Bearer {access_token}
```

---

#### Update Customer

```http
PUT /customers/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Customer Name"
}
```

---

#### Delete Customer

```http
DELETE /customers/{id}
Authorization: Bearer {access_token}
```

---

### Payment Methods

#### Create Payment Method

```http
POST /payment-methods
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Bank Transfer",
  "code": "BT",
  "description": "Payment via bank transfer",
  "is_active": true
}
```

---

#### Get All Payment Methods

```http
GET /payment-methods
Authorization: Bearer {access_token}
```

---

#### Get Payment Methods with Pagination

```http
GET /payment-methods/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Payment Method by ID

```http
GET /payment-methods/{id}
Authorization: Bearer {access_token}
```

---

#### Update Payment Method

```http
PUT /payment-methods/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Payment Method"
}
```

---

#### Delete Payment Method

```http
DELETE /payment-methods/{id}
Authorization: Bearer {access_token}
```

---

### Price Lists

#### Create Price List

```http
POST /price-lists
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Price List 2026",
  "code": "PL001",
  "effective_date": "2026-05-13",
  "is_active": true
}
```

---

#### Get All Price Lists

```http
GET /price-lists
Authorization: Bearer {access_token}
```

---

#### Get Active Price Lists

```http
GET /price-lists/active
Authorization: Bearer {access_token}
```

---

#### Get Price List by ID

```http
GET /price-lists/{id}
Authorization: Bearer {access_token}
```

---

#### Update Price List

```http
PUT /price-lists/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Price List"
}
```

---

#### Delete Price List

```http
DELETE /price-lists/{id}
Authorization: Bearer {access_token}
```

---

### Taxes

#### Create Tax

```http
POST /taxes
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "VAT 10%",
  "code": "VAT10",
  "rate": 10.00,
  "is_active": true
}
```

---

#### Get All Taxes

```http
GET /taxes
Authorization: Bearer {access_token}
```

---

#### Get Taxes with Pagination

```http
GET /taxes/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Tax by ID

```http
GET /taxes/{id}
Authorization: Bearer {access_token}
```

---

#### Update Tax

```http
PUT /taxes/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "rate": 12.00
}
```

---

#### Delete Tax

```http
DELETE /taxes/{id}
Authorization: Bearer {access_token}
```

---

## Warehouses

### Create Warehouse

```http
POST /warehouses
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "code": "WH001",
  "name": "Main Warehouse",
  "address": "Jl. Warehouse Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "is_active": true
}
```

---

### Get All Warehouses

```http
GET /warehouses
Authorization: Bearer {access_token}
```

---

### Get Warehouses with Pagination

```http
GET /warehouses/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

### Get Warehouse by ID

```http
GET /warehouses/{id}
Authorization: Bearer {access_token}
```

---

### Update Warehouse

```http
PUT /warehouses/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Warehouse Name"
}
```

---

### Delete Warehouse

```http
DELETE /warehouses/{id}
Authorization: Bearer {access_token}
```

---

## Master Data

### Master Data Proposals

#### Get All Proposals

```http
GET /master-data
Authorization: Bearer {access_token}
```

---

#### Create Proposal

```http
POST /master-data
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "entity_type": "PRODUCT",
  "entity_id": "uuid",
  "proposal_type": "CREATE",
  "data": {}
}
```

---

#### Get Proposals with Pagination

```http
GET /master-data/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Pending Proposals

```http
GET /master-data/pending
Authorization: Bearer {access_token}
```

---

#### Get Proposals by Entity Type

```http
GET /master-data/entity/{entityType}
Authorization: Bearer {access_token}
```

---

#### Get Proposals by Group

```http
GET /master-data/group/{groupId}
Authorization: Bearer {access_token}
```

---

#### Get Proposal by ID

```http
GET /master-data/{id}
Authorization: Bearer {access_token}
```

---

#### Review Proposal

```http
POST /master-data/{id}/review
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "status": "APPROVED",
  "notes": "Approved by admin"
}
```

---

#### Execute Proposal

```http
POST /master-data/{id}/execute
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "executed_by": "uuid"
}
```

---

#### Bulk Link Product-Supplier

```http
POST /master-data/bulk/product-supplier
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "product_id": "uuid",
  "supplier_ids": ["uuid1", "uuid2"]
}
```

---

## Purchase Orders

### Create Purchase Order

```http
POST /purchase-orders
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "store_id": "uuid",
  "supplier_id": "uuid",
  "delivery_date": "2026-05-20",
  "items": [
    {
      "product_id": "uuid",
      "quantity": 100,
      "unit_price": 10000.00
    }
  ]
}
```

---

### Create Purchase Order from Planning

```http
POST /purchase-orders/from-planning
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "planning_id": "uuid"
}
```

---

### Bulk Create from Approved Planning

```http
POST /purchase-orders/bulk-from-approved
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "planning_ids": ["uuid1", "uuid2"]
}
```

---

### Get Purchase Orders with Pagination

```http
GET /purchase-orders/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

### Get Purchase Orders by Store

```http
GET /purchase-orders/store/{storeId}
Authorization: Bearer {access_token}
```

---

### Get Pending Purchase Orders by Store

```http
GET /purchase-orders/store/{storeId}/pending
Authorization: Bearer {access_token}
```

---

### Get Purchase Order by ID

```http
GET /purchase-orders/{id}
Authorization: Bearer {access_token}
```

---

### Submit Purchase Order

```http
POST /purchase-orders/{id}/submit
Authorization: Bearer {access_token}
```

---

### Approve Purchase Order

```http
POST /purchase-orders/{id}/approve
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "approved_by": "uuid"
}
```

---

### Cancel Purchase Order

```http
POST /purchase-orders/cancel
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "purchase_order_id": "uuid",
  "reason": "Cancellation reason"
}
```

---

## Goods Receipts

### Create Goods Receipt

```http
POST /goods-receipts
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "purchase_order_id": "uuid",
  "warehouse_id": "uuid",
  "receipt_date": "2026-05-13",
  "items": [
    {
      "po_item_id": "uuid",
      "quantity_received": 100
    }
  ]
}
```

---

### Create Goods Receipt with Invoice

```http
POST /goods-receipts/with-invoice
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "purchase_order_id": "uuid",
  "warehouse_id": "uuid",
  "invoice_number": "INV001",
  "invoice_date": "2026-05-13",
  "items": [...]
}
```

---

### Get Goods Receipts with Pagination

```http
GET /goods-receipts/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

### Get Goods Receipts by Purchase Order

```http
GET /goods-receipts/po/{poId}
Authorization: Bearer {access_token}
```

---

### Get Goods Receipts by Warehouse

```http
GET /goods-receipts/warehouse/{warehouseId}
Authorization: Bearer {access_token}
```

---

### Get Goods Receipt by ID

```http
GET /goods-receipts/{id}
Authorization: Bearer {access_token}
```

---

### Confirm Goods Receipt

```http
POST /goods-receipts/{id}/confirm
Authorization: Bearer {access_token}
```

---

### Cancel Goods Receipt

```http
POST /goods-receipts/{id}/cancel
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "reason": "Cancellation reason"
}
```

---

## Purchase Payments

### Create Purchase Payment

Create a new purchase payment to pay supplier invoices.

```http
POST /purchase-payments
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "supplier_id": "uuid",
  "payment_account_id": "uuid",
  "ap_account_id": "uuid",
  "payment_date": "2026-05-14",
  "payment_mode": "TRANSFER",
  "reference_no": "TRF-001",
  "notes": "Payment for invoice PI/2026/04/0001",
  "items": [
    {
      "purchase_invoice_id": "uuid",
      "paid_amount": 5000000.00
    }
  ]
}
```

**Validation:**

- `supplier_id`: Required (UUID)
- `payment_account_id`: Required (UUID) - Cash/Bank account
- `ap_account_id`: Required (UUID) - Accounts Payable account
- `payment_date`: Required
- `payment_mode`: Required, one of: `CASH`, `TRANSFER`, `GIRO`
- `giro_number`: Required when payment_mode is `GIRO`
- `items`: Required, at least one item
- `items[].paid_amount`: Cannot exceed remaining amount on invoice

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Purchase payment created successfully",
  "data": {
    "id": "uuid",
    "payment_number": "PAY/2605/00001",
    "supplier_id": "uuid",
    "payment_account_id": "uuid",
    "ap_account_id": "uuid",
    "payment_date": "2026-05-14",
    "payment_mode": "TRANSFER",
    "total_amount": 5000000.00,
    "status": "DRAFT",
    "created_at": "2026-05-14T10:00:00Z"
  }
}
```

---

### Get Purchase Payments with Pagination

Get all purchase payments with pagination.

```http
GET /purchase-payments/pagination?page=1&limit=10&search=&order_column=created_at&order_dir=desc
Authorization: Bearer {access_token}
```

**Query Parameters:**

- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 10) - Items per page
- `search` (string) - Search term (searches payment_number, reference_no)
- `order_column` (string, default: created_at) - Column to sort by
- `order_dir` (string, default: desc) - Sort direction (asc/desc)

**Response:** `200 OK`

---

### Get Purchase Payment by ID

Get a purchase payment by ID.

```http
GET /purchase-payments/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

---

### Post Purchase Payment

Post a purchase payment. This creates journal entries:
- Debit: Accounts Payable
- Credit: Cash/Bank

```http
POST /purchase-payments/{id}/post
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "notes": "Payment posted"
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "posted"
}
```

---

### Void Purchase Payment

Void a posted purchase payment. This creates reversal journal entries.

```http
POST /purchase-payments/{id}/void
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "reason": "Payment cancelled due to invoice discrepancy"
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "voided"
}
```

---

## Purchase Order Planning

### Calculate Planning

```http
POST /planning/calculate
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "store_id": "uuid",
  "calculation_date": "2026-05-13"
}
```

---

### Get Pending Planning

```http
GET /planning/pending
Authorization: Bearer {access_token}
```

---

### Get All Planning

```http
GET /planning
Authorization: Bearer {access_token}
```

---

### Get Planning by ID

```http
GET /planning/{id}
Authorization: Bearer {access_token}
```

---

### Approve Planning

```http
POST /planning/approve
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "planning_id": "uuid"
}
```

---

### Ignore Planning

```http
POST /planning/ignore
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "planning_id": "uuid",
  "reason": "Ignore reason"
}
```

---

## Roles & Permissions

### Roles

#### Create Role

```http
POST /roles
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Admin",
  "description": "Administrator role",
  "permissions": ["uuid1", "uuid2"]
}
```

---

#### List Roles

```http
GET /roles
Authorization: Bearer {access_token}
```

---

#### List Roles with Pagination

```http
GET /roles/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Role by Name

```http
GET /roles/name/{name}
Authorization: Bearer {access_token}
```

---

#### Get Role by ID

```http
GET /roles/{id}
Authorization: Bearer {access_token}
```

---

#### Update Role

```http
PUT /roles/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "description": "Updated description"
}
```

---

#### Delete Role

```http
DELETE /roles/{id}
Authorization: Bearer {access_token}
```

---

### Permissions

#### Create Permission

```http
POST /permissions
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "CREATE_PRODUCT",
  "description": "Permission to create products",
  "resource": "PRODUCT",
  "action": "CREATE"
}
```

---

#### List Permissions

```http
GET /permissions
Authorization: Bearer {access_token}
```

---

#### List Permissions with Pagination

```http
GET /permissions/pagination?page=1&limit=10
Authorization: Bearer {access_token}
```

---

#### Get Permission by ID

```http
GET /permissions/{id}
Authorization: Bearer {access_token}
```

---

#### Update Permission

```http
PUT /permissions/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "description": "Updated description"
}
```

---

#### Delete Permission

```http
DELETE /permissions/{id}
Authorization: Bearer {access_token}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "code": 400,
  "message": "Error message",
  "errors": [
    {
      "field": "field_name",
      "message": "Error details"
    }
  ]
}
```

### Common Status Codes

| Code | Description                                      |
| ---- | ------------------------------------------------ |
| 200  | OK - Request successful                          |
| 201  | Created - Resource created successfully          |
| 400  | Bad Request - Invalid request body               |
| 401  | Unauthorized - Missing or invalid authentication |
| 403  | Forbidden - Insufficient permissions             |
| 404  | Not Found - Resource not found                   |
| 409  | Conflict - Resource already exists               |
| 422  | Unprocessable Entity - Validation error          |
| 500  | Internal Server Error - Server error             |

---

## Pagination

Paginated endpoints return responses with metadata:

```json
{
  "code": 200,
  "message": "Success message",
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "total_pages": 10
  }
}
```

**Parameters:**

- `page` - Current page (default: 1)
- `limit` - Items per page (default: 10)
- `search` - Search term (optional)
- `order_column` - Sort column (default: id)
- `order_dir` - Sort direction: `asc` or `desc` (default: asc)

---

## Authentication

All protected endpoints require the `Authorization` header with a Bearer token:

```
Authorization: Bearer {access_token}
```

The token is obtained from the `/auth/login` endpoint and is valid for 1 hour. Use the refresh token to obtain a new access token.

---

## Rate Limiting

API rate limiting is applied per user. Contact administrator for specific limits.

---

## Versioning

Current API version: `v1`

---

## Support

For API support and issues, please contact the development team.

---

**Last Updated:** May 13, 2026
