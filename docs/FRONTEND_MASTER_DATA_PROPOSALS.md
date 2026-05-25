# Frontend — Master Data Proposals

## 1. Konsep: Maker-Checker

Setiap perubahan data master (PRODUCT, SUPPLIER, TAX, dll.)
**tidak langsung** mengubah tabel asli, melainkan melalui workflow
**Maker-Checker**:

```
Maker (Staf)         Checker (Manajer)
      |                     |
      v                     v
  CREATE proposal ──> REVIEW ──> APPROVED ──> EXECUTE (otomatis atau manual)
                        |
                        v
                    REJECTED (dengan alasan)
```

### Alur Lengkap

| Step | Actor | Aksi | Status Proposal |
|------|-------|------|-----------------|
| 1 | Maker (Staf) | Ajukan perubahan via `POST /master-data` | `PENDING` |
| 2 | Checker (Manajer) | Lihat antrian pending | `PENDING` |
| 3 | Checker | `POST /master-data/:id/review` (APPROVE/REJECT) | `APPROVED` / `REJECTED` |
| 4 (a) | System | `POST /master-data/:id/execute` — tulis ke tabel asli | `EXECUTED` (opsional) |
| 4 (b) | Maker | Jika REJECTED, bisa diperbaiki (UPDATE) lalu tunggu review lagi | `PENDING` lagi |

---

## 2. Entity Types & Prefix

| Entity Type | Prefix | Deskripsi | Ada di fitur apa |
|-------------|--------|-----------|------------------|
| `PRODUCT` | PRD | Barang / produk | Master Barang |
| `PRODUCT_PRICE` | PRC | Harga barang per price list | Harga Jual |
| `PRODUCT_UOM_CONVERSION` | PUC | Konversi satuan barang | Satuan Barang |
| `SUPPLIER` | SUP | Pemasok | Master Supplier |
| `PRODUCT_SUPPLIER` | PSP | Kontrak barang-supplier | Kontrak Pembelian |
| `CHART_OF_ACCOUNT` | COA | Akun akuntansi | Chart of Accounts |
| `TAX` | TAX | Pajak (PPN dll) | Master Pajak |

**Format Reference Number:** `{PREFIX}/{YYMM}/{NNNNN}`

Contoh: `PRD/2605/00007`, `SUP/2605/00003`

---

## 3. API Endpoints

Base path: `/api/master-data`

