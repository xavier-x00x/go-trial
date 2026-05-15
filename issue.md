# Fitur Update untuk Goods Receipt (GR)

## Deskripsi Singkat
Fitur ini bertujuan untuk memungkinkan pengguna mengubah (update) dokumen `GoodsReceipt` yang statusnya masih `DRAFT` (belum dikonfirmasi/diterima). Jika status sudah `CONFIRMED` atau dibatalkan, maka penerimaan barang tidak boleh diubah lagi. Proses update ini akan mencakup perubahan pada data header (seperti referensi pengiriman, tanggal terima, dan catatan) serta melakukan *replace* (hapus dan buat ulang) pada semua baris item (`Items`) di dalam dokumen tersebut.

## Persyaratan Utama (Requirements)
1. Sesuai dengan standar **Clean Architecture** (DTO -> Handler -> UseCase -> Repository).
2. Wajib menggunakan **Unit of Work (UOW)** karena proses melibatkan *delete & insert* item dalam satu transaksi untuk mencegah data *corrupt*.
3. Hanya GR berstatus `DRAFT` yang diizinkan untuk di-update.
4. Total nilai seperti `Subtotal`, `TaxAmount`, `GrandTotal`, dan sebagainya harus dihitung ulang berdasarkan item yang baru.

---

## Tahapan Implementasi untuk Programmer / AI

### 1. Delivery Layer (DTO)
- Buka file `internal/delivery/http/dto/goods_receipt_dto.go`.
- Buat struct baru bernama `UpdateGoodsReceiptRequest`.
- Isinya mirip dengan `CreateGoodsReceiptRequest`, **TETAPI** pastikan untuk **menghilangkan** field relasi master seperti `PurchaseOrderID` dan `WarehouseID`. Relasi ini dikunci agar tidak merusak konsistensi data yang sudah berjalan.
- Field yang boleh di-update antara lain: `ReceiptDate`, `DeliveryNoteNo`, `Notes`, dan list `Items`.

### 2. Repository Layer (Hapus Item Lama)
- Buka antarmuka (interface) di `internal/domain/repository/goods_receipt_repository.go`.
- Tambahkan method baru: `DeleteItemsByGoodsReceiptID(ctx context.Context, grID string) error`.
- Buka file implementasinya di `internal/infrastructure/repository/goods_receipt_repository_impl.go` dan tuliskan kodenya menggunakan `uow.GetTx(ctx, r.db)` untuk menghapus data di tabel `goods_receipt_items` berdasarkan `goods_receipt_id`.

### 3. UseCase Layer (Logika Bisnis & Transaksi UOW)
- Buka `internal/usecase/goods_receipt_usecase.go`.
- Tambahkan method `Update` pada interface dan struct implementasinya.
- **Langkah logika fungsi Update:**
  1. Cari dokumen GR berdasarkan ID. Jika tidak ada, return `ErrGoodsReceiptNotFound`.
  2. Validasi status dokumen. Pastikan statusnya adalah `DRAFT`. Jika bukan, tolak perubahannya.
  3. Lakukan *mapping* data dari DTO ke entitas GR (timpa header lama dengan data baru).
  4. Lakukan iterasi pada list `Items` yang baru, dan hitung ulang semua nominal (kalkulasi Subtotal, harga, dsb.) layaknya fungsi *Create*. Set ID pada tiap item baru.
  5. Buka transaksi database dengan UOW (`txCtx, err := u.uow.Begin(ctx)`).
  6. Panggil `repo.DeleteItemsByGoodsReceiptID` untuk menghapus baris lama.
  7. Panggil `repo.Update` untuk menyimpan perubahan header GR beserta *array* items yang baru.
  8. Lakukan `uow.Commit(txCtx)`. Pastikan ada `defer uow.Rollback()` jika terjadi error di tengah proses.

### 4. Handler & Routes (Menyediakan API Endpoint)
- Buka `internal/delivery/http/handler/goods_receipt_handler.go`.
- Tambahkan method `Update(c *fiber.Ctx) error` yang bertugas untuk parsing parameter ID, parsing *body* *request* menjadi DTO, dan memanggil fungsi `Update` dari UseCase.
- Buka `internal/delivery/http/route/route.go`.
- Daftarkan route baru berupa `r.Put("/:id", h.Update)` pada blok rute untuk `goods-receipts`.
