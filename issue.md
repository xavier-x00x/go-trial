# Implementasi Validasi Sisa Kuantitas & Over-Receiving Override pada Goods Receipt

## Deskripsi Fitur
Saat ini, sistem mengizinkan pembuatan multi-DRAFT Goods Receipt (GR) untuk satu Purchase Order (PO) yang sama tanpa memvalidasi sisa kuantitas barang secara keseluruhan. Hal ini berpotensi menyebabkan *over-receiving* (penerimaan berlebih) yang tidak disengaja.

Issue ini bertujuan untuk mengimplementasikan:
1. **Validasi Sisa Kuantitas yang Fleksibel**: Mencegah pembuatan GR jika total barang yang akan diterima melebihi pesanan.
2. **Fitur Over-Receiving Override (Bypass)**: Memungkinkan penerimaan berlebih *hanya jika* diotorisasi menggunakan PIN dari Supervisor/Kepala Gudang.
3. **Filter PO**: Memfilter daftar PO yang bisa dibuatkan GR berdasarkan sisa kuantitas.

---

## Aturan Bisnis (Business Rules)
1. **Rumus Sisa Kuantitas**: 
   `Sisa Kuantitas = Qty Dipesan - Qty Diterima (yang sudah Confirm) - Qty di GR Draft lain`
2. **Otorisasi Over-Receiving**: PIN yang digunakan untuk *bypass* harus milik user dengan role Supervisor atau Kepala Gudang.
3. **Standar Kode**: Wajib menggunakan *Clean Architecture* dan pattern *Unit of Work (UOW)* untuk transaksi database.

---

## Tahapan Implementasi (Checklist untuk Programmer)

### 1. Persiapan Database & Entity (Domain Layer)
- [ ] Tambahkan field `IsOverReceivedOverride bool` pada entity `GoodsReceipt`.
- [ ] Tambahkan field `OverrideApprovedByID *uuid.UUID` pada entity `GoodsReceipt` beserta relasinya ke entity `User`.
- [ ] Jalankan/buat file *migration* untuk menambahkan kedua kolom tersebut ke tabel `goods_receipts` di database.

### 2. Modifikasi Data Transfer Object (Delivery/DTO Layer)
- [ ] Buka file `goods_receipt_dto.go`.
- [ ] Tambahkan field opsional `OverridePIN *string` pada struct `CreateGoodsReceiptRequest` dan `UpdateGoodsReceiptRequest`.

### 3. Penambahan Method di Repository
- [ ] Buka `GoodsReceiptRepository`.
- [ ] Tambahkan method untuk menghitung total kuantitas barang yang masih *menggantung* di GR berstatus `DRAFT` untuk suatu PO Item tertentu.
  *(Contoh: `GetTotalDraftQtyByPOItemID(ctx, poItemID) decimal.Decimal`)*.
- [ ] Modifikasi *query* pencarian PO (misal `FindPendingByStoreID` di `PurchaseOrderRepository`) agar **hanya** mengembalikan PO yang memiliki `Sisa Kuantitas > 0`. (Bisa dilakukan via filter query atau perhitungan di level repository/usecase).

### 4. Modifikasi Logika Bisnis (UseCase Layer)
- [ ] Buka `GoodsReceiptUseCase` pada method `Create` dan `Update`.
- [ ] Implementasikan logika perhitungan Sisa Kuantitas:
  - Looping setiap *request item*.
  - Ambil `QtyOrdered` dan `QtyReceived` dari `PO Item`.
  - Ambil total `Qty Draft` dari repository (langkah 3).
  - Hitung: `Sisa = QtyOrdered - QtyReceived - QtyDraft`.
- [ ] Terapkan logika pengecekan **Over-Receiving**:
  - **JIKA** `Qty Request > Sisa` **DAN** `OverridePIN` kosong:
    - *Return error* khusus, misalnya: `errors.New("ERR_OVER_RECEIVE_NEEDS_PIN")` (Nanti di-mapping ke HTTP 400/403 di Handler).
  - **JIKA** `Qty Request > Sisa` **DAN** `OverridePIN` terisi:
    - Verifikasi PIN tersebut ke database `User` (Pastikan user pemilik PIN memiliki akses Supervisor/Kepala Gudang).
    - Jika PIN salah: *Return error* `Invalid PIN`.
    - Jika PIN benar: Set `IsOverReceivedOverride = true` dan `OverrideApprovedByID = ID Pemilik PIN`. Lanjutkan proses simpan.
  - **JIKA** `Qty Request <= Sisa`:
    - Lanjutkan proses penyimpanan data normal.
- [ ] Pastikan seluruh proses penyimpanan data dibungkus dalam `uow.Do` (Unit of Work) untuk menjaga integritas transaksi.

### 5. Penanganan Error di Controller (Handler Layer)
- [ ] Buka `GoodsReceiptHandler`.
- [ ] Tangkap error `ERR_OVER_RECEIVE_NEEDS_PIN` dari UseCase.
- [ ] Format response JSON agar mengembalikan status code yang sesuai (misal 403 Forbidden atau 400 Bad Request) lengkap dengan kode unik agar Frontend mudah mendeteksinya.

---

## Catatan Konsep Alur Frontend (Untuk Integrasi)
*(Informasi ini agar Backend Dev paham bagaimana datanya akan dikonsumsi)*
1. User menginput jumlah barang masuk di UI dan klik Submit.
2. Jika kuantitas melebihi batas, Backend mengembalikan response error `ERR_OVER_RECEIVE_NEEDS_PIN`.
3. Frontend menangkap error tersebut, lalu *tidak* langsung menggagalkan proses, melainkan memunculkan **Pop-up "Penerimaan Berlebih! Masukkan PIN Supervisor:"**.
4. Supervisor memasukkan PIN di perangkat user tersebut.
5. Frontend melakukan resubmit request awal, namun kali ini menyisipkan payload `OverridePIN`.
6. Backend memvalidasi PIN, jika sukses, GR tersimpan.