### List & Filter

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/master-data` | Semua proposal (tanpa pagination) |
| `GET` | `/master-data/pagination?page=1&limit=20&search=&status=PENDING&entity_type=PRODUCT&order_column=created_at&order_dir=desc` | Pagination + filter |
| `GET` | `/master-data/pending` | Khusus proposal PENDING |
| `GET` | `/master-data/entity/{entityType}?status=PENDING` | Filter by entity type |
| `GET` | `/master-data/group/{groupId}` | Proposal dalam 1 grup (bulk) |

### CRUD Proposal

| Method | Path | Aksi |
|--------|------|------|
| `POST` | `/master-data` | **Create** proposal baru |
| `GET` | `/master-data/{id}` | **Detail** proposal (header + items) |
| `PUT` | `/master-data/{id}` | **Update** proposal yang masih PENDING |

### Review & Execute

| Method | Path | Aksi |
|--------|------|------|
| `POST` | `/master-data/{id}/review` | **Review** (APPROVE / REJECT) |
| `POST` | `/master-data/{id}/execute` | **Execute** (proses ke tabel asli) |

### Bulk

| Method | Path | Aksi |
|--------|------|------|
| `POST` | `/master-data/bulk/product-supplier` | Bulk link product-supplier |

---

## 4. Struktur Halaman (Pages)

### a. Proposal List Page `/master-data`

**Tabel dengan kolom:**

| Kolom | Keterangan |
|-------|------------|
| Reference Number | `PRD/2605/00001` |
| Entity Type | Badge: PRODUCT, SUPPLIER, dll |
| Action Type | Badge warna: CREATE (hijau), UPDATE (biru), DELETE (merah) |
| Total Items | Jumlah item dalam 1 proposal |
| Status | Badge: PENDING (kuning), APPROVED (hijau), REJECTED (merah) |
| Proposed By | Nama pengaju |
| Reason | Alasan pengajuan (tooltip) |
| Created At | Tanggal |
| Actions | Lihat Detail, Review (jika checker) |

**Filter:**

- Search (by reference_number)
- Status dropdown: ALL, PENDING, APPROVED, REJECTED
- Entity Type dropdown: ALL, PRODUCT, SUPPLIER, dll
- Date range (opsional)

### b. Proposal Detail Page `/master-data/:id`

Dua seksi:

**Header Info:**

| Field | Value |
|-------|-------|
| Reference Number | `PRD/2605/00001` |
| Entity Type | PRODUCT |
| Action Type | CREATE |
| Status | PENDING |
| Status Timeline | Created → Pending → Approved → Executed |
| Proposed By | Nama staf |
| Reason | "Tambah barang baru Q2" |
| Review Notes | (jika sudah direview) |
| Reviewed By | (jika sudah direview) |
| Reviewed At | (jika sudah direview) |

**Items Table:**

| # | Entity ID | Payload (parsed) | Snapshot |
|---|-----------|-------------------|----------|
| 1 | `-` (CREATE) | SKU: BRG001, Name: Minuman 500ml | - |
| 2 | `-` (CREATE) | SKU: BRG002, Name: Snack 100gr | - |

**Actions (conditional):**

- Jika status PENDING dan user = proposer: **Edit**, **Hapus**
- Jika status PENDING dan user punya permission `master-data:review`: **Approve** / **Reject**
- Jika status APPROVED dan user punya permission `master-data:execute`: **Execute**

### c. Create/Edit Proposal Page

Satu halaman yang menyesuaikan form berdasarkan `entity_type` yang dipilih.

**Step 1: Pilih Entity & Action**

Dropdown:
- Entity Type: PRODUCT / SUPPLIER / TAX / CHART_OF_ACCOUNT / PRODUCT_PRICE / PRODUCT_UOM_CONVERSION / PRODUCT_SUPPLIER
- Action Type: CREATE / UPDATE / DELETE (tampilkan opsi sesuai)

**Step 2: Isi Alasan**

Textarea (opsional)

**Step 3: Isi Items**

Jika `CREATE`: Form dinamis untuk input data baru (1 form = 1 item, bisa tambah banyak).

Jika `UPDATE`: Picker/search entity yang sudah ada → ambil data lama → tampilkan form pre-filled.

Jika `DELETE`: Picker/search entity yang sudah ada → konfirmasi.

**Step 4: Submit**

### d. Review Dialog (Modal)

Tombol Approve (hijau) dan Reject (merah) di detail page.

Jika **Reject** → wajib isi `review_notes` alasan penolakan.

Jika **Approve** → langsung submit.

---

## 5. Permissions

| Permission | Pages yang dilindungi |
|------------|----------------------|
| `master-data:view` | List, Detail |
| `master-data:create` | Halaman Create |
| `master-data:update` | Halaman Edit |
| `master-data:review` | Tombol Approve/Reject |
| `master-data:execute` | Tombol Execute |

**Catatan:** Saat ini backend **belum** memasang middleware `check("master-data:*")` di route, tapi frontend tetap harus sembunyikan tombol berdasarkan permission user.

---

## 6. Seragam Frontend untuk Semua Entity Type

Karena hampir semua entity type menggunakan pola yang sama (7 entity), disarankan:

### a. Komponen Reusable

| Komponen | Fungsi |
|----------|--------|
| `ProposalListTable` | Tabel proposal dengan filter |
| `ProposalDetailCard` | Header + item list |
| `ProposalFormFactory` | Render form dinamis berdasarkan entity type |
| `ReviewModal` | Modal approve/reject |
| `StatusBadge` | Badge warna untuk status |
| `EntityTypeBadge` | Badge untuk entity type |
| `ActionTypeBadge` | Badge untuk action type |

### b. Factory Pattern untuk Form

Buat file mapping entity type → form component:

```
entity-forms/
├── ProductForm.vue
├── SupplierForm.vue
├── TaxForm.vue
├── CoaForm.vue
├── ProductPriceForm.vue
├── ProductUomForm.vue
└── ProductSupplierForm.vue
```

Lalu `ProposalFormFactory` me-render form yang sesuai berdasarkan `entity_type`.

### c. Payload JSON Handler

Setiap item proposal menyimpan data sebagai **JSON string** di `payload_json`.

Pattern:

```typescript
// Konversi object → JSON string sebelum POST
const payload = JSON.stringify(formData);
// Atau parse JSON string → object untuk display
const data = JSON.parse(item.payload_json);
```

---

## 7. Payload JSON per Entity Type

### 7.1 PRODUCT

**CREATE:**
```json
{
  "sku": "BRG001",
  "barcode": "8991234567890",
  "name": "Minuman Kemasan 500ml",
  "category_id": "uuid",
  "base_uom_id": "uuid",
  "is_stockable": true,
  "length": 10.5,
  "width": 5.0,
  "height": 15.0,
  "weight": 500.0,
  "is_stackable": true,
  "max_stack_layer": 8
}
```

**UPDATE:** field sama, semua opsional.

### 7.2 PRODUCT_PRICE

**CREATE:**
```json
{
  "price_list_id": "uuid",
  "product_id": "uuid",
  "uom_id": "uuid",
  "markup_pct": 20.00,
  "sell_price": 15000.00,
  "discount_pct": 5.00
}
```

### 7.3 PRODUCT_UOM_CONVERSION

**CREATE:**
```json
{
  "product_id": "uuid",
  "uom_id": "uuid",
  "conversion_rate": 24.000,
  "barcode": "8991234567890",
  "length": 40.0,
  "width": 30.0,
  "height": 25.0,
  "weight": 5000.0,
  "is_stackable": true,
  "max_stack_layer": 6
}
```

### 7.4 SUPPLIER

**CREATE:**
```json
{
  "code": "SUP001",
  "name": "PT Sumber Makmur",
  "contact_person": "Budi",
  "contact_phone": "08123456789",
  "phone_number": "021-12345678",
  "email": "budi@email.com",
  "preferred_notification_method": "WHATSAPP",
  "address": "Jl. Merdeka No. 123",
  "tax_reg_number": "01.234.567.8-901.000",
  "supplier_category_id": "uuid",
  "is_pkp": true,
  "payment_term_days": 30,
  "payment_mode": "TRANSFER",
  "min_order_amount": 500000,
  "bank_name": "BCA",
  "bank_account": "1234567890",
  "bank_account_name": "PT Sumber Makmur",
  "ap_account_id": "uuid"
}
```

### 7.5 PRODUCT_SUPPLIER

**CREATE:**
```json
{
  "product_id": "uuid",
  "supplier_id": "uuid",
  "store_id": null,
  "supplier_sku": "SBRG-XYZ-001",
  "is_primary": true,
  "is_consignment": false,
  "is_returnable": true,
  "default_lead_time_days": 3,
  "offered_price": 12000.00,
  "min_order_qty": 10.000
}
```

### 7.6 CHART_OF_ACCOUNT

**CREATE:**
```json
{
  "account_code": "1110.01",
  "name": "Kas di Tangan",
  "account_type": "ASSET",
  "normal_balance": "DEBIT"
}
```

**Enum values:**
- `account_type`: `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`
- `normal_balance`: `DEBIT`, `CREDIT`

### 7.7 TAX

**CREATE:**
```json
{
  "name": "PPN 11%",
  "rate_percentage": 11.00,
  "tax_account_id": "uuid"
}
```

---

## 8. Tips Implementasi

### a. JSON Display

Item `payload_json` adalah string — perlu di-parse ke object untuk ditampilkan sebagai tabel:

```typescript
interface ProposalItemDisplay {
  seqNo: number;
  fields: { key: string; value: string }[];
}
```

### b. Snapshot

Saat UPDATE, backend akan mengisi `snapshot_json` dengan data **sebelum** perubahan — ini berguna untuk menampilkan diff (sebelum vs sesudah) saat review.

### c. Bulk Product-Supplier

Halaman khusus: pick 1 supplier + pilih banyak produk → submit sebagai 1 bulk request. Backend akan memecah jadi banyak proposal (1 per item).

### d. Form State

Untuk masing-masing action:

| Action | Entity ID | Payload |
|--------|-----------|---------|
| CREATE | jangan kirim | form CREATE (isi semua) |
| UPDATE | wajib diisi | form UPDATE (isi yg berubah saja, sisanya dari snapshot) |
| DELETE | wajib diisi | `{}` atau tidak perlu payload berarti |

### e. Status Flow

```
PENDING ──> APPROVED ──> EXECUTED
   │
   └──> REJECTED ──> (bisa diedit jadi PENDING lagi)
```

Update hanya bisa:
- Oleh **proposer** (staf yang membuat)
- Status **PENDING**
- Semua items akan **diganti total** (bukan patch, tapi replace)

### f. Edge Cases

1. **User bukan proposer coba update** → backend return error
2. **Review proposal yang bukan PENDING** → backend return error
3. **Execute proposal yang bukan APPROVED** → backend return error
4. **Reference number sudah dipakai** → backend auto-generate unik, frontend tidak perlu urus
