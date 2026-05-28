# Dokumentasi & Alur Kerja Lengkap: Master Data Proposals (Backend - Frontend)

Dokumen ini merangkum alur kerja terintegrasi dari fitur **Master Data Proposals**, mulai dari perancangan sisi backend (API, proses, validasi) hingga implementasinya di sisi frontend (halaman, form, pengelolaan JSON payload).

---

## 1. Konsep Utama: Maker-Checker

Sistem menerapkan prinsip **Maker-Checker** untuk semua perubahan pada Master Data (seperti Produk, Supplier, Pajak, Chart of Account, dll).
Perubahan **tidak akan langsung disimpan** ke tabel master data utama, melainkan masuk ke sebuah tabel `master_data_proposals` terlebih dahulu untuk melalui proses peninjauan (review) dan persetujuan (approval).

> [!NOTE]
> **Ilustrasi Alur:**
> Maker (Staf/User) $\rightarrow$ **CREATE Proposal** $\rightarrow$ [Status: PENDING] $\rightarrow$ Checker (Manajer) $\rightarrow$ **REVIEW** $\rightarrow$ [Status: APPROVED] $\rightarrow$ **EXECUTE** (sistem menulis ke tabel asli)

### Flowchart Status

- `PENDING`: Proposal baru dibuat, menunggu review.
- `APPROVED`: Proposal telah disetujui, siap dieksekusi.
- `REJECTED`: Proposal ditolak dengan alasan tertentu. Bisa diperbarui kembali menjadi `PENDING`.
- `EXECUTED`: Proposal yang sudah dieksekusi ke tabel asli (Telah selesai).

---

## 2. Entitas & Prefiks

Sistem mendukung pengajuan (proposal) untuk beberapa entitas master data. Setiap proposal memiliki **Reference Number** unik dengan format `{PREFIX}/{YYMM}/{NNNNN}`.

| Tipe Entitas             | Prefix | Keterangan                |
| :----------------------- | :----- | :------------------------ |
| `PRODUCT`                | `PRD`  | Master Data Barang        |
| `PRODUCT_PRICE`          | `PRC`  | Harga jual per price list |
| `PRODUCT_UOM_CONVERSION` | `PUC`  | Konversi satuan           |
| `SUPPLIER`               | `SUP`  | Pemasok                   |
| `PRODUCT_SUPPLIER`       | `PSP`  | Kontrak/Harga pemasok     |
| `CHART_OF_ACCOUNT`       | `COA`  | Kode akun akuntansi       |
| `TAX`                    | `TAX`  | Pengaturan pajak          |

---

## 3. Implementasi Backend

### 3.1. Struktur Data Proposal

Satu proposal (Header) bisa memiliki banyak item pengajuan (Items).
Data perubahan sebenarnya disimpan di level `Item` di dalam kolom **`payload_json`**.

- Jika pengajuan `CREATE`, `payload_json` berisi data lengkap objek baru.
- Jika pengajuan `UPDATE`, `payload_json` berisi field yang diedit, dan sistem menyimpan data awal di kolom **`snapshot_json`**.
- Jika pengajuan `DELETE`, `payload_json` dapat kosong atau hanya memuat instruksi penghapusan.

### 3.2. API Endpoints

Semua endpoint berawalan `/api/master-data`.

**Daftar Endpoint:**

- `GET /master-data` : Mengambil list proposal.
- `GET /master-data/pagination` : Pagination dan filter list proposal (by status, entity_type).
- `GET /master-data/{id}` : Mengambil detail proposal dan items-nya.
- `POST /master-data` : **Create** proposal baru (Maker).
- `PUT /master-data/{id}` : **Update** proposal (Hanya bisa saat status `PENDING` dan oleh _Proposer_).
- `POST /master-data/{id}/review` : **Review** proposal (Approve / Reject oleh Checker).
- `POST /master-data/{id}/execute` : **Execute** proposal (Memproses data JSON ke tabel entitas aslinya).
- `POST /master-data/bulk/product-supplier` : Membuat proposal bulk untuk Product-Supplier.

