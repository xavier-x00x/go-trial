# Konsep Dynamic Safety Stock & Reorder Point (ROP) dalam ERP

Dokumen ini berisi rangkuman konsep arsitektur Supply Chain Management (SCM) untuk sistem ERP Retail Toko Kelontong Pro, khusus untuk algoritma Pemesanan Otomatis (Auto-Replenishment).

## 1. Anatomi Reorder Point (Titik Pemesanan Kembali)

Reorder Point (ROP) adalah titik jumlah stok di mana sistem harus secara otomatis membuat Draft Purchase Order (PO) kepada supplier agar barang tidak kehabisan.

**Rumus Dasar ROP:**
`ROP = Lead Time Demand + Safety Stock`

- **Lead Time Demand:** Jumlah rata-rata barang yang terjual selama kita menunggu truk dari supplier datang.
  *(Rumus: Rata-rata Penjualan Harian × LeadTimeDays)*
- **Safety Stock:** Stok pengaman / bantalan untuk mengantisipasi jika penjualan mendadak meledak, atau truk supplier terlambat datang.

---

## 2. Static Safety Stock vs Dynamic Safety Stock

### A. Static Safety Stock (Manual)
Pendekatan tradisional di mana Manajer Toko memasukkan angka mutlak ke dalam sistem.
- Di database kita, ini disimpan di tabel `Product` pada kolom `MinStockAlert`.
- **Kelemahan:** Jika tren penjualan berubah (misal menjelang lebaran), angka ini tidak relevan lagi. Jika manajer lupa update, toko kehabisan barang. Jika angka terlalu besar, uang toko nyangkut di gudang (Dead Stock).

### B. Dynamic Safety Stock (Otomatis)
Pendekatan A.I. (Machine Learning) dasar di mana sistem menghitung ulang angka bantalan setiap malam secara mandiri menggunakan ilmu statistika. Sistem menyesuaikan bantalan berdasarkan dua fluktuasi utama:

1. **Fluktuasi Penjualan (Demand Variability)**
   Sistem membaca histori penjualan 30 hari terakhir. Jika tren grafik menunjukkan kenaikan (musim ramai), sistem **menaikkan** Safety Stock secara otomatis. Jika musim sepi, sistem **menurunkan** Safety Stock untuk mengamankan *Cash Flow* perusahaan.

2. **Fluktuasi Kinerja Supplier (Lead Time Variability)**
   Sistem membaca histori kedatangan barang (PO vs Goods Receipt). Jika Supplier A bulan ini sering telat kirim akibat musim hujan (dari biasanya 3 hari menjadi 7 hari), sistem secara otomatis menebalkan Safety Stock khusus untuk barang dari Supplier A.

**Rumus Matematis (Statistika):**
`Dynamic Safety Stock = Z-Score × Akar Kuadrat [ (Rata² Lead Time × Variansi Penjualan) + (Rata² Penjualan × Variansi Lead Time) ]`
*(Z-Score adalah Konstanta Service Level, misal Z = 1.65 untuk kepastian 95% barang selalu ready).*

---

## 3. Korelasi dengan Struktur Database Kita

Karena Dynamic Safety Stock bersifat *On-The-Fly* (selalu berubah setiap hari berdasarkan transaksi), nilainya **TIDAK** direkam secara statis di dalam Master Data. Sebaliknya, Master Data kita bertugas menyediakan metrik konstan untuk diolah oleh algoritma ini:

1. **`entities/product_supplier.go`** menyediakan `LeadTimeDays` (sebagai Rata-rata Lead Time Dasar).
2. **`entities/product.go`** menyediakan `MinStockAlert` (sebagai Manual Fallback jika mode dinamis dimatikan).
3. **Tabel Transaksi Penjualan (SalesInvoiceItem)** akan menyediakan `Rata-rata Penjualan` & `Variansi Penjualan`.
4. **Tabel Transaksi Pembelian (PurchaseReceipt)** akan menyediakan `Variansi Lead Time` (Seberapa sering supplier ini ingkar janji waktu pengiriman).

Dengan memisahkan parameter dinamis dan statis seperti di atas, arsitektur database Master Data ini sudah siap 100% untuk ditancapkan algoritma pintar oleh tim Golang Backend!
