# Dokumen Persyaratan Produk (PRD)
# Sistem ERP Toko Kelontong Pro

> **Enterprise Resource Planning untuk Operasional Toko Kelontong Modern**

---

| Atribut | Keterangan |
|---|---|
| **Versi Dokumen** | v1.2 – Penambahan Master UoM & Jurnal Manual |
| **Tanggal** | April 2025 (Diperbarui) |
| **Disusun Oleh** | Tim Product Management – Divisi Fintech Ritel |
| **Status** | DRAFT – Menunggu Persetujuan Engineering Lead |
| **Ditujukan Kepada** | Engineering, UI/UX Designer, QA, Business Analyst |

---

## Daftar Isi

1. [Product Overview & Value Proposition](#1-product-overview--value-proposition)
2. [User Personas & Use Cases](#2-user-personas--use-cases)
3. [Detailed Functional Requirements](#3-detailed-functional-requirements)
4. [Data Model, Logika UoM & Logika Pajak](#4-data-model-logika-uom--logika-pajak)
   - [4.0 Master Measurement Unit (UoM)](#40-master-measurement-unit-uom)
5. [Journaling Logic – Tabel Penjurnalan Otomatis](#5-journaling-logic--tabel-penjurnalan-otomatis)
   - [D. Jurnal Manual (Non-Otomatis)](#d-jurnal-manual-non-otomatis)
6. [User Access & Security](#6-user-access--security)
7. [Acceptance Criteria & Success Metrics](#7-acceptance-criteria--success-metrics)
8. [Laporan Histori Barang Keluar Masuk per Bulan](#8-laporan-histori-barang-keluar-masuk-per-bulan)

---

## 1. Product Overview & Value Proposition

### 1.1 Latar Belakang & Masalah yang Diselesaikan

Toko kelontong di Indonesia merupakan tulang punggung ekonomi ritel skala mikro dengan estimasi lebih dari **3,5 juta unit usaha**. Namun, mayoritas masih dioperasikan secara manual — pencatatan stok di buku tulis, perhitungan laba rugi berdasarkan intuisi, dan tidak ada visibilitas terhadap arus kas. Permasalahan ini mengakibatkan:

- **Kebocoran stok** yang tidak terdeteksi (susut, kecurian, kesalahan input).
- **Ketidakmampuan memantau profitabilitas** per produk secara real-time.
- **Hutang ke supplier tidak terlacak**, berisiko menimbulkan masalah likuiditas.
- **Laporan keuangan tidak ada**, sehingga pemilik sulit mengajukan kredit ke bank.
- **Ketergantungan penuh pada satu orang** (pemilik) tanpa sistem delegasi yang aman.

### 1.2 Solusi: ERP Toko Kelontong Pro

**Sistem ERP Toko Kelontong Pro** adalah platform manajemen operasional terintegrasi berbasis web dan mobile yang dirancang khusus untuk kebutuhan toko kelontong skala kecil hingga menengah. Sistem ini menghubungkan seluruh rantai operasional — dari pemesanan barang ke supplier hingga laporan keuangan — dalam satu ekosistem digital yang akurat secara akuntansi.

> ### 🎯 Core Value Proposition
>
> | # | Proposisi | Deskripsi |
> |---|---|---|
> | 1 | **Akurasi Stok 99%+** | Setiap barang masuk dan keluar tercatat otomatis. Tidak ada celah kebocoran. |
> | 2 | **Laporan Keuangan Instan** | Laba Rugi, Neraca, dan Arus Kas tersedia real-time tanpa perlu akuntan manual. |
> | 3 | **Kontrol Penuh, Delegasi Aman** | Role-based access memastikan kasir hanya bisa kasir, gudang hanya bisa gudang. |
> | 4 | **Siap Pajak** | PPN masukan & keluaran dihitung otomatis, siap untuk pelaporan SPT. |
> | 5 | **Audit Trail Tak Terhapus** | Setiap perubahan data tercatat. Tidak ada yang bisa bermain curang tanpa jejak. |

### 1.3 Target Segmen Pengguna

- Toko kelontong dengan omzet **Rp 5 juta – Rp 500 juta per bulan**.
- Minimarket independen (non-franchise) dengan **1–5 karyawan**.
- Toko sembako dengan minimal **1 gudang dan 1 kasir** terpisah.
- Pemilik toko yang mulai tertarik pada **digitalisasi dan pelaporan keuangan formal**.

---

## 2. User Personas & Use Cases

### 2.1 Persona 1 – Owner / Pemilik Toko

> **👤 Pak Budi, 45 Tahun — Pemilik Toko Sembako "Maju Jaya"**
>
> - **Latar Belakang:** Pemilik toko sembako dengan 3 karyawan. Tidak berlatar belakang akuntan.
> - **Goals:** Ingin tahu toko laba atau rugi setiap hari tanpa harus hitung manual. Ingin kontrol karyawan dari jarak jauh.
> - **Pain Points:** Sering stok habis tiba-tiba, tidak tahu produk mana yang paling untung, takut karyawan korupsi.
> - **Perangkat:** Smartphone Android (utama), kadang laptop.
> - **Skenario Harian:** Membuka dashboard saat sarapan untuk cek omzet kemarin, laba kotor, dan alert stok hampir habis.

**Use Cases:**

- `UC-01` Melihat dashboard ringkasan harian (omzet, laba kotor, transaksi count).
- `UC-02` Menerima notifikasi stok hampir habis dan membuat PO ke supplier.
- `UC-03` Melihat laporan Laba Rugi bulan ini dibanding bulan lalu.
- `UC-04` Mengatur harga jual produk dan melihat margin keuntungan.

---

### 2.2 Persona 2 – Manajer Toko

> **👤 Mbak Sari, 32 Tahun — Manajer Operasional Harian**
>
> - **Latar Belakang:** Lulusan SMK Akuntansi. Mengelola toko 6 hari seminggu.
> - **Goals:** Memastikan stok selalu tersedia, proses pembelian rapi, dan laporan harian akurat.
> - **Pain Points:** Kesulitan rekap stok manual. Sering ada selisih antara catatan dan fisik.
> - **Skenario Harian:** Buka sistem → cek PO yang belum diterima → proses GR → approve penyesuaian stok.

**Use Cases:**

- `UC-05` Membuat dan mengirim Purchase Order ke supplier dari sistem.
- `UC-06` Menginput Good Receive dan mencocokkan dengan PO (3-way matching).
- `UC-07` Menyetujui penyesuaian stok hasil stock opname.
- `UC-08` Melihat histori pembelian per supplier dan membandingkan harga.

---

### 2.3 Persona 3 – Kasir

> **👤 Dini, 22 Tahun — Kasir Shift Pagi**
>
> - **Latar Belakang:** Lulusan SMA. Baru bekerja 3 bulan di toko.
> - **Goals:** Melayani pelanggan cepat tanpa antrian panjang. Tidak salah hitung kembalian.
> - **Pain Points:** Sering bingung jika ada retur barang. Takut salah input harga.
> - **Skenario Harian:** Login → buka sesi kasir → scan barcode → konfirmasi pembayaran → cetak struk.

**Use Cases:**

- `UC-09` Memproses transaksi penjualan dengan barcode scanner.
- `UC-10` Menghitung kembalian uang tunai secara otomatis.
- `UC-11` Memproses retur penjualan dengan referensi nomor transaksi.
- `UC-12` Menutup sesi kasir dan rekap total penjualan shift.

---

### 2.4 Persona 4 – Admin Gudang

> **👤 Topan, 28 Tahun — Admin Gudang & Penerimaan Barang**
>
> - **Latar Belakang:** Lulusan SMA. Bertanggung jawab penuh atas fisik gudang.
> - **Goals:** Pastikan barang yang datang sesuai pesanan. Stok tersusun rapi dan akurat.
> - **Pain Points:** Sering ada selisih antara yang dipesan dan yang datang. Bingung harus lapor ke siapa.
> - **Skenario Harian:** Terima barang → input GR di sistem → sistem otomatis flag jika ada selisih → lapor ke manajer.

**Use Cases:**

- `UC-13` Menerima notifikasi PO yang sedang dikirim supplier.
- `UC-14` Input Good Receive dengan scan barcode atau input manual.
- `UC-15` Melihat lokasi rak stok untuk setiap produk.
- `UC-16` Input data stok opname (hitungan fisik).

---

## 3. Detailed Functional Requirements

> **Keterangan Prioritas:**
> - `P0` **(Must Have / Kritis)** — Wajib ada di MVP. Tanpa fitur ini, sistem tidak bisa beroperasi.
> - `P1` **(Should Have / Penting)** — Sangat disarankan masuk Sprint 2. Meningkatkan efisiensi signifikan.
> - `P2` **(Nice to Have)** — Fitur tambahan untuk iterasi berikutnya setelah validasi pasar.

### 3.1 Modul Manajemen Produk

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 1 | **Kategori & Sub-Kategori** | Pengelompokan produk ke dalam hierarki 2 level (misal: Sembako > Beras). Mendukung pencarian & filter berbasis kategori pada POS dan laporan. | `P0` | BE + FE |
| 2 | **Multi UoM & Konversi Otomatis** | Produk dapat memiliki UoM pembelian (Karton/Dusun) dan UoM penjualan (Pcs/Bungkus). Sistem konversi otomatis saat PO, GR, dan transaksi kasir. | `P0` | BE + FE |
| 3 | **Manajemen Pajak (PPN)** | Konfigurasi PPN fleksibel per produk/kelompok produk. Mendukung pajak masukan (PM) dan pajak keluaran (PK) dengan persentase yang dapat diubah (default 11%). | `P0` | BE + Akuntan |
| 4 | **Manajemen Harga Jual & HPP** | Penetapan harga jual dengan margin target. HPP dihitung secara otomatis menggunakan metode FIFO/Average. Alert saat harga jual di bawah HPP. | `P0` | BE + FE |

### 3.2 Modul Pembelian (Purchase Cycle)

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 5 | **Purchase Order (PO)** | Pembuatan dokumen PO resmi ke supplier dengan status: Draft → Approved → Sent → Received. Mendukung partial fulfillment dan cancellation dengan alasan. | `P0` | BE + FE |
| 6 | **Good Receive (GR) & 3-Way Matching** | Penerimaan barang dicocokkan dengan PO (kuantitas & harga). Sistem otomatis mendeteksi selisih dan membutuhkan approval sebelum stok diperbarui. | `P0` | BE + Gudang |
| 7 | **Manajemen Supplier** | Database supplier lengkap: nama, NPWP, alamat, kontak, termin pembayaran, histori harga per produk. Perbandingan harga antar supplier otomatis. | `P1` | BE + FE |
| 8 | **Hutang Dagang (AP)** | Pencatatan hutang ke supplier dari GR. Jadwal jatuh tempo, pembayaran parsial/penuh, rekonsiliasi otomatis dengan jurnal pelunasan hutang. | `P0` | BE + Akuntan |

### 3.3 Modul Penjualan (POS)

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 9 | **Point of Sale (POS) – Kasir** | Interface kasir cepat: scan barcode, pencarian produk, perhitungan kembalian otomatis. Mendukung multi-metode bayar (tunai, QRIS, kartu debit). | `P0` | FE + BE |
| 10 | **Manajemen Retur Penjualan** | Proses retur barang dari pelanggan dengan referensi transaksi asal. Stok otomatis bertambah, jurnal reversal dibuat secara otomatis. | `P1` | BE + FE |

### 3.4 Modul Inventaris

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 11 | **Manajemen Stok Real-Time** | Stok diperbarui real-time pada setiap transaksi (pembelian, penjualan, retur, penyesuaian). Tampilan stok per lokasi/rak gudang. | `P0` | BE |
| 12 | **Stok Opname & Penyesuaian** | Fitur stok opname periodik: input hitungan fisik, sistem hitung selisih, approval manajer, lalu jurnal penyesuaian stok otomatis dibuat. | `P1` | BE + Gudang |
| 13 | **Reorder Point & Alert Stok** | Konfigurasi minimum stok (reorder point) per produk. Notifikasi otomatis (in-app dan email) ketika stok mencapai atau di bawah batas minimum. | `P1` | BE + FE |

### 3.5 Modul Keuangan & Akuntansi

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 14 | **Chart of Accounts (CoA)** | Daftar akun standar ritel Indonesia (5 kelompok: Aset, Liabilitas, Ekuitas, Pendapatan, Beban). Dapat dikustomisasi dan ditambah sub-akun oleh akuntan. | `P0` | BE + Akuntan |
| 15 | **Automated Double-Entry Journal** | Setiap transaksi menghasilkan jurnal Debit/Kredit otomatis yang tidak dapat diedit (immutable). Koreksi hanya melalui jurnal reversal. | `P0` | BE + Akuntan |
| 16 | **Laporan Keuangan Real-Time** | Laba Rugi, Neraca (Balance Sheet), dan Arus Kas yang dapat difilter per periode, kategori, dan dapat diekspor ke Excel/PDF. | `P0` | BE + FE |

### 3.6 Modul User & Keamanan

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 17 | **Role-Based Access Control (RBAC)** | 4 role bawaan: Owner, Manajer Toko, Kasir, Admin Gudang. Setiap role memiliki izin akses granular per modul (Create/Read/Update/Delete). | `P0` | BE + DevOps |
| 18 | **Audit Trail** | Setiap perubahan data (stok, harga, transaksi, akun) dicatat: siapa, kapan, nilai sebelum dan sesudah perubahan. Tidak dapat dihapus oleh siapapun. | `P0` | BE |

### 3.7 Modul Master Measurement Unit (UoM)

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 19 | **Master Satuan Ukur (UoM Library)** | Halaman administrasi terpusat untuk mendefinisikan semua satuan ukur yang tersedia di sistem (Pcs, Kg, Gram, Liter, ml, Lusin, Karton, dst.). Setiap satuan memiliki tipe dimensi (Jumlah, Berat, Volume) dan dapat dinonaktifkan tanpa menghapus data historis. | `P0` | BE + FE |
| 20 | **Kelompok Konversi UoM (UoM Category)** | Satuan yang dapat saling dikonversi dikelompokkan dalam satu UoM Category (misal: Kategori *Berat* mencakup Gram, Ons, Kg, Kuintal, Ton). Konversi antar satuan dalam satu kategori bersifat sistem (tidak perlu input manual per produk). | `P0` | BE |
| 21 | **Konversi Kustom per Produk** | Untuk UoM lintas kategori (misal: Karton → Pcs yang tidak memiliki hubungan dimensi standar), konversi didefinisikan secara spesifik per produk. Nilai faktor konversi dapat diubah oleh admin; perubahan berlaku untuk transaksi baru (tidak retroaktif). | `P0` | BE + FE |
| 22 | **Validasi Konsistensi UoM** | Sistem memvalidasi bahwa UoM beli dan UoM jual suatu produk berada dalam kelompok konversi yang sama, atau memiliki konversi kustom yang terdefinisi. Jika tidak, sistem menolak penyimpanan produk dengan pesan error yang jelas. | `P0` | BE |

### 3.8 Modul Jurnal Manual

| No | Fitur | Deskripsi | Prioritas | PIC |
|---|---|---|:---:|---|
| 23 | **Entry Jurnal Manual** | Form input jurnal Debit/Kredit bebas untuk mencatat transaksi yang tidak dihasilkan otomatis oleh sistem: jasa perbaikan, biaya tak terduga, koreksi akuntansi, amortisasi, depresiasi aset, dll. Mendukung multi-baris (satu entry, banyak akun). | `P0` | BE + FE |
| 24 | **Kategori Jurnal Manual** | Setiap jurnal manual wajib dikategorikan berdasarkan template yang tersedia: *Biaya Jasa & Servis*, *Koreksi Akuntansi*, *Depresiasi Aset*, *Amortisasi*, *Transfer Antar Akun*, *Lain-Lain*. Memudahkan filter dan analisis laporan. | `P1` | BE + FE |
| 25 | **Alur Approval Jurnal Manual** | Jurnal manual yang dibuat oleh Manajer Toko berstatus `Pending Approval` dan harus disetujui oleh Owner sebelum di-posting ke buku besar. Jurnal yang dibuat langsung oleh Owner langsung berstatus `Posted`. | `P0` | BE |
| 26 | **Jurnal Reversal Manual** | Setiap jurnal manual yang sudah di-posting dapat di-reverse oleh Owner. Sistem membuat jurnal baru dengan nilai Debit/Kredit terbalik secara otomatis, dengan referensi ke jurnal asal. Jurnal asal tetap tersimpan (tidak dihapus). | `P0` | BE |
| 27 | **Template Jurnal Berulang** | Owner/Manajer dapat menyimpan jurnal manual sebagai template (misal: pembayaran sewa bulanan). Template dapat dijadwalkan (otomatis dibuat tiap bulan pada tanggal tertentu) atau dipanggil manual saat dibutuhkan. | `P1` | BE + FE |

---

## 4. Data Model, Logika UoM & Logika Pajak

### 4.0 Master Measurement Unit (UoM)

#### 4.0.1 Konsep & Hierarki UoM

Sistem menggunakan arsitektur UoM dua lapis: **UoM Category** (kelompok dimensi) yang menaungi satu atau lebih **UoM** (satuan konkret). Setiap produk wajib memiliki tepat satu `uom_base` (satuan penyimpanan stok terkecil) yang tidak berubah selama masa hidup produk.

```
UoM Category
    └── UoM (satuan konkret)
            └── Faktor Konversi ke Base Unit dalam Kategori yang sama

Produk
    ├── uom_base       → satuan penyimpanan stok (tidak berubah)
    ├── uom_beli       → satuan pada PO & GR
    ├── uom_jual       → satuan pada POS & nota
    └── konversi_kustom[] → opsional, jika uom_beli/jual lintas kategori
```

#### 4.0.2 Daftar Master UoM Bawaan Sistem

| Kode UoM | Nama Satuan | UoM Category | Faktor ke Base Unit Kategori | Tipe Dimensi | Keterangan |
|---|---|---|:---:|---|---|
| `PCS` | Pieces / Buah | Jumlah | **1** *(base)* | Jumlah | Satuan cacah tunggal. Base unit kategori Jumlah. |
| `LUSIN` | Lusin | Jumlah | 12 | Jumlah | 1 Lusin = 12 Pcs |
| `KODI` | Kodi | Jumlah | 20 | Jumlah | 1 Kodi = 20 Pcs |
| `GROSS` | Gros | Jumlah | 144 | Jumlah | 1 Gros = 144 Pcs |
| `PACK` | Pack / Bungkus | Jumlah | *kustom per produk* | Jumlah | Isi per pack ditentukan per produk. |
| `BOX` | Box / Kotak | Jumlah | *kustom per produk* | Jumlah | Isi per box ditentukan per produk. |
| `KARTON` | Karton / Dus | Jumlah | *kustom per produk* | Jumlah | Isi per karton ditentukan per produk. |
| `SAK` | Sak / Kantong | Jumlah | *kustom per produk* | Jumlah | Isi per sak ditentukan per produk. |
| `BOTOL` | Botol | Jumlah | *kustom per produk* | Jumlah | Ukuran botol ditentukan per produk. |
| `GRAM` | Gram | Berat | **1** *(base)* | Berat | Base unit kategori Berat. |
| `ONS` | Ons | Berat | 100 | Berat | 1 Ons = 100 Gram |
| `KG` | Kilogram | Berat | 1.000 | Berat | 1 Kg = 1.000 Gram |
| `KUINTAL` | Kuintal | Berat | 100.000 | Berat | 1 Kuintal = 100.000 Gram |
| `TON` | Ton | Berat | 1.000.000 | Berat | 1 Ton = 1.000.000 Gram |
| `ML` | Mililiter | Volume | **1** *(base)* | Volume | Base unit kategori Volume. |
| `CL` | Centiliter | Volume | 10 | Volume | 1 Cl = 10 ml |
| `DL` | Desiliter | Volume | 100 | Volume | 1 Dl = 100 ml |
| `LITER` | Liter | Volume | 1.000 | Volume | 1 Liter = 1.000 ml |
| `M3` | Meter Kubik | Volume | 1.000.000 | Volume | 1 m³ = 1.000.000 ml |

> **Catatan:** UoM dengan faktor *kustom per produk* tidak memiliki nilai konversi sistem yang tetap karena bergantung pada spesifikasi kemasan masing-masing produk. Nilai konversi ini wajib diisi saat pendaftaran produk baru.

#### 4.0.3 Skema Database: Tabel `uom` & `uom_category`

```sql
-- Tabel Kelompok Dimensi UoM
CREATE TABLE uom_category (
  id          INT PRIMARY KEY AUTO_INCREMENT,
  kode        VARCHAR(20)  NOT NULL UNIQUE,  -- 'JUMLAH', 'BERAT', 'VOLUME'
  nama        VARCHAR(50)  NOT NULL,
  keterangan  VARCHAR(255),
  is_active   BOOLEAN DEFAULT TRUE,
  created_at  DATETIME DEFAULT NOW()
);

-- Tabel Master Satuan Ukur
CREATE TABLE uom (
  id              INT PRIMARY KEY AUTO_INCREMENT,
  kode            VARCHAR(20)  NOT NULL UNIQUE,  -- 'PCS', 'KG', 'LITER', dst.
  nama            VARCHAR(50)  NOT NULL,
  category_id     INT NOT NULL REFERENCES uom_category(id),
  faktor_ke_base  DECIMAL(20,6) NOT NULL DEFAULT 1,
  -- faktor_ke_base = NULL berarti konversi kustom (ditentukan per produk)
  -- faktor_ke_base = 1 berarti UoM ini adalah base unit kategorinya
  is_base_unit    BOOLEAN DEFAULT FALSE,
  is_custom       BOOLEAN DEFAULT FALSE, -- TRUE jika faktor tergantung per produk
  is_active       BOOLEAN DEFAULT TRUE,
  is_system       BOOLEAN DEFAULT FALSE, -- TRUE = bawaan sistem, tidak bisa dihapus
  created_by      INT REFERENCES users(id),
  created_at      DATETIME DEFAULT NOW(),

  INDEX idx_category (category_id),
  INDEX idx_kode (kode)
);

-- Tabel Konversi Kustom per Produk (untuk UoM dengan is_custom = TRUE)
CREATE TABLE product_uom_conversion (
  id                INT PRIMARY KEY AUTO_INCREMENT,
  product_id        INT NOT NULL REFERENCES products(id),
  uom_id            INT NOT NULL REFERENCES uom(id),     -- UoM yang dikonversi (misal: KARTON)
  uom_base_id       INT NOT NULL REFERENCES uom(id),     -- UoM base produk (misal: PCS)
  faktor_konversi   DECIMAL(20,6) NOT NULL,               -- misal: 40 (1 Karton = 40 Pcs)
  berlaku_mulai     DATE NOT NULL,
  berlaku_sampai    DATE,                                  -- NULL = berlaku selamanya
  created_by        INT REFERENCES users(id),
  created_at        DATETIME DEFAULT NOW(),

  UNIQUE KEY uq_product_uom_date (product_id, uom_id, berlaku_mulai),
  INDEX idx_product (product_id)
);
```

> **Aturan Perubahan Faktor Konversi:** Jika supplier mengubah isi karton (misal dari 40 Pcs menjadi 48 Pcs mulai tanggal tertentu), admin membuat record baru di `product_uom_conversion` dengan `berlaku_mulai` yang baru. Record lama tetap tersimpan untuk menjaga integritas historis PO dan GR sebelumnya.

#### 4.0.4 Aturan Bisnis Validasi UoM

| # | Aturan | Konsekuensi Pelanggaran |
|---|---|---|
| V1 | UoM beli dan UoM jual harus berada dalam UoM Category yang sama, ATAU memiliki konversi kustom yang terdefinisi di `product_uom_conversion`. | Sistem menolak simpan produk dengan error: *"UoM beli dan jual tidak kompatibel. Tambahkan konversi kustom terlebih dahulu."* |
| V2 | `uom_base` produk tidak boleh diubah setelah ada transaksi. | Tombol edit `uom_base` dikunci (disabled) setelah produk memiliki minimal 1 transaksi. Perubahan hanya melalui tiket ke System Admin. |
| V3 | UoM sistem (`is_system = TRUE`) tidak dapat dihapus atau dinonaktifkan selama ada produk aktif yang menggunakannya. | Sistem menampilkan daftar produk yang terpengaruh dan meminta admin untuk memigrasikannya terlebih dahulu. |
| V4 | Faktor konversi wajib bernilai > 0. | Validasi frontend dan backend menolak nilai ≤ 0. |
| V5 | Perubahan faktor konversi tidak berlaku retroaktif. | Transaksi historis tetap menggunakan faktor lama yang di-snapshot saat transaksi dibuat. |

---

### 4.1 Logika Konversi Multi-Unit of Measurement (UoM)

Sistem menggunakan konsep **Base Unit** (satuan terkecil penjualan) sebagai anchor konversi. Semua stok disimpan dalam Base Unit. Faktor konversi didefinisikan per produk oleh admin.

#### 4.1.1 Formula Konversi

```
Stok Bertambah (GR)  :  Qty_GR × Faktor_Konversi_Beli  =  Penambahan Stok (Base Unit)
Stok Berkurang (POS) :  Qty_Jual × Faktor_Konversi_Jual  =  Pengurangan Stok (Base Unit)
HPP per Base Unit    :  Harga_Beli / Faktor_Konversi_Beli
Harga Jual per Unit  :  HPP_BaseUnit × (1 + Margin%) × (1 + PPN%)
```

**Contoh Kalkulasi:**
```
Beli Indomie 1 Karton (40 Pcs) @ Rp 120.000
  → HPP per Pcs         = Rp 120.000 / 40                   = Rp 3.000
  → Harga Jual          = Rp 3.000 × 1.20 (margin) × 1.11 (PPN)
                        = Rp 3.996 ≈ dibulatkan Rp 4.000
```

#### 4.1.2 Tabel Contoh Konfigurasi Multi-UoM

| Produk | UoM Beli | Faktor Konversi | UoM Jual | Stok Disimpan Dalam | Contoh Kalkulasi |
|---|---|:---:|---|---|---|
| Indomie Goreng | Karton | 1 Karton = **40 Pcs** | Pcs | Pcs (base unit) | Beli 2 Karton → Stok +80 Pcs |
| Gula Pasir | Sak | 1 Sak = **25 Kg** | Kg / Ons | Gram (base unit) | Beli 3 Sak → Stok +75.000 gram |
| Minyak Goreng | Dus (12 Btl) | 1 Dus = **12 Botol** | Botol / Liter | Botol (base unit) | Beli 5 Dus → Stok +60 Botol |
| Beras Premium | Ton | 1 Ton = **1.000.000 gram** | Kg | Gram (base unit) | Beli 0,5 Ton → Stok +500 Kg |

---

### 4.2 Logika Perhitungan Pajak PPN

#### 4.2.1 Pajak Masukan (PM) – Saat Pembelian

Pajak masukan adalah PPN yang dibayar toko saat membeli barang dari supplier yang merupakan **Pengusaha Kena Pajak (PKP)**. PM ini dicatat sebagai aset sementara dan dapat dikreditkan (dikurangkan) dari pajak keluaran.

```
Dasar Pengenaan Pajak (DPP)  =  Total Nilai Beli (ekskl. PPN)
PPN Masukan                  =  DPP × Tarif PPN (default 11%)
Total Hutang ke Supplier     =  DPP + PPN Masukan

Contoh: Beli barang DPP Rp 10.000.000
  → PPN Masukan  = Rp 10.000.000 × 11%          = Rp 1.100.000
  → Total AP     = Rp 10.000.000 + Rp 1.100.000 = Rp 11.100.000

  Jurnal:
    Debit  Persediaan Barang Dagang   Rp 10.000.000
    Debit  PPN Masukan                Rp  1.100.000
    Kredit Hutang Dagang (AP)                       Rp 11.100.000
```

#### 4.2.2 Pajak Keluaran (PK) – Saat Penjualan

Pajak keluaran adalah PPN yang **dipungut toko dari pelanggan** atas penjualan barang. PK ini merupakan kewajiban toko kepada negara dan akan dibayarkan setelah dikurangi pajak masukan.

```
Harga Jual Ekskl. PPN         =  Harga yang ditetapkan admin
PPN Keluaran                   =  Harga Jual Ekskl. PPN × Tarif PPN (default 11%)
Harga Jual Inkl. PPN (bayar)  =  Harga Jual Ekskl. PPN + PPN Keluaran

Contoh: Jual Indomie 1 Pcs @ Rp 3.600 (ekskl. PPN)
  → PPN Keluaran = Rp 3.600 × 11%              = Rp 396
  → Harga Bayar  = Rp 3.600 + Rp 396           = Rp 3.996 (≈ Rp 4.000)

  Jurnal:
    Debit  Kas                                    Rp  4.000
    Kredit Pendapatan Penjualan                               Rp 3.600
    Kredit PPN Keluaran                                       Rp   400  ← (setelah pembulatan)
```

---

### 4.3 Entitas Data Utama & Relasi

| Entitas | Field Utama |
|---|---|
| **UomCategory** | `id`, `kode`, `nama`, `keterangan`, `is_active` |
| **Uom** | `id`, `kode`, `nama`, `category_id`, `faktor_ke_base`, `is_base_unit`, `is_custom`, `is_active`, `is_system` |
| **ProductUomConversion** | `id`, `product_id`, `uom_id`, `uom_base_id`, `faktor_konversi`, `berlaku_mulai`, `berlaku_sampai` |
| **Product** | `id`, `nama`, `kode_sku`, `kategori_id`, `uom_base_id`, `uom_beli_id`, `uom_jual_id`, `harga_jual`, `harga_beli_terakhir`, `reorder_point`, `pajak_id`, `is_active` |
| **PurchaseOrder** | `id`, `nomor_po`, `supplier_id`, `tanggal_po`, `status` *(draft/approved/sent/received)*, `total_nilai`, `created_by`, `approved_by` |
| **POItem** | `id`, `po_id`, `product_id`, `qty_pesan` *(dalam UoM beli)*, `uom_id`, `faktor_konversi_snapshot`, `harga_satuan`, `subtotal` |
| **GoodReceive** | `id`, `nomor_gr`, `po_id`, `tanggal_gr`, `received_by`, `status` *(draft/confirmed)* |
| **GRItem** | `id`, `gr_id`, `product_id`, `qty_pesan`, `qty_diterima`, `qty_selisih`, `uom_id`, `faktor_konversi_snapshot`, `harga_satuan` |
| **SalesTransaction** | `id`, `nomor_transaksi`, `tanggal`, `kasir_id`, `metode_bayar`, `total_bayar`, `total_kembalian`, `status` |
| **SalesItem** | `id`, `transaction_id`, `product_id`, `qty`, `uom_id`, `harga_jual`, `diskon`, `ppn_keluaran`, `hpp_saat_jual` |
| **JournalEntry** | `id`, `nomor_jurnal`, `tanggal`, `tipe` *(AUTO/MANUAL)*, `kategori_manual`, `referensi`, `deskripsi`, `status` *(draft/pending_approval/posted/reversed)*, `created_by`, `approved_by`, `is_reversed`, `reversal_of_id` |
| **JournalLine** | `id`, `journal_id`, `account_id`, `debit`, `kredit`, `keterangan` |
| **JournalTemplate** | `id`, `nama_template`, `kategori`, `jadwal_hari`, `is_recurring`, `created_by`, `is_active` |
| **JournalTemplateLine** | `id`, `template_id`, `account_id`, `debit_default`, `kredit_default`, `keterangan_default` |
| **Account (CoA)** | `id`, `kode_akun`, `nama_akun`, `tipe` *(Aset/Liabilitas/Ekuitas/Pendapatan/Beban)*, `is_active` |

---

## 5. Journaling Logic – Tabel Penjurnalan Otomatis

Sistem mengimplementasikan **double-entry bookkeeping** secara penuh. Setiap transaksi menghasilkan satu `Journal Entry` dengan dua atau lebih `Journal Lines` yang selalu memenuhi kondisi:

```
∑ Debit = ∑ Kredit
```

Jurnal bersifat **immutable** — tidak dapat diubah atau dihapus. Koreksi dilakukan melalui **jurnal reversal** yang menghasilkan jurnal baru bernilai negatif.

---

### A. Siklus Pembelian (Purchase Cycle)

| No | Event / Transaksi | Akun Debit | Akun Kredit | Keterangan | Modul Pemicu |
|:---:|---|---|---|---|---|
| 1 | **Penerimaan Barang / Good Receive (GR)** | Persediaan Barang Dagang `(1-1300)` | Hutang Dagang / AP `(2-1100)` | Nilai = Qty × Harga PO (ekskl. PPN). Stok langsung bertambah. | GR Module |
| 2 | **PPN Masukan pada GR** | PPN Masukan `(1-1400)` | Hutang Dagang / AP `(2-1100)` | Nilai = Harga Beli × Tarif PPN (default 11%). Diklaim sebagai kredit pajak. | GR Module |
| 3 | **Pembayaran Hutang ke Supplier** | Hutang Dagang / AP `(2-1100)` | Kas / Bank `(1-1100 / 1-1200)` | Pembayaran penuh atau sebagian. Selisih bayar dicatat sebagai diskon atau denda. | AP Module |
| 4 | **Retur Pembelian ke Supplier** | Hutang Dagang / AP `(2-1100)` | Persediaan Barang Dagang `(1-1300)` | Jurnal kebalikan dari GR. Stok berkurang. PPN Masukan di-reverse. | Purchase Return |

---

### B. Siklus Penjualan (Sales Cycle)

| No | Event / Transaksi | Akun Debit | Akun Kredit | Keterangan | Modul Pemicu |
|:---:|---|---|---|---|---|
| 5 | **Penjualan Tunai – Kasir POS** | Kas `(1-1100)` | Pendapatan Penjualan `(4-1000)` + PPN Keluaran `(2-1300)` | Kas debit = total termasuk PPN. Pendapatan diakui ekskl. PPN. | POS Module |
| 6 | **HPP saat Penjualan (COGS Recognition)** | Harga Pokok Penjualan / HPP `(5-1000)` | Persediaan Barang Dagang `(1-1300)` | Otomatis bersamaan jurnal penjualan. Nilai HPP dari FIFO/Average. | POS Module |
| 7 | **Penjualan Non-Tunai (QRIS / Debit)** | Piutang Pembayaran Digital `(1-1250)` | Pendapatan Penjualan `(4-1000)` + PPN Keluaran `(2-1300)` | Diakui saat settlement dana masuk ke rekening bank toko. | POS + Payment |
| 8 | **Retur Penjualan dari Pelanggan** | Retur & Potongan Penjualan `(4-1100)` | Kas `(1-1100)` + PPN Keluaran `(2-1300)` *(reverse)* | Jurnal debit kontra-pendapatan. Stok bertambah. HPP di-reverse. | Sales Return |

---

### C. Biaya Operasional & Penyesuaian Stok

| No | Event / Transaksi | Akun Debit | Akun Kredit | Keterangan | Modul Pemicu |
|:---:|---|---|---|---|---|
| 9 | **Pembayaran Biaya Operasional** *(Listrik, Sewa, dll.)* | Beban Operasional `(5-2000 s.d. 5-2900)` | Kas / Bank `(1-1100 / 1-1200)` | Dikategorikan per jenis beban sesuai CoA. Dilakukan oleh Owner/Manajer. | Expense Module |
| 10 | **Penyesuaian Stok – Selisih Minus** *(Susut/Kecurian)* | Beban Penyusutan Stok `(5-3000)` | Persediaan Barang Dagang `(1-1300)` | Dibuat setelah stock opname diapprove Manajer. Nilai = HPP × Qty selisih. | Inventory Adj. |
| 11 | **Penyesuaian Stok – Selisih Plus** *(Surplus)* | Persediaan Barang Dagang `(1-1300)` | Pendapatan Lain-Lain `(4-2000)` | Stok lebih dari catatan diakui sebagai pendapatan lain-lain. | Inventory Adj. |
| 12 | **Setoran Modal Awal / Tambahan Modal** | Kas / Bank `(1-1100 / 1-1200)` | Modal Pemilik `(3-1000)` | Dilakukan saat setup awal atau penambahan modal oleh owner. | Capital Module |

---

### D. Jurnal Manual (Non-Otomatis)

> Jurnal manual dibuat secara bebas oleh Owner atau Manajer Toko untuk mencatat transaksi yang **tidak dihasilkan oleh alur sistem otomatis**. Setiap jurnal manual wajib memenuhi `∑ Debit = ∑ Kredit` sebelum dapat di-submit. Jurnal oleh Manajer membutuhkan approval Owner sebelum posting ke buku besar.

| No | Kategori | Contoh Transaksi | Akun Debit (Umum) | Akun Kredit (Umum) | Catatan Penting |
|:---:|---|---|---|---|---|
| 13 | **Biaya Jasa & Servis** | Bayar jasa perbaikan kulkas toko | Beban Pemeliharaan & Perbaikan `(5-2300)` | Kas `(1-1100)` | Wajib lampirkan nomor nota/invoice servis sebagai bukti. |
| 14 | **Biaya Jasa & Servis – Kredit** | Jasa perbaikan belum dibayar (hutang jasa) | Beban Pemeliharaan & Perbaikan `(5-2300)` | Hutang Lain-Lain `(2-1200)` | Saat pembayaran: Debit Hutang Lain-Lain, Kredit Kas. |
| 15 | **Depresiasi Aset Tetap** | Penyusutan bulanan mesin kasir, rak, dll. | Beban Depresiasi Aset `(5-4000)` | Akumulasi Depresiasi `(1-2100)` | Dijadwalkan otomatis tiap bulan via template jurnal berulang. |
| 16 | **Amortisasi Biaya Dibayar Dimuka** | Amortisasi sewa toko yang dibayar 1 tahun di muka | Beban Sewa Tempat `(5-2100)` | Biaya Dibayar Dimuka `(1-1500)` | Dibuat tiap bulan sejumlah 1/12 dari total sewa tahunan. |
| 17 | **Koreksi Akuntansi** | Memperbaiki akun yang salah catat di jurnal sebelumnya | *(Akun yang benar – Debit)* | *(Akun yang benar – Kredit)* | Wajib sertakan referensi nomor jurnal yang dikoreksi. Jurnal asal di-reverse terlebih dahulu. |
| 18 | **Transfer Antar Akun Kas/Bank** | Setor uang kas toko ke rekening bank | Bank – Rekening Toko `(1-1200)` | Kas `(1-1100)` | Nilai transfer tidak termasuk PPN. Dilakukan oleh Owner saja. |
| 19 | **Prive / Pengambilan Modal Pribadi** | Owner menarik uang toko untuk kebutuhan pribadi | Prive / Penarikan Modal `(3-2000)` | Kas / Bank `(1-1100 / 1-1200)` | Menurunkan ekuitas. Wajib ada persetujuan Owner sendiri. |
| 20 | **Pendapatan Lain Non-Operasional** | Menerima uang ganti rugi dari vendor | Kas `(1-1100)` | Pendapatan Lain-Lain `(4-2000)` | Untuk penerimaan di luar penjualan barang dagangan. |
| 21 | **Beban Non-Operasional** | Denda keterlambatan pembayaran ke supplier | Beban Non-Operasional `(5-5000)` | Kas / Bank `(1-1100 / 1-1200)` | Dipisahkan dari beban operasional untuk kejernihan laporan. |
| 22 | **Pencatatan Aset Baru (Pembelian Tunai)** | Beli rak gudang baru seharga Rp 5.000.000 | Aset Tetap – Peralatan `(1-2000)` | Kas `(1-1100)` | Tidak melalui modul pembelian barang dagang (PO). |
| 23 | **Penghapusan Aset (Disposal)** | Aset rusak total / dijual murah | Akumulasi Depresiasi `(1-2100)` + Rugi Penjualan Aset `(5-4100)` | Aset Tetap – Peralatan `(1-2000)` | Sisa nilai buku dibebankan sebagai kerugian. |

#### 5.2 Alur Status Jurnal Manual

```
[Owner / Manajer]
      │
      ▼
  ┌─────────┐     Submit      ┌──────────────────┐    Approve (Owner)   ┌────────┐
  │  DRAFT  │ ─────────────▶ │ PENDING APPROVAL  │ ─────────────────▶  │ POSTED │
  └─────────┘                └──────────────────┘                      └────────┘
      ▲                              │ Tolak (Owner)                        │
      └──────── Kembalikan ──────────┘                                      │ Reverse
                                                                            ▼
                                                                      ┌──────────┐
                                                                      │ REVERSED │
                                                                      └──────────┘

Catatan:
  • Jurnal yang dibuat LANGSUNG oleh Owner melewati tahap PENDING APPROVAL → langsung POSTED.
  • Hanya Owner yang dapat APPROVE, TOLAK, dan REVERSE jurnal manual.
  • Status POSTED dan REVERSED bersifat final — tidak ada penghapusan.
```

#### 5.3 Contoh Entry Jurnal Manual: Jasa Perbaikan Kulkas

```
Tanggal    : 10 April 2025
Nomor      : JM-2025-0042
Kategori   : Biaya Jasa & Servis
Deskripsi  : Pembayaran jasa servis kulkas display oleh CV Teknik Jaya
             Nota No. TJ-2025/0410
Dibuat oleh: Manajer Toko (Sari)
Status     : Posted (disetujui Owner jam 11:23)

┌────┬─────────────────────────────────────────┬──────────────┬──────────────┐
│ No │ Akun                                    │    Debit     │   Kredit     │
├────┼─────────────────────────────────────────┼──────────────┼──────────────┤
│  1 │ 5-2300  Beban Pemeliharaan & Perbaikan  │ Rp  350.000  │      -       │
│  2 │ 1-1400  PPN Masukan (jika PKP)          │ Rp   38.500  │      -       │
│  3 │ 1-1100  Kas                             │      -       │ Rp  388.500  │
├────┼─────────────────────────────────────────┼──────────────┼──────────────┤
│    │ TOTAL                                   │ Rp  388.500  │ Rp  388.500  │
└────┴─────────────────────────────────────────┴──────────────┴──────────────┘
✅ Jurnal Balance: Debit = Kredit = Rp 388.500
```

#### 5.4 CoA Standar & Diperluas (Lengkap)

Tabel berikut adalah daftar akun lengkap setelah penambahan akun untuk jurnal manual. Akun bertanda *(baru)* adalah penambahan dari versi sebelumnya.

| Kode Akun | Nama Akun | Tipe | Keterangan |
|---|---|---|---|
| `1-1100` | Kas | Aset Lancar | Uang tunai di laci kasir |
| `1-1200` | Bank – Rekening Toko | Aset Lancar | Saldo rekening bank toko |
| `1-1250` | Piutang Pembayaran Digital | Aset Lancar | QRIS/GoPay yang belum settle |
| `1-1300` | Persediaan Barang Dagang | Aset Lancar | Nilai stok barang di gudang |
| `1-1400` | PPN Masukan | Aset Lancar | Pajak masukan yang dapat dikreditkan |
| `1-1500` | Biaya Dibayar Dimuka | Aset Lancar | Sewa/asuransi yang dibayar di muka *(baru)* |
| `1-2000` | Aset Tetap – Peralatan | Aset Tetap | Rak, kulkas, mesin kasir, dst. *(baru)* |
| `1-2100` | Akumulasi Depresiasi – Peralatan | Aset Tetap (kontra) | Akumulasi penyusutan aset tetap *(baru)* |
| `2-1100` | Hutang Dagang (AP) | Liabilitas Lancar | Hutang ke supplier |
| `2-1200` | Hutang Lain-Lain | Liabilitas Lancar | Hutang jasa/servis yang belum dibayar *(baru)* |
| `2-1300` | PPN Keluaran | Liabilitas Lancar | PPN yang harus disetorkan ke negara |
| `3-1000` | Modal Pemilik | Ekuitas | Modal yang disetor oleh pemilik |
| `3-1100` | Laba Ditahan | Ekuitas | Akumulasi laba/rugi periode lalu |
| `3-2000` | Prive / Penarikan Modal | Ekuitas (kontra) | Pengambilan dana toko oleh owner *(baru)* |
| `4-1000` | Pendapatan Penjualan | Pendapatan | Pendapatan dari penjualan barang |
| `4-1100` | Retur & Potongan Penjualan | Pendapatan (kontra) | Pengurangan pendapatan dari retur |
| `4-2000` | Pendapatan Lain-Lain | Pendapatan | Surplus stok, ganti rugi vendor, dll. |
| `5-1000` | Harga Pokok Penjualan (HPP) | Beban | Biaya barang yang terjual |
| `5-2000` | Beban Gaji & Upah | Beban | Gaji kasir dan karyawan |
| `5-2100` | Beban Sewa Tempat | Beban | Biaya sewa toko bulanan |
| `5-2200` | Beban Listrik & Air | Beban | Biaya utilitas |
| `5-2300` | Beban Pemeliharaan & Perbaikan | Beban | Jasa servis peralatan toko *(baru)* |
| `5-3000` | Beban Penyusutan Stok | Beban | Selisih minus hasil stock opname |
| `5-4000` | Beban Depresiasi Aset | Beban | Penyusutan aset tetap per periode *(baru)* |
| `5-4100` | Rugi Penjualan / Penghapusan Aset | Beban | Kerugian disposal aset tetap *(baru)* |
| `5-5000` | Beban Non-Operasional | Beban | Denda, biaya bank, penalti, dll. *(baru)* |

---

## 6. User Access & Security

### 6.1 Role-Based Access Control (RBAC) Matrix

| Modul / Aksi | Owner | Manajer Toko | Kasir | Admin Gudang |
|---|:---:|:---:|:---:|:---:|
| Dashboard & Laporan Keuangan | ✅ Full | ✅ Full | ❌ | ❌ |
| POS – Transaksi Penjualan | ✅ Full | ✅ Full | ✅ Full | ❌ |
| Void / Batalkan Transaksi POS | ✅ Full | ✅ Dengan PIN | ❌ | ❌ |
| Purchase Order & GR | ✅ Full | ✅ Full | ❌ | ✅ View + GR |
| Manajemen Supplier | ✅ Full | ✅ Full | ❌ | 👁 View Only |
| Stok Opname & Penyesuaian | ✅ Full | ✅ Approve | ❌ | ✅ Input Data |
| Chart of Accounts & Jurnal Otomatis | ✅ Full | 👁 View Only | ❌ | ❌ |
| **Jurnal Manual – Buat & Edit Draft** | ✅ Full | ✅ Full | ❌ | ❌ |
| **Jurnal Manual – Approve & Posting** | ✅ Full | ❌ | ❌ | ❌ |
| **Jurnal Manual – Reverse** | ✅ Full | ❌ | ❌ | ❌ |
| **Template Jurnal Berulang** | ✅ Full | ✅ Full | ❌ | ❌ |
| **Master UoM – Lihat** | ✅ Full | ✅ Full | ❌ | ✅ View Only |
| **Master UoM – Tambah / Edit** | ✅ Full | ✅ Full | ❌ | ❌ |
| **Konversi UoM per Produk** | ✅ Full | ✅ Full | ❌ | ❌ |
| Manajemen User & Role | ✅ Full | ❌ | ❌ | ❌ |
| Konfigurasi Pajak & UoM | ✅ Full | ✅ Full | ❌ | ❌ |
| Audit Trail | ✅ Full | 👁 View Only | ❌ | ❌ |

### 6.2 Ketentuan Keamanan Sistem

- **Autentikasi:** Username + Password (min. 8 karakter, kombinasi huruf & angka) + opsional 2FA via OTP SMS/Email.
- **Sesi:** Token JWT dengan masa berlaku 8 jam. Auto-logout setelah 30 menit tidak aktif.
- **Enkripsi:** Semua data sensitif (password hash bcrypt, NPWP, data keuangan) dienkripsi AES-256 di database.
- **Audit Trail:** Setiap aksi `CREATE` / `UPDATE` / `DELETE` dicatat di tabel `audit_log` dengan field: `user_id`, `action`, `table_name`, `record_id`, `old_value` (JSON), `new_value` (JSON), `ip_address`, `timestamp`.
- **Backup:** Otomatis setiap 24 jam. Data disimpan minimal 90 hari terakhir dengan enkripsi at-rest.
- **Akses Void/Reversal:** Hanya Owner. Membutuhkan konfirmasi alasan tertulis dan PIN khusus 6 digit.

---

## 7. Acceptance Criteria & Success Metrics

### 7.1 Acceptance Criteria per Fitur (Given-When-Then)

| No | Fitur | Acceptance Criteria | Success Metric |
|:---:|---|---|---|
| 1 | **Multi UoM Konversi** | **GIVEN** produk memiliki konversi 1 Karton = 40 Pcs, **WHEN** admin gudang input GR 3 Karton, **THEN** stok bertambah 120 Pcs secara otomatis tanpa input manual. | 0% error konversi dalam 1.000 transaksi uji coba. |
| 2 | **Automated Journaling** | **GIVEN** transaksi penjualan tunai Rp 100.000 (PPN 11%), **WHEN** transaksi dikonfirmasi, **THEN** sistem membuat jurnal: Debit Kas Rp 111.000 / Kredit Pendapatan Rp 100.000 + PPN Keluaran Rp 11.000. | 100% akurasi D=K pada setiap jurnal. 0 jurnal *unbalanced* dalam 30 hari produksi. |
| 3 | **GR 3-Way Matching** | **GIVEN** PO dibuat untuk 100 Pcs Indomie @ Rp 3.000, **WHEN** GR input 95 Pcs, **THEN** sistem flag 5 Pcs selisih dan membutuhkan approval manajer sebelum memperbarui stok. | 100% GR dengan selisih > 2% harus disetujui manajer sebelum posting. |
| 4 | **Laporan Keuangan Real-Time** | **GIVEN** transaksi penjualan dilakukan jam 14:00, **WHEN** manajer membuka dashboard Laba Rugi, **THEN** data terbaru muncul dalam waktu < 5 detik tanpa perlu refresh manual. | Latency laporan < 5 detik. Neraca selalu balance: Aset = Liabilitas + Ekuitas. |
| 5 | **RBAC & Audit Trail** | **GIVEN** kasir mencoba mengakses menu Chart of Accounts, **WHEN** login sebagai role Kasir, **THEN** menu tersembunyi atau *access denied*. Setiap perubahan dicatat di audit log. | 0 *unauthorized access* dalam penetration testing. 100% aksi ter-log di audit trail. |
| 6 | **Performa Sistem POS** | **GIVEN** kasir melakukan scan barcode produk, **WHEN** item ditemukan di database, **THEN** produk tampil di layar kasir dalam waktu < 1 detik. Sistem mampu handle 50 transaksi/jam. | Response time POS < 1 detik. Uptime sistem > 99,5%. Zero data loss saat offline. |

---

### 7.2 KPI & Success Metrics (3 Bulan Post-Launch)

| Metrik | Baseline (Saat Ini) | Target 3 Bulan | Cara Ukur |
|---|---|---|---|
| Akurasi stok vs fisik (stok opname) | ~85% akurasi | **> 98% akurasi** | Perbandingan sistem vs hitungan fisik bulanan |
| Waktu proses PO hingga GR (end-to-end) | 2–3 hari manual | **< 24 jam digital** | Timestamp PO `created_at` vs GR `confirmed_at` |
| Waktu tutup buku bulanan | 3–5 hari manual | **< 1 jam (otomatis)** | Durasi dari akhir bulan ke laporan final tersedia |
| Response time POS (scan to cart) | N/A (manual) | **< 1 detik** | Load test 100 concurrent users |
| User adoption rate (login aktif) | N/A | **> 80%** user terdaftar aktif tiap minggu | DAU/WAU dari sistem analytics |
| Kepuasan pengguna (SUS Score) | N/A | **> 70 / 100** | Kuesioner System Usability Scale post-onboarding |

---

### 7.3 Non-Functional Requirements

| Kategori | Requirement |
|---|---|
| **Performance** | POS response time < 1 detik untuk 95% request. Laporan keuangan render < 5 detik. |
| **Availability** | Uptime minimal **99,5%** per bulan (downtime maksimal ~3,6 jam/bulan). |
| **Scalability** | Sistem harus mampu menangani hingga **10.000 SKU** dan **1.000 transaksi/hari** tanpa penurunan performa. |
| **Offline Mode** *(MVP+)* | POS harus tetap berfungsi saat koneksi internet terputus. Data disinkronkan saat koneksi kembali. |
| **Data Retention** | Data transaksi dan jurnal disimpan minimal **10 tahun** (sesuai peraturan pajak Indonesia UU KUP). |
| **Compliance** | Sistem harus menghasilkan data yang kompatibel dengan format **e-Faktur DJP** untuk pelaporan PPN. |

---

### 7.4 Asumsi & Batasan

#### Asumsi

- Toko beroperasi di wilayah dengan jaringan internet minimal 4G (kecuali fitur offline mode diaktifkan).
- Supplier mayoritas bukan PKP sehingga tidak semua pembelian menghasilkan faktur pajak — sistem harus mendukung GR tanpa faktur pajak.
- Metode perhitungan HPP default menggunakan **Average Cost** (dapat diubah ke FIFO via konfigurasi sistem).

#### Batasan — Out of Scope (MVP)

| Fitur | Rencana |
|---|---|
| Integrasi langsung dengan e-Faktur DJP | Sprint 4+ |
| Fitur multi-cabang / multi-outlet | v2.0 |
| Modul payroll / penggajian karyawan | v2.0 |
| Integrasi marketplace (Tokopedia, Shopee) | v2.5 |
| Aplikasi mobile native (iOS & Android) | v2.0 |

---

## 8. Laporan Histori Barang Keluar Masuk per Bulan

### 8.1 Deskripsi & Tujuan Fitur

Laporan Histori Barang Keluar Masuk (selanjutnya disebut **Laporan Mutasi Stok**) adalah laporan periodik yang menyajikan setiap pergerakan stok — baik masuk maupun keluar — untuk setiap produk dalam rentang waktu tertentu (harian, mingguan, atau bulanan). Laporan ini menjadi **sumber kebenaran tunggal (single source of truth)** untuk rekonsiliasi stok fisik vs sistem dan analisis pola konsumsi barang.

**Tujuan bisnis laporan ini:**

- Memberikan visibilitas penuh kepada Owner dan Manajer atas pergerakan setiap SKU.
- Mendeteksi anomali stok (penurunan mendadak tanpa transaksi penjualan tercatat).
- Menjadi dasar keputusan pemesanan ulang dan negosiasi volume dengan supplier.
- Mendukung proses rekonsiliasi stok opname dengan membandingkan mutasi sistem vs hitungan fisik.
- Memenuhi kebutuhan audit internal maupun eksternal dengan jejak data yang lengkap.

---

### 8.2 Spesifikasi Fungsional

#### 8.2.1 Jenis Mutasi Stok yang Dicatat

Setiap baris pada laporan merepresentasikan satu **event mutasi stok** yang dikelompokkan berdasarkan tipe transaksi pemicunya:

| Kode Tipe | Label Tampilan | Arah Stok | Modul Pemicu | Keterangan |
|---|---|:---:|---|---|
| `IN-GR` | Penerimaan Barang (GR) | ⬆️ Masuk | Purchase / GR | Barang diterima dari supplier, GR sudah confirmed. |
| `IN-RET-CUST` | Retur dari Pelanggan | ⬆️ Masuk | Sales Return | Pelanggan mengembalikan barang ke toko. |
| `IN-ADJ-PLUS` | Penyesuaian Stok (+) | ⬆️ Masuk | Inventory Adj. | Stok opname: fisik lebih banyak dari sistem. |
| `IN-OPENING` | Stok Awal Periode | ⬆️ Masuk | System / Setup | Saldo stok awal saat pertama kali produk diinput. |
| `OUT-POS` | Penjualan Kasir | ⬇️ Keluar | POS | Barang terjual melalui kasir (tunai/non-tunai). |
| `OUT-RET-SUP` | Retur ke Supplier | ⬇️ Keluar | Purchase Return | Barang dikembalikan ke supplier karena cacat/salah kirim. |
| `OUT-ADJ-MIN` | Penyesuaian Stok (−) | ⬇️ Keluar | Inventory Adj. | Stok opname: fisik lebih sedikit dari sistem (susut/hilang). |
| `OUT-DAMAGE` | Barang Rusak / Kadaluarsa | ⬇️ Keluar | Inventory Adj. | Penghapusan barang tidak layak jual. |

---

#### 8.2.2 Parameter Filter & Konfigurasi Laporan

Pengguna dapat mengkonfigurasi laporan melalui panel filter sebelum generate:

| Parameter | Tipe Input | Nilai Default | Deskripsi |
|---|---|---|---|
| **Periode** | Date range picker | Bulan berjalan | Rentang tanggal laporan. Bisa harian, mingguan, bulanan, atau custom range. |
| **Granularitas** | Dropdown | Bulanan | Pilihan: Harian / Mingguan / Bulanan. |
| **Produk** | Multi-select autocomplete | Semua Produk | Filter satu atau beberapa SKU spesifik. |
| **Kategori** | Multi-select dropdown | Semua Kategori | Filter berdasarkan kategori/sub-kategori produk. |
| **Tipe Mutasi** | Checkbox group | Semua Tipe | Filter hanya mutasi masuk, hanya keluar, atau tipe spesifik. |
| **Tampilkan UoM** | Toggle | UoM Jual (Base) | Pilih tampilkan dalam UoM beli (Karton) atau UoM jual (Pcs). |
| **Tampilkan Nilai (Rp)** | Toggle | Aktif | Tampilkan kolom nilai rupiah berdasarkan HPP saat transaksi. |

---

#### 8.2.3 Struktur Output Laporan — Ringkasan per Produk per Bulan

Tampilan utama laporan disajikan dalam format **rangkuman per SKU per periode** (default: bulanan):

| Kolom | Tipe Data | Sumber Data | Keterangan |
|---|---|---|---|
| `Kode SKU` | String | `product.kode_sku` | Kode unik produk. |
| `Nama Produk` | String | `product.nama` | Nama lengkap produk. |
| `Kategori` | String | `product.kategori` | Nama kategori produk. |
| `Stok Awal Periode` | Decimal | Kalkulasi | Saldo stok pada awal bulan (Pcs/Base Unit). |
| `Total Masuk` | Decimal | `SUM(IN-*)` | Total semua mutasi masuk dalam periode. |
| `Total Keluar` | Decimal | `SUM(OUT-*)` | Total semua mutasi keluar dalam periode. |
| `Stok Akhir Periode` | Decimal | Kalkulasi | `Stok Awal + Total Masuk − Total Keluar`. |
| `Nilai Masuk (Rp)` | Currency | HPP × Qty Masuk | Nilai rupiah total barang masuk berdasarkan HPP saat GR. |
| `Nilai Keluar (Rp)` | Currency | HPP × Qty Keluar | Nilai rupiah total barang keluar berdasarkan HPP saat transaksi. |
| `Selisih Stok` | Decimal | Kalkulasi | Stok sistem vs fisik (diisi saat ada stok opname di periode tsb). |

**Contoh tampilan ringkasan bulanan (Maret 2025):**

```
┌─────────────┬──────────────────┬──────────┬────────────┬──────────────┬─────────────┬─────────────┬──────────────────┬─────────────────┐
│ Kode SKU    │ Nama Produk      │ Kategori │ Stok Awal  │ Total Masuk  │ Total Keluar│ Stok Akhir  │ Nilai Masuk (Rp) │ Nilai Keluar (Rp│
├─────────────┼──────────────────┼──────────┼────────────┼──────────────┼─────────────┼─────────────┼──────────────────┼─────────────────┤
│ SKU-001     │ Indomie Goreng   │ Sembako  │ 240 Pcs    │ 480 Pcs      │ 510 Pcs     │ 210 Pcs     │ Rp 1.440.000     │ Rp 1.530.000    │
│ SKU-002     │ Gula Pasir 1 Kg  │ Sembako  │ 85 Kg      │ 150 Kg       │ 120 Kg      │ 115 Kg      │ Rp 2.100.000     │ Rp 1.680.000    │
│ SKU-003     │ Minyak Goreng 1L │ Sembako  │ 48 Btl     │ 120 Btl      │ 135 Btl     │ 33 Btl      │ Rp 1.920.000     │ Rp 2.160.000    │
│ SKU-004     │ Aqua 600ml       │ Minuman  │ 300 Pcs    │ 600 Pcs      │ 580 Pcs     │ 320 Pcs     │ Rp 1.500.000     │ Rp 1.450.000    │
└─────────────┴──────────────────┴──────────┴────────────┴──────────────┴─────────────┴─────────────┴──────────────────┴─────────────────┘
```

---

#### 8.2.4 Tampilan Detail Mutasi per Transaksi (Drill-Down)

Pengguna dapat mengklik baris produk manapun untuk membuka **tampilan detail** seluruh transaksi individual dalam periode tersebut:

| Kolom | Tipe Data | Sumber Data | Keterangan |
|---|---|---|---|
| `Tanggal & Waktu` | DateTime | `transaction.created_at` | Timestamp transaksi terjadi (WIB). |
| `No. Referensi` | String | `transaction.nomor` | Nomor PO, GR, transaksi POS, atau nomor penyesuaian. |
| `Tipe Mutasi` | Enum | `stock_movement.tipe` | Lihat tabel kode tipe di 8.2.1. |
| `Keterangan` | String | `stock_movement.keterangan` | Deskripsi otomatis (misal: "Penjualan kasir – Shift Pagi"). |
| `Qty Masuk` | Decimal | `stock_movement.qty_in` | Jumlah unit masuk. Kosong jika mutasi keluar. |
| `Qty Keluar` | Decimal | `stock_movement.qty_out` | Jumlah unit keluar. Kosong jika mutasi masuk. |
| `Saldo Stok` | Decimal | Running balance | Saldo stok setelah transaksi ini (running total). |
| `HPP per Unit` | Currency | `stock_movement.hpp_saat_itu` | HPP yang digunakan saat transaksi. |
| `Nilai Transaksi (Rp)` | Currency | Qty × HPP | Nilai rupiah mutasi ini. |
| `Dilakukan Oleh` | String | `user.nama` | Nama user yang memicu transaksi. |

**Contoh tampilan detail drill-down Indomie Goreng – Maret 2025:**

```
Tanggal              │ No. Referensi  │ Tipe           │ Keterangan                    │ Masuk │ Keluar │ Saldo │ HPP/Pcs  │ Nilai (Rp)
─────────────────────┼────────────────┼────────────────┼───────────────────────────────┼───────┼────────┼───────┼──────────┼───────────
01/03/2025 08:00     │ GR-2025-0301   │ IN-GR          │ Terima dari CV Maju Supplier   │ +240  │   -    │  240  │ Rp 3.000 │ Rp 720.000
05/03/2025 09:23     │ POS-2025-1042  │ OUT-POS        │ Penjualan kasir – Shift Pagi   │   -   │  -12   │  228  │ Rp 3.000 │ Rp  36.000
05/03/2025 14:11     │ POS-2025-1087  │ OUT-POS        │ Penjualan kasir – Shift Siang  │   -   │   -8   │  220  │ Rp 3.000 │ Rp  24.000
10/03/2025 10:05     │ GR-2025-0310   │ IN-GR          │ Terima dari CV Maju Supplier   │ +240  │   -    │  460  │ Rp 3.000 │ Rp 720.000
15/03/2025 16:44     │ ADJ-2025-0315  │ OUT-ADJ-MIN    │ Stok opname – selisih minus 2  │   -   │   -2   │  458  │ Rp 3.000 │ Rp   6.000
...                  │ ...            │ ...            │ ...                            │  ...  │  ...   │  ...  │ ...      │ ...
31/03/2025 23:59     │ SYSTEM         │ CLOSING        │ Saldo akhir bulan Maret 2025   │   -   │   -    │  210  │    -     │     -
```

---

### 8.3 Logika Kalkulasi Stok

#### 8.3.1 Formula Saldo Running (Kartu Stok Digital)

```
Saldo_Setelah_Transaksi(n) = Saldo_Setelah_Transaksi(n-1) + Qty_In(n) − Qty_Out(n)

Stok_Awal_Periode          = Saldo_Akhir_Periode_Sebelumnya
Stok_Akhir_Periode         = Stok_Awal + Σ(Qty_In) − Σ(Qty_Out) dalam periode
```

#### 8.3.2 Validasi Integritas Data

Sistem wajib memvalidasi bahwa saldo stok **tidak pernah negatif** pada setiap titik transaksi. Jika ada kondisi stok negatif yang terdeteksi (akibat data historis atau bug), sistem menampilkan **flag peringatan** `⚠️ STOK NEGATIF TERDETEKSI` pada baris yang bersangkutan tanpa memblokir tampilan laporan.

```
CONSTRAINT: Saldo_Setelah_Transaksi(n) >= 0  untuk semua n
JIKA dilanggar → tampilkan flag warning, catat ke error_log, notifikasi Manajer
```

#### 8.3.3 Entitas Data Baru: `stock_movement`

Tabel baru yang perlu ditambahkan ke skema database untuk mendukung laporan ini:

```sql
CREATE TABLE stock_movement (
  id              BIGINT PRIMARY KEY AUTO_INCREMENT,
  product_id      INT NOT NULL,
  tanggal         DATETIME NOT NULL,
  tipe            ENUM('IN-GR','IN-RET-CUST','IN-ADJ-PLUS','IN-OPENING',
                       'OUT-POS','OUT-RET-SUP','OUT-ADJ-MIN','OUT-DAMAGE') NOT NULL,
  referensi_id    INT,                   -- FK ke PO/GR/SalesTransaction/AdjustmentID
  referensi_tipe  VARCHAR(50),           -- 'GR' | 'POS' | 'ADJ' | dsb.
  qty_in          DECIMAL(15,4) DEFAULT 0,
  qty_out         DECIMAL(15,4) DEFAULT 0,
  saldo_setelah   DECIMAL(15,4) NOT NULL, -- running balance
  hpp_per_unit    DECIMAL(15,4),          -- snapshot HPP saat transaksi
  nilai_rp        DECIMAL(15,2),          -- qty × hpp_per_unit
  keterangan      VARCHAR(255),
  created_by      INT NOT NULL,           -- FK ke users
  created_at      DATETIME DEFAULT NOW(),

  INDEX idx_product_tanggal (product_id, tanggal),
  INDEX idx_tipe (tipe),
  INDEX idx_referensi (referensi_tipe, referensi_id)
);
```

> **Catatan Engineering:** Tabel `stock_movement` bersifat **append-only** — tidak ada UPDATE atau DELETE. Setiap koreksi menghasilkan baris baru dengan tipe `IN-ADJ-PLUS` atau `OUT-ADJ-MIN`. Ini memastikan jejak audit stok 100% dapat ditelusuri.

---

### 8.4 Tampilan UI & UX

#### 8.4.1 Halaman Utama Laporan Mutasi Stok

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  📦 Laporan Mutasi Stok                              [Export Excel] [Export PDF]│
├─────────────────────────────────────────────────────────────────────────────┤
│  Periode: [Maret 2025 ▼]  Produk: [Semua ▼]  Kategori: [Semua ▼]          │
│  Tipe: [✅ Masuk  ✅ Keluar]  UoM: [Pcs ▼]  Tampilkan Nilai: [ON]  [Generate]│
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────── Ringkasan Bulan ────────────────────────┐│
│  │  📥 Total Masuk: 1.590 Pcs   📤 Total Keluar: 1.345 Pcs               ││
│  │  💰 Nilai Masuk: Rp 6.960.000  💸 Nilai Keluar: Rp 6.820.000          ││
│  └──────────────────────────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────────────────────────┤
│  [Tabel Ringkasan per SKU — lihat contoh di 8.2.3]                          │
│  Klik baris → buka drawer detail transaksi (drill-down)                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 8.4.2 Ketentuan UX

- **Kode warna mutasi:** Baris masuk ditampilkan dengan aksen **hijau**, baris keluar dengan aksen **merah/oranye**, penyesuaian dengan aksen **kuning**.
- **Highlight anomali:** Baris dengan selisih stok > 5% dari total keluar ditampilkan dengan latar kuning dan ikon `⚠️`.
- **Sticky header:** Kolom `Kode SKU` dan `Nama Produk` dikunci saat scroll horizontal pada tabel lebar.
- **Pagination:** Default 50 baris per halaman. Opsi tampil 100/200 baris.
- **Grafik opsional:** Toggle tampilan grafik batang (bar chart) mutasi masuk vs keluar per produk per minggu dalam bulan tersebut.

---

### 8.5 Hak Akses Laporan Mutasi Stok

| Role | Hak Akses | Batasan |
|---|:---:|---|
| **Owner** | ✅ Full | Dapat melihat semua produk, semua periode, export semua format. |
| **Manajer Toko** | ✅ Full | Dapat melihat semua produk dan periode. Dapat export. |
| **Kasir** | ❌ Tidak Ada | Tidak memiliki akses ke laporan stok manapun. |
| **Admin Gudang** | 👁 Terbatas | Dapat melihat laporan stok (view only). Tidak dapat export ke PDF/Excel. |

---

### 8.6 Format Ekspor

| Format | Konten yang Disertakan | Keterangan |
|---|---|---|
| **Excel (.xlsx)** | Sheet 1: Ringkasan per SKU. Sheet 2: Detail semua transaksi. Sheet 3: Pivot per kategori. | Cocok untuk analisis lanjutan oleh owner atau akuntan. |
| **PDF** | Ringkasan per SKU + header toko + tanda tangan manajer (untuk audit fisik). | Cocok untuk laporan tercetak / arsip bulanan. |
| **CSV** | Raw data detail transaksi tanpa formatting. | Untuk keperluan import ke sistem akuntansi eksternal. |

---

### 8.7 Acceptance Criteria – Laporan Mutasi Stok

| No | Skenario | Given | When | Then | Metric |
|:---:|---|---|---|---|---|
| 1 | **Akurasi Running Balance** | Produk memiliki 500 transaksi dalam satu bulan | Pengguna generate laporan detail | Setiap baris `saldo_setelah` = baris sebelumnya + qty_in − qty_out, tanpa exception | 0 error kalkulasi dalam 10.000 baris uji. |
| 2 | **Kelengkapan Tipe Mutasi** | Dalam satu bulan terjadi GR, penjualan POS, retur, dan stok opname | Generate laporan dengan filter "Semua Tipe" | Semua 8 tipe mutasi tercatat dengan tipe yang benar, tidak ada yang missing | 100% event stok terklasifikasi. |
| 3 | **Performa Generate** | Produk memiliki 10.000 baris mutasi dalam satu bulan | Pengguna klik tombol Generate | Laporan ringkasan muncul dalam < 3 detik. Detail drill-down muncul < 5 detik | P95 response time < 3 detik. |
| 4 | **Konsistensi dengan Jurnal** | Setiap mutasi stok memiliki jurnal akuntansi yang sesuai | Bandingkan `stock_movement` dengan `journal_entry` bulan yang sama | Nilai total `OUT-POS` pada laporan stok = nilai HPP pada jurnal akun 5-1000 periode yang sama | Selisih rekonsiliasi = Rp 0. |
| 5 | **Ekspor Excel** | Pengguna memfilter laporan Maret 2025 untuk 3 kategori | Klik Export Excel | File .xlsx terunduh dalam < 10 detik, berisi 3 sheet dengan data yang identik dengan tampilan layar | 100% data match antara layar dan file. |
| 6 | **Deteksi Stok Negatif** | Terdapat data historis dengan saldo stok negatif (bug lama) | Generate laporan detail | Baris dengan saldo negatif ditampilkan dengan flag `⚠️ STOK NEGATIF` berwarna merah, laporan tetap dapat dilihat | 100% anomali ter-flag, 0 crash. |

---

*— Akhir Dokumen PRD v1.2 —*

*Dokumen ini bersifat **CONFIDENTIAL** dan hanya untuk konsumsi internal tim produk.*
