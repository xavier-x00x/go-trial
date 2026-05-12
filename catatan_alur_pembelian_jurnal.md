# Alur Siklus Pembelian (Procure to Pay) & Penjurnalan

Dokumen ini merangkum arsitektur alur dokumen transaksi dari mulai perencanaan pesanan (PO) hingga pencatatan jurnal akuntansi di Toko Kelontong Pro. Sistem ini mengadopsi standar ERP **3-Way Matching** (Pencocokan 3 Dokumen: PO, GR, Invoice) untuk mencegah *fraud* dan memastikan data keuangan valid.

## 1. Peta Dokumen & Efek ke Sistem

| Tahap | Entitas | Aksi Utama | Efek Stok | Efek Hutang | Jurnal Akuntansi |
|---|---|---|---|---|---|
| 1. Plan | `PurchaseOrderPlanning` | Rekomendasi oleh Cron | ❌ | ❌ | ❌ |
| 2. Order | `PurchaseOrder` | Dokumen resmi PO | ❌ | ❌ | ❌ |
| 3. Terima | `GoodsReceipt` | Cek fisik (CONFIRMED) | ✅ **Bertambah** (HPP update) | ❌ | ❌ |
| 4. Tagih | `PurchaseInvoice` | Terima Faktur (POSTED) | ❌ | ✅ **Bertambah** | ✅ Persediaan vs Hutang |
| 5. Bayar | `PurchasePayment` | Lunas/Cicil (POSTED) | ❌ | ✅ **Berkurang** | ✅ Hutang vs Kas/Bank |
| 6. Retur | `PurchaseReturn` | Kembalikan (CONFIRMED) | ✅ **Berkurang** | ✅ **Berkurang** | ✅ Hutang vs Persediaan |

*(Catatan Khusus)*:
| Pembelian Non-Stok | `ExpenseVoucher` | Bayar ATK/Jasa (POSTED) | ❌ | ❌ (CASH) / ✅ (CREDIT) | ✅ Beban vs Kas/Hutang |

---

## 2. Alur Kerja Detail (Step-by-Step)

### Step 1: Perencanaan Otomatis (Auto-Replenishment)
- **Kapan:** Setiap malam.
- **Proses:** Sistem menjumlahkan total stok semua gudang dalam satu toko. Jika total stok $\le$ Reorder Point, sistem membuat `PurchaseOrderPlanning` dengan status `PENDING`.
- **Aktor:** Sistem (Cron Job).

### Step 2: Pemesanan (Purchase Order)
- **Kapan:** Staf Purchasing meng-approve rekomendasi atau input PO manual.
- **Proses:** Dokumen `PurchaseOrder` (PO) diterbitkan ke Supplier. Status awal `DRAFT` $\rightarrow$ `SUBMITTED` $\rightarrow$ `APPROVED`. PO mengikat janji tapi belum memotong hutang atau menambah stok.
- **Aktor:** Staf Purchasing / Manajer.

### Step 3: Penerimaan Barang (Goods Receipt / BTB)
- **Kapan:** Truk supplier tiba bawa fisik barang.
- **Proses:** Staf gudang mengecek fisik. Jika oke, `GoodsReceipt` di-`CONFIRMED`.
- **Efek:** 
  - `InventoryStock.Quantity` bertambah.
  - HPP (`AverageBuyPrice`) dihitung ulang.
  - `PurchaseOrderItem.QtyReceived` diupdate.
- **Aktor:** Kepala Gudang.
- **Jurnal:** Belum dijurnal (menunggu faktur).

### Step 4: Faktur Pembelian (Purchase Invoice)
- **Kapan:** Faktur/tagihan kertas diterima oleh tim Keuangan.
- **Proses (3-Way Match):** Finance mencocokkan baris Faktur $\leftrightarrow$ baris PO $\leftrightarrow$ baris GR. Jika cocok, `PurchaseInvoice` di-`POSTED`.
- **Aktor:** Staf Finance / Akuntansi.
- **Efek Hutang:** Hutang Usaha (Account Payable) resmi diakui di sistem. Jatuh tempo dihitung mulai dari sini.
- **Jurnal Otomatis:**
  - **Debit:** Persediaan Barang Dagangan (Inventory)
  - **Kredit:** Hutang Usaha (Account Payable)

### Step 5: Pembayaran (Purchase Payment)
- **Kapan:** Tiba masa jatuh tempo faktur.
- **Proses:** Finance membuat `PurchasePayment`. Satu pembayaran bisa dialokasikan untuk memotong banyak Faktur sekaligus (Batch Payment). Status $\rightarrow$ `POSTED`.
- **Aktor:** Staf Finance.
- **Jurnal Otomatis:**
  - **Debit:** Hutang Usaha (Account Payable)
  - **Kredit:** Kas / Bank

### Step 6: Retur Pembelian (Purchase Return)
- **Kapan:** Barang diterima rusak / expired dan dikembalikan.
- **Proses:** Gudang menerbitkan `PurchaseReturn`. Status $\rightarrow$ `CONFIRMED`.
- **Aktor:** Kepala Gudang.
- **Efek:** Stok berkurang. Nota debit otomatis mengurangi total hutang kita ke supplier.
- **Jurnal Otomatis:**
  - **Debit:** Hutang Usaha (Account Payable)
  - **Kredit:** Persediaan Barang Dagangan (Inventory)

---

## 3. Pengecualian: Mode Langsung (Barang & Faktur Bersamaan)
Untuk toko retail, seringkali barang fisik datang bersamaan dengan faktur kertas dari sopir truk.
- **Solusi UX:** Terdapat endpoint khusus `CreateGoodsReceiptWithInvoiceRequest`.
- **Backend:** Dalam 1x klik, Backend secara otomatis membuat `GoodsReceipt` dan langsung mem-posting `PurchaseInvoice`. Database tetap bersih karena menyimpa 2 entitas, tapi pengguna merasa hanya menginput 1 dokumen.

---

## 4. Pembelian Keperluan Toko (Expense Voucher)
Jika membeli barang yang BUKAN untuk dijual (ATK, Sabun Pel, Kertas Struk, Jasa Service AC), **JANGAN pakai PO**. Gunakan `ExpenseVoucher`.
- Bebas menentukan `VendorName` tanpa harus daftarkan ke master Supplier.
- Langsung masuk ke Akun Beban (Expense).
- Mendukung Pembayaran `CASH` (langsung potong Kas) atau `CREDIT` (catat sebagai Hutang Lain-lain).
- **Jurnal Otomatis (Jika CASH):**
  - **Debit:** Beban ATK / Beban Kebersihan / dll.
  - **Kredit:** Kas / Bank

---

## 5. Mesin Jurnal (`JournalEntry` & `JournalEntryLine`)
Setiap dokumen yang berstatus **POSTED** (atau CONFIRMED pada kasus retur) akan diubah menjadi `JournalEntry` melalui sistem *Double-Entry*.
- Setiap dokumen punya rujukan balik (`SourceDocumentID` & `SourceDocumentType`).
- Jika dokumen dibatalkan (`VOID`/`CANCEL`), sistem tidak akan menghapus jurnal (Audit Trail), melainkan membuat **Jurnal Balik (Reversal)** yang menunjuk `ReversalOfID` ke jurnal aslinya.
