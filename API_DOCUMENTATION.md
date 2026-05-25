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
2. [Users](#users)
3. [Stores](#stores)
4. [Products](#products)
5. [Product Categories](#product-categories)
6. [Units of Measure (UOM)](#units-of-measure-uom)
7. [Suppliers](#suppliers)
8. [Chart of Accounts](#chart-of-accounts)
9. [Warehouses](#warehouses)
10. [Master Data](#master-data)
11. [Purchase Orders](#purchase-orders)
12. [Goods Receipts](#goods-receipts)
13. [Purchase Payments](#purchase-payments)
14. [Purchase Invoices](#purchase-invoices)
15. [Purchase Returns](#purchase-returns)
16. [Expense Vouchers](#expense-vouchers)
17. [Planning](#planning)
18. [Roles & Permissions](#roles--permissions)

---

## Grouping

### Data Master

- [Master Data Proposals](#master-data)
- [Users](#users)
- [Stores](#stores)
- [Products](#products)
- [Product Categories](#product-categories)
- [Suppliers](#suppliers)
- [Warehouses](#warehouses)

### Finance

- [Chart of Accounts](#chart-of-accounts)
- [Customers](#finance)
- [Payment Methods](#finance)
- [Price Lists](#finance)
- [Taxes](#finance)
- [Purchase Invoices](#purchase-invoices)
- [Purchase Returns](#purchase-returns)
- [Expense Vouchers](#expense-vouchers)

### Transactions

- [Purchase Orders](#purchase-orders)
- [Goods Receipts](#goods-receipts)
- [Purchase Payments](#purchase-payments)
- [Purchase Invoices](#purchase-invoices)
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

### Google OAuth Redirect

Redirect user to Google login page.

```http
GET /auth/google
```

**Response:** `302 Redirect` to Google OAuth URL

---

### Google OAuth Callback

Callback handler after Google authentication.

```http
GET /auth/google/callback?code={authorization_code}
```

**Response:** `200 OK`
```json
{
  "code": 200,
  "message": "Login successful via Google",
  "data": {
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "user@gmail.com",
      "avatar_url": "https://..."
    }
  }
}
```

**Error:** `400 Bad Request` - Code not found

---

### Google Token Login

Login menggunakan Google Access Token atau ID Token dari frontend.

```http
POST /auth/google/token
Content-Type: application/json

{
  "token": "google-access-token-or-id-token",
  "token_type": "access"
}
```

**Validation:**
- `token`: Required
- `token_type`: Required, one of `access`, `id`

**Response:** `200 OK`
```json
{
  "code": 200,
  "message": "Login successful via Google Token",
  "data": {
    "access_token": "eyJhbGc...",
    "user": { ... }
  }
}
```

**Error:** `401 Unauthorized` - Invalid credentials

---

### Register with Google

Mendaftar akun baru menggunakan Google token.

```http
POST /auth/google/register
Content-Type: application/json

{
  "token": "google-access-token-or-id-token",
  "token_type": "access"
}
```

**Response:** `201 Created`
```json
{
  "code": 201,
  "message": "User registered successfully via Google",
  "data": {
    "access_token": "eyJhbGc...",
    "user": { ... }
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

## Users

Manajemen pengguna sistem. Endpoint ini hanya dapat diakses oleh admin.

### Get All Users

```http
GET /users
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "store_id": "uuid",
      "store_name": "Store 001",
      "name": "John Doe",
      "username": "johndoe",
      "email": "john@example.com",
      "phone": "021123456",
      "role": "ADMIN",
      "avatar_url": null,
      "is_active": true,
      "last_login_at": "2026-05-13T10:00:00Z",
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z",
      "permissions": ["products:view", "products:create"]
    }
  ]
}
```

---

### Get Users with Pagination

```http
GET /users/pagination?page=1&limit=10&search=&order_column=id&order_dir=asc
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
  "message": "Users retrieved successfully",
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

### Get User by ID

```http
GET /users/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "ADMIN",
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `404 Not Found` - User not found

---

### Update User

Update user by ID. Semua field bersifat opsional.

```http
PUT /users/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "John Updated",
  "username": "johnupdated",
  "email": "john.updated@example.com",
  "phone": "08123456789",
  "role": "SUPERVISOR",
  "is_active": true
}
```

**Validation:**
- `store_id`: Required (UUID) - ID toko yang akan diassign
- `name`: Optional, max 255
- `username`: Optional, max 100
- `email`: Optional, valid email, max 255
- `phone`: Optional, max 20
- `password`: Optional, min 8
- `role`: Optional, max 50
- `is_active`: Optional (boolean)

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "id": "uuid",
    "name": "John Updated",
    "email": "john.updated@example.com",
    "role": "SUPERVISOR",
    "updated_at": "2026-05-13T11:00:00Z"
  }
}
```

**Error:** `404 Not Found` - User not found

---

### Delete User

Soft-delete user by ID.

```http
DELETE /users/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "User deleted successfully",
  "data": null
}
```

**Error:** `404 Not Found` - User not found

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

```json
{
  "code": 200,
  "message": "Store retrieved successfully",
  "data": {
    "id": "uuid",
    "code": "STR001",
    "name": "Store 001",
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Products retrieved successfully",
  "data": [
    {
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
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Product by ID

Get a product by ID.

```http
GET /products/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Product retrieved successfully",
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

Mengelola hierarki kategori produk menggunakan pola **Adjacency List** yang mendukung kategori bertingkat tak terbatas. Setiap kategori memiliki slug untuk URL-friendly identifier dan default markup percentage untuk auto-pricing.

### Create Category

Membuat kategori produk baru. ParentID dapat dikosongkan untuk kategori utama (root).

```http
POST /categories
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "parent_id": "uuid-or-null",
  "name": "Sembako",
  "slug": "sembako",
  "default_markup_pct": 10.00
}
```

**Validation:**
- `name`: Required, min 1, max 100 characters
- `slug`: Required, min 1, max 120 characters, must be unique
- `parent_id`: Optional (UUID), null untuk root category
- `default_markup_pct`: Optional, decimal(5,2), default 0

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Category created successfully",
  "data": {
    "id": "uuid",
    "parent_id": null,
    "parent": null,
    "name": "Sembako",
    "slug": "sembako",
    "default_markup_pct": 10.00,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `409 Conflict` - Slug already exists

---

### Get All Categories

```http
GET /categories
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "parent_id": null,
      "parent": null,
      "name": "Sembako",
      "slug": "sembako",
      "default_markup_pct": 10.00,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ]
}
```

---

### Get Categories with Pagination

```http
GET /categories/pagination?page=1&limit=10&search=&order_column=id&order_dir=asc
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
  "message": "Categories retrieved successfully",
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

### Get Category by ID

```http
GET /categories/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Category retrieved successfully",
  "data": {
    "id": "uuid",
    "parent_id": "uuid",
    "parent": {
      "id": "uuid",
      "name": "Induk Kategori",
      "slug": "induk-kategori"
    },
    "name": "Sub Kategori",
    "slug": "sub-kategori",
    "default_markup_pct": 5.00,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `404 Not Found` - Category not found

---

### Update Category

```http
PUT /categories/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "parent_id": "uuid-or-null",
  "name": "Sembako Updated",
  "slug": "sembako-updated",
  "default_markup_pct": 12.50
}
```

Semua field bersifat opsional pada update.

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Category updated successfully",
  "data": {
    "id": "uuid",
    "parent_id": null,
    "name": "Sembako Updated",
    "slug": "sembako-updated",
    "default_markup_pct": 12.50,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T11:00:00Z"
  }
}
```

**Error:** `409 Conflict` - Slug already exists

---

### Delete Category

```http
DELETE /categories/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Category deleted successfully",
  "data": null
}
```

**Error:** `404 Not Found` - Category not found

---

## Units of Measure (UOM)

Unit of Measure (UOM) merepresentasikan satuan ukuran dasar yang digunakan untuk produk. Setiap UOM memiliki kode unik (max 10 karakter) dan nama (max 50 karakter). UOM tidak memiliki endpoint DELETE; untuk menonaktifkan cukup tidak digunakan lagi.

### Create UOM

Membuat satuan ukuran baru.

```http
POST /uoms
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "code": "PC",
  "name": "Piece"
}
```

**Validation:**
- `code`: Required, min 1, max 10 characters, must be unique
- `name`: Required, min 1, max 50 characters

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "UOM created successfully",
  "data": {
    "id": "uuid",
    "code": "PC",
    "name": "Piece",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `400 Bad Request` - Validation error (missing required fields)
```json
{
  "code": 400,
  "message": "Validation failed",
  "errors": [
    { "field": "code", "message": "Code is required" },
    { "field": "name", "message": "Name is required" }
  ]
}
```

**Error:** `409 Conflict` - Code already exists
```json
{
  "code": 409,
  "message": "Conflict",
  "errors": [
    { "field": "code", "message": "code 'PC' already exists" }
  ]
}
```

---

### Get All UOMs

Mengambil seluruh daftar satuan ukuran (tanpa pagination).

```http
GET /uoms
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "UOMs retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "code": "PC",
      "name": "Piece",
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ]
}
```

---

### Get UOMs with Pagination

```http
GET /uoms/pagination?page=1&limit=10&search=&order_column=id&order_dir=asc
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 10) - Items per page
- `search` (string) - Search term (searches by code or name)
- `order_column` (string, default: id) - Column to sort by
- `order_dir` (string, default: asc) - Sort direction (asc/desc)

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "UOMs retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "code": "PC",
      "name": "Piece",
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get UOM by ID

```http
GET /uoms/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "UOM retrieved successfully",
  "data": {
    "id": "uuid",
    "code": "PC",
    "name": "Piece",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `404 Not Found` - UOM not found

---

### Update UOM

Semua field bersifat opsional pada update.

```http
PUT /uoms/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "code": "PCS",
  "name": "Pieces"
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "UOM updated successfully",
  "data": {
    "id": "uuid",
    "code": "PCS",
    "name": "Pieces",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T11:00:00Z"
  }
}
```

**Error:** `404 Not Found` - UOM not found
**Error:** `409 Conflict` - Code already exists

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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Product price retrieved successfully",
  "data": {
    "id": "uuid",
    "product_id": "uuid",
    "effective_date": "2026-05-13",
    "cost_price": 10000.00,
    "selling_price": 15000.00,
    "minimum_quantity": 1,
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

---

### Get Prices by Product ID

```http
GET /product-prices/product/{productId}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Product prices retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "effective_date": "2026-05-13",
      "cost_price": 10000.00,
      "selling_price": 15000.00,
      "minimum_quantity": 1,
      "is_active": true,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ]
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "UOM conversion retrieved successfully",
  "data": {
    "id": "uuid",
    "product_id": "uuid",
    "from_uom_id": "uuid",
    "to_uom_id": "uuid",
    "conversion_factor": 12.0,
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

---

### Get Conversions by Product ID

```http
GET /product-uoms/product/{productId}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "UOM conversions retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "from_uom_id": "uuid",
      "to_uom_id": "uuid",
      "conversion_factor": 12.0,
      "is_active": true,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ]
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Suppliers retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "code": "SUP001",
      "name": "Supplier Name",
      "contact_person": "John Doe",
      "phone_number": "021123456",
      "email": "supplier@example.com",
      "address": "Jl. Supplier Street",
      "is_active": true,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Supplier by ID

```http
GET /suppliers/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Supplier retrieved successfully",
  "data": {
    "id": "uuid",
    "code": "SUP001",
    "name": "Supplier Name",
    "contact_person": "John Doe",
    "phone_number": "021123456",
    "email": "supplier@example.com",
    "address": "Jl. Supplier Street",
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

Chart of Accounts (COA) mendukung hierarki bertingkat menggunakan pola **Adjacency List** dengan field `parent_id`. Setiap akun dapat memiliki akun induk untuk pengelompokan laporan keuangan.

#### Create Account

Membuat akun baru. `parent_id` dapat dikosongkan untuk akun root (induk utama).

```http
POST /accounts
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "account_code": "1110.01",
  "name": "Kas di Tangan",
  "account_type": "ASSET",
  "normal_balance": "DEBIT",
  "parent_id": "uuid-akun-induk-atau-null"
}
```

**Validation:**
- `account_code`: Required, max 20 characters, unique
- `name`: Required, max 100 characters
- `account_type`: Required, one of `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`
- `normal_balance`: Required, one of `DEBIT`, `CREDIT`
- `parent_id`: Optional (UUID), null untuk root account

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Chart of account created successfully",
  "data": {
    "id": "uuid",
    "account_code": "1110.01",
    "name": "Kas di Tangan",
    "account_type": "ASSET",
    "normal_balance": "DEBIT",
    "is_active": true,
    "parent_id": "uuid-parent",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `409 Conflict` - Account code already exists

---

#### Get All Accounts

```http
GET /accounts
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of accounts retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "account_code": "1110.01",
      "name": "Kas di Tangan",
      "account_type": "ASSET",
      "normal_balance": "DEBIT",
      "is_active": true,
      "parent_id": "uuid-parent",
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ]
}
```

---

#### Get Accounts with Pagination

```http
GET /accounts/pagination?page=1&limit=10&search=&order_column=id&order_dir=asc
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (integer, default: 1)
- `limit` (integer, default: 10)
- `search` (string) - Search by account_code or name
- `order_column` (string, default: id)
- `order_dir` (string, default: asc)

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of accounts retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "account_code": "1110.01",
      "name": "Kas di Tangan",
      "account_type": "ASSET",
      "normal_balance": "DEBIT",
      "is_active": true,
      "parent_id": null,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

#### Get Account by ID

```http
GET /accounts/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of account retrieved successfully",
  "data": {
    "id": "uuid",
    "account_code": "1110.01",
    "name": "Kas di Tangan",
    "account_type": "ASSET",
    "normal_balance": "DEBIT",
    "is_active": true,
    "parent_id": "uuid-parent",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

**Error:** `404 Not Found`

---

#### Get Accounts by Type

```http
GET /accounts/type/{type}
Authorization: Bearer {access_token}
```

Account Types: `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of accounts retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "account_code": "1110.01",
      "name": "Kas di Tangan",
      "account_type": "ASSET",
      "normal_balance": "DEBIT",
      "is_active": true,
      "parent_id": null,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ]
}
```

---

#### Get Accounts Tree

Mengembalikan seluruh hierarki akun dalam format tree (parent → children rekursif).

```http
GET /accounts/tree
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of accounts tree retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "account_code": "1000",
      "name": "ASET",
      "account_type": "ASSET",
      "normal_balance": "DEBIT",
      "is_active": true,
      "children": [
        {
          "id": "uuid",
          "account_code": "1100",
          "name": "ASET LANCAR",
          "account_type": "ASSET",
          "normal_balance": "DEBIT",
          "is_active": true,
          "children": [
            {
              "id": "uuid",
              "account_code": "1110.01",
              "name": "Kas di Tangan",
              "account_type": "ASSET",
              "normal_balance": "DEBIT",
              "is_active": true,
              "children": []
            }
          ]
        }
      ]
    }
  ]
}
```

---

#### Update Account

```http
PUT /accounts/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "account_code": "1110.02",
  "name": "Kas di Bank",
  "account_type": "ASSET",
  "normal_balance": "DEBIT",
  "is_active": true,
  "parent_id": "uuid-parent-atau-kosongkan-untuk-hapus-parent"
}
```

Semua field bersifat opsional. Kirim `parent_id: ""` untuk menghapus relasi induk.

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of account updated successfully",
  "data": {
    "id": "uuid",
    "account_code": "1110.02",
    "name": "Kas di Bank",
    "account_type": "ASSET",
    "normal_balance": "DEBIT",
    "is_active": true,
    "parent_id": null,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T11:00:00Z"
  }
}
```

**Error:** `404 Not Found` | `409 Conflict` (duplicate code)

---

#### Import Accounts from Excel

Import akun secara massal dari file Excel (.xlsx/.xls).

```http
POST /accounts/import
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

file: account_data.xlsx
```

**Struktur Excel (5 kolom):**
| account_code | name | account_type | normal_balance | parent_code |
|-------------|------|-------------|---------------|-------------|
| 1110.01 | Kas di Tangan | ASSET | DEBIT | 1100 |
| 1110.02 | Kas di Bank | ASSET | DEBIT | 1100 |

- `parent_code` bersifat opsional, diisi dengan `account_code` dari akun induk
- Validasi duplikat, tipe akun, dan normal balance dilakukan per baris
- Baris error akan di-skip dan dilaporkan

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Import completed",
  "data": {
    "total_rows": 10,
    "success_rows": 8,
    "error_rows": 2,
    "errors": [
      {
        "row_number": 3,
        "message": "account code '1110.01' already exists"
      },
      {
        "row_number": 7,
        "message": "invalid account_type 'INCOME', must be one of: ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE"
      }
    ]
  }
}
```

---

#### Download Import Template

Download template Excel untuk import akun. File berisi header dengan dropdown validation untuk `account_type` dan `normal_balance`.

```http
GET /accounts/import/template
Authorization: Bearer {access_token}
```

**Response:** `200 OK` (file .xlsx)

```text
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="account_import_template.xlsx"
```

---

#### Delete Account

Hapus akun. Tidak dapat menghapus akun yang memiliki children (sub-akun).

```http
DELETE /accounts/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Chart of account deleted successfully",
  "data": null
}
```

**Error:** `404 Not Found` - Account not found
**Error:** `409 Conflict` - Account has children, cannot delete

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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Customers retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "code": "CUST001",
      "name": "Customer Name",
      "phone_number": "021123456",
      "email": "customer@example.com",
      "address": "Jl. Customer Street",
      "is_active": true,
      "point_balance": 0,
      "credit_limit": 0,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

#### Get Customer by ID

```http
GET /customers/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Customer retrieved successfully",
  "data": {
    "id": "uuid",
    "code": "CUST001",
    "name": "Customer Name",
    "phone_number": "021123456",
    "email": "customer@example.com",
    "address": "Jl. Customer Street",
    "is_active": true,
    "point_balance": 0,
    "credit_limit": 0,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Payment methods retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "code": "BT",
      "name": "Bank Transfer",
      "mdr_percentage": 0,
      "is_active": true,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

#### Get Payment Method by ID

```http
GET /payment-methods/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Payment method retrieved successfully",
  "data": {
    "id": "uuid",
    "code": "BT",
    "name": "Bank Transfer",
    "mdr_percentage": 0,
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Price list retrieved successfully",
  "data": {
    "id": "uuid",
    "code": "PL001",
    "name": "Price List 2026",
    "effective_date": "2026-05-13",
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Taxes retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "VAT 10%",
      "rate_percentage": 10.00,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

#### Get Tax by ID

```http
GET /taxes/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Tax retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "VAT 10%",
    "rate_percentage": 10.00,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Warehouses retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "store_id": "uuid",
      "code": "WH001",
      "name": "Main Warehouse",
      "is_active": true,
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Warehouse by ID

```http
GET /warehouses/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Warehouse retrieved successfully",
  "data": {
    "id": "uuid",
    "store_id": "uuid",
    "code": "WH001",
    "name": "Main Warehouse",
    "is_active": true,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Proposals retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "reference_number": "PROP/2026/05/00001",
      "entity_type": "PRODUCT",
      "action_type": "CREATE",
      "total_items": 1,
      "status": "PENDING",
      "proposed_by_id": "uuid",
      "reason": "Barang baru Q2",
      "created_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Proposal retrieved successfully",
  "data": {
    "id": "uuid",
    "reference_number": "PROP/2026/05/00001",
    "entity_type": "PRODUCT",
    "action_type": "CREATE",
    "total_items": 1,
    "status": "PENDING",
    "proposed_by_id": "uuid",
    "reason": "Barang baru Q2",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "entity_id": null,
        "payload_json": "{\"sku\":\"PROD002\",\"name\":\"Product Baru\"}"
      }
    ]
  }
}
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

#### Update Proposal

Update proposal yang masih berstatus PENDING.

```http
PUT /master-data/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "reason": "Revisi data barang",
  "items": [
    {
      "entity_id": null,
      "payload_json": "{\"sku\":\"PROD003\",\"name\":\"Product Updated\"}"
    }
  ]
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Proposal updated successfully",
  "data": { ... }
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase orders retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "po_number": "PO/2026/05/00001",
      "supplier_id": "uuid",
      "supplier_name": "Supplier Name",
      "store_id": "uuid",
      "store_name": "Store 001",
      "warehouse_id": "uuid",
      "warehouse_name": "Main Warehouse",
      "order_date": "2026-05-13T10:00:00Z",
      "total_amount": 1000000,
      "status": "DRAFT",
      "created_by_id": "uuid",
      "created_by_name": "John Doe",
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase order retrieved successfully",
  "data": {
    "id": "uuid",
    "po_number": "PO/2026/05/00001",
    "reference_no": "REF-001",
    "supplier": {
      "id": "uuid",
      "code": "SUP001",
      "name": "Supplier Name"
    },
    "store": {
      "id": "uuid",
      "code": "STR001",
      "name": "Store 001"
    },
    "warehouse": {
      "id": "uuid",
      "code": "WH001",
      "name": "Main Warehouse"
    },
    "order_date": "2026-05-13T10:00:00Z",
    "expected_delivery": "2026-05-20T10:00:00Z",
    "payment_term_days": 30,
    "payment_mode": "TRANSFER",
    "total_amount": 1000000,
    "status": "DRAFT",
    "created_by_id": "uuid",
    "created_by_name": "John Doe",
    "notes": "PO untuk restock",
    "supplier_notes": null,
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "product_id": "uuid",
        "product_sku": "PROD001",
        "product_name": "Product Name",
        "uom_id": "uuid",
        "uom_name": "PC",
        "qty_ordered": 100,
        "qty_received": 0,
        "unit_price": 10000,
        "subtotal": 1000000
      }
    ]
  }
}
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

### Update Purchase Order

Update PO yang masih berstatus DRAFT.

```http
PUT /purchase-orders/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "reference_no": "REF-001",
  "expected_delivery": "2026-05-25T10:00:00Z",
  "payment_term_days": 30,
  "payment_mode": "TRANSFER",
  "notes": "Updated notes",
  "supplier_notes": "Notes for supplier",
  "items": [
    {
      "product_id": "uuid",
      "uom_id": "uuid",
      "qty_ordered": 100,
      "unit_price": 10000,
      "notes": "Item notes"
    }
  ]
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase order updated successfully",
  "data": { ... }
}
```

---

## Goods Receipts

### Create Goods Receipt

Digunakan untuk mencatat penerimaan barang fisik dari supplier berdasarkan Purchase Order (PO).

```http
POST /goods-receipts
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "purchase_order_id": "uuid",
  "warehouse_id": "uuid",
  "receipt_date": "2026-05-13",
  "delivery_note_no": "SJ/123",
  "notes": "Barang datang lengkap",
  "override_pin": "123456",
  "items": [
    {
      "purchase_order_item_id": "uuid",
      "product_id": "uuid",
      "uom_id": "uuid",
      "qty_received": 100,
      "qty_rejected": 0,
      "unit_price": 50000,
      "reject_reason": null,
      "notes": null
    }
  ]
}
```

**Over-Receiving Validation:**
Sistem akan menghitung sisa kuantitas: `Sisa = Qty Dipesan - Qty Diterima (Confirm) - Qty Draft Lainnya`.
Jika `qty_received > Sisa`, sistem akan mengembalikan error `ERR_OVER_RECEIVE_NEEDS_PIN`. Untuk melanjutkan, user harus menyertakan `override_pin` milik Supervisor/Kepala Gudang.

---

### Update Goods Receipt

Mengubah data Goods Receipt yang masih berstatus `DRAFT`.

```http
PUT /goods-receipts/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "receipt_date": "2026-05-13",
  "delivery_note_no": "SJ/123-REV",
  "notes": "Revisi kuantitas",
  "override_pin": "123456",
  "items": [...]
}
```

---

### Create Goods Receipt with Invoice (Direct Mode)

Shortcut untuk membuat Goods Receipt sekaligus Purchase Invoice (otomatis CONFIRMED & POSTED).

```http
POST /goods-receipts/with-invoice
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "purchase_order_id": "uuid",
  "warehouse_id": "uuid",
  "receipt_date": "2026-05-13",
  "delivery_note_no": "SJ/123",
  "supplier_invoice_number": "INV/SUPP/001",
  "invoice_date": "2026-05-13",
  "ap_account_id": "uuid",
  "payment_term_days": 30,
  "discount_amount": 0,
  "override_pin": null,
  "items": [...]
}
```

---

### Get Goods Receipts with Pagination

```http
GET /goods-receipts/pagination?page=1&limit=10&search=GR/2026&warehouse_id={uuid}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Goods receipts retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "gr_number": "GR/2026/05/00001",
      "purchase_order_id": "uuid",
      "po_number": "PO/2026/05/00001",
      "warehouse_id": "uuid",
      "warehouse_name": "Main Warehouse",
      "receipt_date": "2026-05-13",
      "delivery_note_no": "SJ/123",
      "status": "DRAFT",
      "received_by_id": "uuid",
      "supplier_name": "Supplier Name",
      "store_name": "Store 001",
      "created_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
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
GET /goods-receipts/warehouse/{warehouseId}?status=DRAFT
Authorization: Bearer {access_token}
```

---

### Get Goods Receipt by ID

```http
GET /goods-receipts/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Goods receipt retrieved successfully",
  "data": {
    "id": "uuid",
    "gr_number": "GR/2026/05/00001",
    "purchase_order_id": "uuid",
    "po_number": "PO/2026/05/00001",
    "warehouse_id": "uuid",
    "warehouse_name": "Main Warehouse",
    "receipt_date": "2026-05-13",
    "delivery_note_no": "SJ/123",
    "status": "DRAFT",
    "received_by_id": "uuid",
    "supplier_name": "Supplier Name",
    "supplier_code": "SUP001",
    "store_name": "Store 001",
    "subtotal": 5000000,
    "discount_amount": 0,
    "tax_amount": 0,
    "grand_total": 5000000,
    "notes": "Barang datang lengkap",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "product_id": "uuid",
        "product_name": "Product Name",
        "product_sku": "PROD001",
        "uom_id": "uuid",
        "uom_code": "PC",
        "qty_ordered": 100,
        "qty_received": 100,
        "qty_rejected": 0,
        "unit_price": 50000,
        "net_unit_price": 50000
      }
    ]
  }
}
```

---

### Confirm Goods Receipt

Mengubah status dari `DRAFT` menjadi `CONFIRMED`. Stok akan bertambah dan HPP dihitung.

```http
POST /goods-receipts/{id}/confirm
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "notes": "Confirm by Warehouse Manager"
}
```

---

### Cancel Goods Receipt

Membatalkan GR (hanya jika masih `DRAFT`).

```http
POST /goods-receipts/{id}/cancel
Authorization: Bearer {access_token}
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

```json
{
  "code": 200,
  "message": "Purchase payments retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "payment_number": "PAY/2605/00001",
      "supplier_id": "uuid",
      "supplier_name": "Supplier Name",
      "payment_date": "2026-05-14",
      "payment_mode": "TRANSFER",
      "total_amount": 5000000,
      "status": "DRAFT",
      "created_at": "2026-05-14T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Purchase Payment by ID

Get a purchase payment by ID.

```http
GET /purchase-payments/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase payment retrieved successfully",
  "data": {
    "id": "uuid",
    "payment_number": "PAY/2605/00001",
    "reference_no": "TRF-001",
    "supplier_id": "uuid",
    "supplier_name": "Supplier Name",
    "payment_account_id": "uuid",
    "ap_account_id": "uuid",
    "payment_date": "2026-05-14",
    "payment_mode": "TRANSFER",
    "total_amount": 5000000,
    "status": "POSTED",
    "created_by_id": "uuid",
    "posted_by_id": "uuid",
    "posted_at": "2026-05-14T11:00:00Z",
    "notes": "Payment for invoice PI/2026/04/0001",
    "created_at": "2026-05-14T10:00:00Z",
    "updated_at": "2026-05-14T11:00:00Z",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "purchase_invoice_id": "uuid",
        "invoice_number": "PI/2026/04/0001",
        "document_amount": 5000000,
        "paid_amount": 5000000
      }
    ]
  }
}
```

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

## Purchase Invoices

Faktur pembelian dari supplier. Mendukung alur workflow: `DRAFT → SUBMITTED → VERIFIED → POSTED → PARTIALLY_PAID → PAID`.

### Create Purchase Invoice

```http
POST /purchase-invoices
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "store_id": "uuid",
  "warehouse_id": "uuid",
  "purchase_order_id": "uuid",
  "supplier_id": "uuid",
  "supplier_invoice_number": "INV/SUPP/001",
  "reference_no": "REF-001",
  "ap_account_id": "uuid",
  "inventory_account_id": "uuid",
  "invoice_date": "2026-05-13",
  "received_date": "2026-05-13",
  "payment_term_days": 30,
  "payment_mode": "TRANSFER",
  "discount_amount": 0,
  "freight_amount": 0,
  "other_cost_amount": 0,
  "is_tax_inclusive": false,
  "notes": "Invoice notes",
  "items": [
    {
      "product_id": "uuid",
      "uom_id": "uuid",
      "qty_invoiced": 100,
      "unit_price": 50000,
      "discount_1_pct": 0,
      "discount_2_pct": 0,
      "discount_3_pct": 0,
      "discount_amount": 0,
      "tax_pct": 10,
      "notes": null
    }
  ]
}
```

**Validation:**
- `store_id`, `warehouse_id`, `purchase_order_id`, `supplier_id`: Required (UUID)
- `supplier_invoice_number`: Required, max 50
- `ap_account_id`, `inventory_account_id`: Required (UUID)
- `invoice_date`, `received_date`: Required
- `payment_mode`: Optional, one of `CASH`, `TRANSFER`, `GIRO`
- `items`: Required, min 1
- `items[].product_id`, `uom_id`: Required
- `items[].qty_invoiced`: Required
- `items[].unit_price`: Required, min 0

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Purchase invoice created successfully",
  "data": {
    "id": "uuid",
    "invoice_number": "PI/2026/05/00001",
    "supplier_invoice_number": "INV/SUPP/001",
    "status": "DRAFT",
    "grand_total": 5500000,
    "created_at": "2026-05-13T10:00:00Z"
  }
}
```

---

### List Purchase Invoices

```http
GET /purchase-invoices?page=1&limit=10&search=&order_by=created_at&order_dir=desc
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (integer, default: 1)
- `limit` (integer, default: 10)
- `search` (string) - Search by invoice number
- `order_by` (string, default: created_at)
- `order_dir` (string, default: desc)

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase invoices retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "invoice_number": "PI/2026/05/00001",
      "supplier_invoice_number": "INV/SUPP/001",
      "po_number": "PO/2026/05/00001",
      "supplier_name": "Supplier Name",
      "supplier_code": "SUP001",
      "store_name": "Store 001",
      "warehouse_name": "Main Warehouse",
      "invoice_date": "2026-05-13",
      "due_date": "2026-06-12",
      "grand_total": 5500000,
      "paid_amount": 0,
      "remaining_amount": 5500000,
      "status": "DRAFT",
      "created_by_name": "John Doe",
      "created_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Purchase Invoice by ID

```http
GET /purchase-invoices/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase invoice retrieved successfully",
  "data": {
    "id": "uuid",
    "invoice_number": "PI/2026/05/00001",
    "supplier_invoice_number": "INV/SUPP/001",
    "po_number": "PO/2026/05/00001",
    "supplier": { "id": "uuid", "code": "SUP001", "name": "Supplier Name" },
    "store": { "id": "uuid", "code": "STR001", "name": "Store 001" },
    "warehouse": { "id": "uuid", "code": "WH001", "name": "Main Warehouse" },
    "invoice_date": "2026-05-13",
    "received_date": "2026-05-13",
    "due_date": "2026-06-12",
    "payment_term_days": 30,
    "payment_mode": "TRANSFER",
    "subtotal": 5000000,
    "discount_amount": 0,
    "tax_amount": 500000,
    "grand_total": 5500000,
    "paid_amount": 0,
    "remaining_amount": 5500000,
    "status": "DRAFT",
    "created_by_name": "John Doe",
    "notes": "Invoice notes",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "product_id": "uuid",
        "product_name": "Product Name",
        "product_sku": "PROD001",
        "uom_code": "PC",
        "qty_invoiced": 100,
        "unit_price": 50000,
        "subtotal": 5000000,
        "tax_amount": 500000,
        "net_unit_price": 50000
      }
    ]
  }
}
```

**Error:** `404 Not Found`

---

### Get by Invoice Number

```http
GET /purchase-invoices/invoice-number/{invoice_number}
Authorization: Bearer {access_token}
```

**Response:** `200 OK` (same as detail response)

---

### Get by Store ID

```http
GET /purchase-invoices/store/{storeId}?status=DRAFT
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `status` (string, optional) - Filter by status

---

### Update Purchase Invoice

Hanya dapat mengubah invoice yang masih berstatus DRAFT.

```http
PUT /purchase-invoices/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "supplier_invoice_number": "INV/SUPP/001-REV",
  "ap_account_id": "uuid",
  "inventory_account_id": "uuid",
  "invoice_date": "2026-05-13",
  "received_date": "2026-05-13",
  "payment_term_days": 30,
  "payment_mode": "TRANSFER",
  "items": [...]
}
```

**Response:** `200 OK`

---

### Submit Invoice

`DRAFT → SUBMITTED`

```http
POST /purchase-invoices/{id}/submit
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "message": "submitted"
}
```

---

### Approve Invoice

`SUBMITTED → VERIFIED`

```http
POST /purchase-invoices/{id}/approve
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "message": "approved"
}
```

---

### Verify Invoice

Sama dengan Approve. `SUBMITTED → VERIFIED`

```http
POST /purchase-invoices/{id}/verify
Authorization: Bearer {access_token}
```

---

### Post Invoice

`VERIFIED → POSTED`. Membuat jurnal akuntansi: Debit Inventory / Credit AP.

```http
POST /purchase-invoices/{id}/post
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "message": "posted"
}
```

---

### Cancel Invoice

Batalkan invoice (hanya sebelum POSTED).

```http
POST /purchase-invoices/{id}/cancel
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "reason": "Invoice rejected due to discrepancy"
}
```

**Response:** `200 OK`
```json
{
  "message": "cancelled"
}
```

---

## Purchase Returns

Retur pembelian untuk mengembalikan barang ke supplier. Alur: `DRAFT → POSTED`.

### Create Purchase Return

```http
POST /purchase-returns
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "purchase_invoice_id": "uuid",
  "return_date": "2026-05-13",
  "notes": "Barang rusak",
  "items": [
    {
      "purchase_invoice_item_id": "uuid",
      "product_id": "uuid",
      "uom_id": "uuid",
      "qty_return": 10,
      "notes": "Rusak saat pengiriman"
    }
  ]
}
```

**Validation:**
- `purchase_invoice_id`: Required
- `return_date`: Required
- `items`: Required, min 1
- `items[].qty_return`: Required, must be > 0

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Purchase return created successfully",
  "data": {
    "id": "uuid",
    "return_number": "PR/2026/05/00001",
    "return_date": "2026-05-13",
    "invoice_number": "PI/2026/05/00001",
    "supplier_name": "Supplier Name",
    "grand_total": 500000,
    "status": "DRAFT"
  }
}
```

---

### Get All Purchase Returns

```http
GET /purchase-returns?page=1&limit=10
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase returns retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "return_number": "PR/2026/05/00001",
      "return_date": "2026-05-13",
      "invoice_number": "PI/2026/05/00001",
      "supplier_name": "Supplier Name",
      "grand_total": 500000,
      "status": "DRAFT"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Purchase Return by ID

```http
GET /purchase-returns/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase return retrieved successfully",
  "data": {
    "id": "uuid",
    "return_number": "PR/2026/05/00001",
    "return_date": "2026-05-13",
    "invoice_number": "PI/2026/05/00001",
    "supplier_name": "Supplier Name",
    "store_name": "Store 001",
    "warehouse_name": "Main Warehouse",
    "subtotal": 500000,
    "grand_total": 500000,
    "status": "DRAFT",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "product_name": "Product Name",
        "uom_code": "PC",
        "qty_return": 10,
        "unit_price": 50000,
        "subtotal": 500000
      }
    ]
  }
}
```

---

### Post Purchase Return

`DRAFT → POSTED`. Mengurangi stok dan membuat jurnal: Debit AP / Credit Inventory.

```http
POST /purchase-returns/{id}/post
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Purchase return posted successfully",
  "data": null
}
```

---

## Expense Vouchers

Voucher pengeluaran/biaya operasional. Alur: `DRAFT → POSTED → VOIDED`.

### Create Expense Voucher

```http
POST /expense-vouchers
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "voucher_date": "2026-05-13",
  "vendor_name": "Vendor Name",
  "payment_type": "CASH",
  "credit_account_id": "uuid",
  "notes": "Biaya operasional",
  "items": [
    {
      "description": "Biaya transportasi",
      "expense_account_id": "uuid",
      "amount": 500000
    }
  ]
}
```

**Validation:**
- `voucher_date`: Required
- `vendor_name`: Required
- `payment_type`: Required, one of `CASH`, `CREDIT`
- `credit_account_id`: Required (UUID)
- `items`: Required, min 1
- `items[].description`: Required
- `items[].expense_account_id`: Required (UUID)
- `items[].amount`: Required, must be > 0

**Response:** `201 Created`

```json
{
  "code": 201,
  "message": "Expense voucher created successfully",
  "data": {
    "id": "uuid",
    "voucher_number": "EV/DRAFT/20260513100405",
    "voucher_date": "2026-05-13",
    "vendor_name": "Vendor Name",
    "payment_type": "CASH",
    "grand_total": 500000,
    "status": "DRAFT"
  }
}
```

---

### Get All Expense Vouchers

```http
GET /expense-vouchers?page=1&limit=10&search=
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Expense vouchers retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "voucher_number": "EV/2026/05/00001",
      "voucher_date": "2026-05-13",
      "vendor_name": "Vendor Name",
      "payment_type": "CASH",
      "grand_total": 500000,
      "status": "DRAFT"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### Get Expense Voucher by ID

```http
GET /expense-vouchers/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Expense voucher retrieved successfully",
  "data": {
    "id": "uuid",
    "voucher_number": "EV/2026/05/00001",
    "voucher_date": "2026-05-13",
    "vendor_name": "Vendor Name",
    "payment_type": "CASH",
    "credit_account_id": "uuid",
    "credit_account_name": "Kas di Bank",
    "grand_total": 500000,
    "status": "DRAFT",
    "notes": "Biaya operasional",
    "items": [
      {
        "id": "uuid",
        "seq_no": 1,
        "description": "Biaya transportasi",
        "expense_account_id": "uuid",
        "expense_account_name": "Biaya Transportasi",
        "amount": 500000
      }
    ]
  }
}
```

---

### Update Expense Voucher

Hanya dapat mengubah voucher yang masih DRAFT.

```http
PUT /expense-vouchers/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "voucher_date": "2026-05-13",
  "vendor_name": "Updated Vendor",
  "payment_type": "CREDIT",
  "credit_account_id": "uuid",
  "notes": "Updated notes",
  "items": [...]
}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Expense voucher updated successfully",
  "data": null
}
```

---

### Post Expense Voucher

`DRAFT → POSTED`. Membuat jurnal akuntansi (debit tiap expense account, credit account dikredit).

```http
POST /expense-vouchers/{id}/post
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Expense voucher posted successfully"
}
```

---

### Cancel Expense Voucher

Batalkan voucher. Jika DRAFT → langsung VOIDED. Jika POSTED → reversal journal.

```http
POST /expense-vouchers/{id}/cancel
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Expense voucher cancelled successfully"
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Planning retrieved successfully",
  "data": {
    "id": "uuid",
    "store_id": "uuid",
    "product_id": "uuid",
    "product_name": "Product Name",
    "product_sku": "PROD001",
    "suggested_qty": 100,
    "status": "PENDING",
    "created_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Roles retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "Admin",
      "permissions": ["uuid1", "uuid2"],
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

#### Get Role by Name

```http
GET /roles/name/{name}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Role retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "Admin",
    "permissions": ["uuid1", "uuid2"],
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
```

---

#### Get Role by ID

```http
GET /roles/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Role retrieved successfully",
  "data": {
    "id": "uuid",
    "name": "Admin",
    "permissions": ["uuid1", "uuid2"],
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Permissions retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "path": "products:view",
      "name": "Lihat Produk",
      "created_at": "2026-05-13T10:00:00Z",
      "updated_at": "2026-05-13T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

#### Get Permission by ID

```http
GET /permissions/{id}
Authorization: Bearer {access_token}
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Permission retrieved successfully",
  "data": {
    "id": "uuid",
    "path": "products:view",
    "name": "Lihat Produk",
    "created_at": "2026-05-13T10:00:00Z",
    "updated_at": "2026-05-13T10:00:00Z"
  }
}
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

#### Sync Permissions

Sinkronisasi permission dari batch data (biasanya dari frontend `permissions_bulk_data.json`).

```http
POST /permissions/sync
Authorization: Bearer {access_token}
Content-Type: application/json

[
  {
    "path": "products:view",
    "name": "Lihat Produk"
  },
  {
    "path": "products:create",
    "name": "Tambah Produk"
  }
]
```

**Response:** `200 OK`

```json
{
  "code": 200,
  "message": "Permissions synced successfully",
  "data": {
    "created": 10,
    "updated": 2,
    "deleted": 0
  }
}
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
