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

## 2026-05-23

- Added `github.com/xuri/excelize/v2` dependency - Library Excel untuk import/export
- Modified `internal/domain/repository/chart_of_account_repository.go` - Tambah method `BulkCreate` di interface
- Modified `internal/infrastructure/repository/chart_of_account_repository_impl.go` - Implementasi `BulkCreate`
- Modified `internal/delivery/http/dto/chart_of_account_dto.go` - Tambah DTO `AccountImportRow`, `AccountImportResult`, `ImportRowError`
- Modified `internal/usecase/chart_of_account_usecase.go` - Tambah method `Import` (parse & bulk insert dari Excel) dan `GenerateTemplate` (download template Excel)
- Modified `internal/delivery/http/handler/coa_handler.go` - Tambah handler `Import` (POST /accounts/import) dan `DownloadTemplate` (GET /accounts/import/template)
- Modified `internal/delivery/http/route/route_finance.go` - Daftarkan route import dan template
- Modified `permissions_bulk_data.json` - Tambah permission `accounts:import`
- Modified `internal/domain/entity/chart_of_account.go` - Tambah field `ParentID` (*uuid.UUID) dan relasi `Parent` untuk hierarki akun
- Modified `internal/delivery/http/dto/request_dto.go` - Tambah field `ParentID` di `CreateCOARequest` dan `UpdateCOARequest`
- Modified `internal/delivery/http/dto/response_dto.go` - Tambah field `ParentID` di `ChartOfAccountResponse` dan DTO `ChartOfAccountTreeResponse`
- Modified `internal/domain/repository/chart_of_account_repository.go` - Tambah method `FindByParentID` di interface
- Modified `internal/infrastructure/repository/chart_of_account_repository_impl.go` - Implementasi `FindByParentID`
- Modified `internal/usecase/chart_of_account_usecase.go` - Handle `ParentID` di `Create`, `Update`, `Import`; tambah error `ErrCOAParentNotFound`; tambah method `GetTree` untuk hierarki
- Modified `internal/delivery/http/handler/coa_handler.go` - Tambah handler `GetTree` (GET /accounts/tree)
- Modified `internal/delivery/http/route/route_finance.go` - Daftarkan route `GET /accounts/tree`
- Modified `API_DOCUMENTATION.md` - Update dokumentasi Chart of Accounts: tambah `parent_id`, endpoint import, template, tree

## 2026-05-24

- Fixed `internal/infrastructure/repository/base_repository_impl.go` - Perbaiki bug shared pointer `*gorm.DB` pada `PaginateAndFilter` yang menyebabkan `total = 0` di page 2+ karena `Offset` dari pagination bocor ke query count total
- Created `docs/FRONTEND_MASTER_DATA_PROPOSALS.md` - Catatan implementasi Master Data Proposal untuk frontend (API endpoints, payload format per entity type, page structure, permission mapping)