### 3.3. Akses & Permission

Akses API dan tombol di UI dijaga melalui rules permission (contoh `permissions_bulk_data.json`):

- `master-data:view` : Akses daftar dan detail proposal.
- `master-data:create` : Mengajukan proposal baru.
- `master-data:update` : Mengubah proposal miliknya (`PENDING`).
- `master-data:review` : Hak Checker untuk eksekusi Approve/Reject.
- `master-data:execute` : Hak untuk menjalankan Execute proposal.

---

## 4. Implementasi Frontend

Terdapat struktur halaman dan pola form yang seragam karena semua 7 tipe entitas menggunakan alur yang identik.

### 4.1. Struktur Halaman (Pages)

1. **Halaman List (`/master-data`)**
   - Menampilkan tabel Data Proposal dengan kolom: `Reference Number`, `Entity Type`, `Action Type` (Create/Update/Delete), `Status`, dan `Created At`.
   - Filter lengkap: Status, Pencarian, Rentang Tanggal, dan Tipe Entitas.

2. **Halaman Detail (`/master-data/:id`)**
   - Menampilkan **Header Info** (Status, Pembuat, Review Notes, dll).
   - Menampilkan **Items Table** yang parsing nilai dari `payload_json` menjadi tampilan tabel yang mudah dibaca (misal menggunakan formater JSON).
   - Menampilkan **Diff/Perbandingan** (Before vs After) menggunakan data `snapshot_json` untuk proposal tipe `UPDATE`.
   - Tombol Aksi Dinamis:
     - Tampil tombol _Edit/Hapus_ jika user adalah pengaju & status `PENDING`.
     - Tampil tombol _Approve/Reject_ jika user punya akses Review.
     - Tampil tombol _Execute_ jika status `APPROVED`.

3. **Halaman Create / Edit Proposal**
   - Menggunakan sistem _Factory Pattern_ di mana form di-render dinamis sesuai dengan _Entity Type_ yang dipilih.
   - Form mengelola konversi input user $\rightarrow$ String JSON (`payload_json`) saat submit ke Backend.

### 4.2. Komponen Reusable

Untuk menyederhanakan kode Vue/Nuxt, dibuat komponen-komponen yang dapat digunakan di semua entitas:

- `ProposalListTable`: Tabel data master data proposal.
- `ProposalDetailCard`: Header pembungkus detail proposal.
- `ProposalFormFactory`: Komponen dinamis pemilih Form (contoh render `ProductForm.vue` vs `SupplierForm.vue`).
- `ReviewModal`: Modal konfirmasi dengan _textarea_ mandatory untuk catatan penolakan (`Reject`).
- `StatusBadge` & `EntityTypeBadge`: Pewarnaan seragam untuk status dan tipe di tabel.

### 4.3. Penanganan Data JSON di Frontend

> [!IMPORTANT]
> Saat merender nilai Payload/Snapshot di tabel Items, Frontend harus me-restore JSON tersebut menjadi UI yang rapi. Jangan tampilkan raw JSON string ke user.

**Contoh Pattern Handling JSON:**

```javascript
// POST / UPDATE (Submit data)
const payloadString = JSON.stringify(formDataObject);

// Tampil Detail (Render data)
const itemDisplayData = JSON.parse(item.payload_json);
```

---

## 5. Ringkasan Edge Cases & Validasi Bisnis

1. **Pembaruan Terbatas:** Hanya "Proposer" (Staf pembuat) yang dapat merubah (Update) proposal miliknya jika ditolak (`REJECTED`) atau masih `PENDING`.
2. **Review Validation:** Sistem menolak proses review jika proposal sedang dalam state selain `PENDING`.
3. **Execute Validation:** Proses eksekusi ke tabel original tidak akan terjadi jika proposal belum berstatus `APPROVED`.
4. **Reference Numbers:** Nomor Referensi digenerate secara otomatis di backend untuk mencegah bentrok/duplikasi.

Dengan arsitektur Maker-Checker ini, kita mendapat audit trail yang baik atas setiap perubahan tabel vital di aplikasi ERP.
