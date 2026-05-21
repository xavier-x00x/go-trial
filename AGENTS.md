# AGENTS.md - Go Trial Project

## Project Overview

Go Fiber REST API with GORM + MySQL. ERP system for retail/store management.

## Build & Run

```bash
go build ./...
go run cmd/api/main.go
```

## Key Conventions

### Entity Pattern

- Entities in `internal/domain/entity/` with `BaseModel` (ID, CreatedAt, UpdatedAt, DeletedAt)
- Status constants defined in entity file (e.g., `POStatusDraft = "DRAFT"`)
- Always include snapshot fields for master data at transaction time

### Repository Pattern

- Interface in `internal/domain/repository/`
- Implementation in `internal/infrastructure/repository/`
- Use `QueryFilter` from `repository` package for pagination (NOT entity.QueryFilter)
- Use `uow.GetTx(ctx, db)` for transaction-aware operations

### UseCase Pattern

- Define interface in `internal/usecase/`
- Use `Config` struct for DI dependencies
- Use repository's `FindAllWithPagination` with `repository.QueryFilter`
- MetaRequest uses `OrderColumn` (not `OrderBy`)
- Always use UOW pattern with Begin/Commit/Rollback for multi-step operations
- Map entity to DTO at the end, never pass raw entities to handlers

### Handler Pattern

- Get `userID` from `c.Locals("user_id")` for authenticated endpoints
- Use `c.Params("id")` for path parameters
- Use `c.Query("key")` for query parameters
- Return formatted responses, not raw entities

### DTO Pattern

- Separate Request and Response DTOs
- Use proper validation tags (`validate:"required"`)
- Use `decimal.Decimal` from `github.com/shopspring/decimal` for money

### Journal Entry Pattern

- For POST actions, create `JournalEntry` with `JournalSourcePurchasePayment`, `JournalSourcePurchaseInvoice`, etc.
- Always set `SourceDocumentID` and `SourceDocumentNo` for audit trail
- Reverse with opposite debit/credit amounts and prepend "VOID - " to description

### Route Registration

- Handlers registered in `internal/delivery/http/route/route.go`
- Add setup function (e.g., `setupPurchasePayments`), wire in `setupFinanceRoutes` or similar
- Pattern: `r.Group("/resource")` then chain handlers

### Permission Sync (permissions_bulk_data.json)

File `permissions_bulk_data.json` adalah master list semua permission di sistem. Digunakan frontend untuk sync ke backend via `POST /api/permissions/sync`.

**Cara membuat/update:**
1. Buka `internal/delivery/http/route/route.go`
2. Identifikasi semua resource group dan action endpoints
3. Setiap resource butuh permission: `view`, `create`, `update`, `delete`
4. Resource dengan workflow (submit/approve/post/cancel) butuh permission tambahan sesuai action-nya
5. Format path: `<resource>:<action>` (kebab-case, colon separator)
6. Format name: Bahasa Indonesia, diawali huruf kapital, gunakan istilah yang sudah dipakai di route

**Aturan penamaan path:**
- Resource: kebab-case plural (purchase-orders, goods-receipts, master-data)
- Action: kebab-case (create-with-invoice)
- Contoh: `purchase-orders:view`, `goods-receipts:create-with-invoice`

**Saat nambah resource baru di route.go:**
1. Tambahkan permission entries di `permissions_bulk_data.json`
2. Jangan lupa tambah `check("resource:action")` middleware di route

## Common Pitfalls

1. **QueryFilter type mismatch**: Repository uses `repository.QueryFilter`, not `entity.QueryFilter`
2. **MetaRequest field names**: Use `OrderColumn`, not `OrderBy`
3. **Unused imports**: Always remove unused imports after refactoring
4. **Transaction rollback**: Always use defer with rollback in same function
5. **UUID parsing**: Parse userID from string using `uuid.Parse()` for database operations
6. **User.ID is string**: User entity uses string ID, need `uuid.Parse()` for JournalEntry.PostedByID
7. **Entity value types**: Some relations are value types (not pointer), check before nil comparison

## Architecture Layers

- `cmd/` - Entry points (main.go, gen_po, import_products)
- `internal/config/` - Configuration loading
- `internal/delivery/http/` - HTTP handlers, DTOs, routes
- `internal/domain/entity/` - Domain entities
- `internal/domain/repository/` - Repository interfaces
- `internal/infrastructure/repository/` - Repository implementations
- `internal/infrastructure/uow/` - Unit of Work for transactions
- `internal/usecase/` - Business logic
- `internal/registry/` - Dependency injection wiring
- `pkg/` - Shared utilities (response, validator, jwt)

## Response Agent

- Setiap Jawaban Gunakan Bahasa Indonesia
