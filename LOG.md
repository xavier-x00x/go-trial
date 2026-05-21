# LOG.md

## 2026-05-14

- Created `internal/domain/repository/purchase_payment_repository.go` - Interface repository untuk PurchasePayment
- Created `internal/infrastructure/repository/purchase_payment_repository_impl.go` - Implementasi repository PurchasePayment
- Created `internal/infrastructure/repository/purchase_invoice_repository_impl.go` - Implementasi repository PurchaseInvoice (missing sebelumnya)
- Created `internal/usecase/purchase_payment_usecase.go` - Use case dengan logika Create, Post (create journal entry), dan Void (reversal journal)
- Created `internal/delivery/http/handler/purchase_payment_handler.go` - HTTP handler untuk PurchasePayment
- Modified `internal/registry/finance.go` - Tambah DI wiring untuk PurchasePayment
- Modified `internal/delivery/http/route/route.go` - Tambah route registration untuk PurchasePayment endpoints
- Modified `AGENTS.md` - Tambah pattern Journal Entry, Route Registration, dan common pitfalls baru

## 2026-05-20

- Created `permissions_bulk_data.json` - File JSON daftar lengkap semua permissions untuk dikirim dari frontend
- Modified `internal/delivery/http/dto/role_permission_dto.go` - Tambah DTO `SyncPermissionItem` dan `SyncPermissionsRequest`
- Modified `internal/delivery/http/handler/role_permission_handler.go` - Tambah handler `SyncPermissions`
- Modified `internal/usecase/role_permission_usecase.go` - Tambah method `SyncPermissions` di use case
- Modified `internal/domain/repository/role_permission_repository.go` - Tambah method `FindByPaths` di interface
- Modified `internal/infrastructure/repository/role_permission_repository_impl.go` - Implementasi `FindByPaths`
- Modified `internal/delivery/http/route/route.go` - Daftarkan route `POST /permissions/sync`

## 2026-05-21

- Modified `internal/delivery/http/handler/auth_handler.go` - Tambah validator di handler UpdateUser
- Modified `pkg/validator/validator.go` - Perbaiki fungsi toSnakeCase agar StoreID menjadi store_id (bukan store_i_d)