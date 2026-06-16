# Dokumen Rancangan Database: Sistem ERP Toko Kelontong Pro
**Fase 1: Data Master (Master Data)**

## 1. Pendahuluan
Dokumen ini mendefinisikan skema arsitektur database relasional (RDBMS) untuk "Sistem ERP Toko Kelontong Pro". Fokus pada fase ini adalah **Data Master** yang mencakup manajemen entitas inti seperti inventaris, pemasok, pelanggan, konfigurasi toko, dan fondasi akuntansi.

Desain ini menggunakan pendekatan *System-of-Record* dengan konvensi penamaan standar industri (`snake_case` untuk tabel dan kolom) dan mendukung *soft deletes* untuk menjaga integritas historis dan audit log pada data finansial.

### Konvensi dan Standar Tipe Data
- `id`: Tipe data `UUID` v4 direkomendasikan untuk keamanan dan skalabilitas sistem terdistribusi, menghindari enumerasi ID yang mudah ditebak.
- `created_at`, `updated_at`: `TIMESTAMP WITH TIME ZONE` (TIMESTAMPTZ) untuk mencatat waktu operasional secara presisi.
- `deleted_at`: `TIMESTAMP WITH TIME ZONE` bernilai `NULL` (Implementasi *Soft Delete*).
- `created_by`, `updated_by`, `deleted_by`: `UUID` yang merujuk ke tabel pengguna untuk *Audit Trail* yang akuntabel.
- **Finansial (Harga/Uang)**: Menggunakan `DECIMAL(19,4)` untuk akurasi finansial tinggi (menghindari *floating-point error* yang fatal di aplikasi akuntansi).

---

## 2. Skema Tabel Data Master

### 2.1. Otorisasi dan Manajemen Pengguna

Manajemen akses untuk kasir, kepala toko, petugas gudang, dan manajer akuntansi.

#### Tabel: `roles` (Peran Pengguna)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier unik peran |
| `name` | `VARCHAR(50)` | UNIQUE, NOT NULL | Nama peran (cth: `SuperAdmin`, `Cashier`, `Manager`) |
| `description` | `TEXT` | NULL | Deskripsi fungsional cakupan peran |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | `created_at`, `updated_at`, `deleted_at` dll |

#### Tabel: `users` (Pengguna/Karyawan)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier unik pengguna |
| `role_id` | `UUID` | FK (`roles.id`), NOT NULL | Hak akses / peran pengguna dalam sistem ERP |
| `username` | `VARCHAR(50)` | UNIQUE, NOT NULL | *Username* untuk otentikasi kasir/karyawan |
| `email` | `VARCHAR(100)`| UNIQUE, NOT NULL | Email valid pengguna |
| `password_hash` | `VARCHAR(255)`| NOT NULL | *Hash* kata sandi menggunakan algoritma kuat (cth: `Argon2id` / `Bcrypt`) |
| `full_name` | `VARCHAR(100)`| NOT NULL | Nama entitas karyawan sesungguhnya |
| `is_active` | `BOOLEAN` | DEFAULT TRUE | Otomasi pemblokiran akses tanpa menghapus data (`FALSE` jika *resign*) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

---

### 2.2. Manajemen Toko dan Lokasi

Memfasilitasi arsitektur multi-cabang (skalabilitas).

#### Tabel: `stores` (Toko/Cabang Utama)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier cabang toko |
| `code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Kode unik identifikasi internal toko (cth: `TK-PUSAT`, `TK-01`) |
| `name` | `VARCHAR(100)`| NOT NULL | Nama legal/publik toko cabang |
| `address` | `TEXT` | NOT NULL | Alamat lokasi fisik |
| `tax_id` | `VARCHAR(50)` | NULL | NPWP Cabang / Toko (Bila ada entitas PKP terpisah) |
| `default_price_list_id` | `UUID` | FK (`price_lists.id`), NULL | **Daftar harga standar** yang dipakai oleh toko ini jika tidak ada promo/member |
| `is_active` | `BOOLEAN` | DEFAULT TRUE | Status operasional cabang (Aktif/Tutup permanen) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `warehouses` (Gudang/Lokasi Penyimpanan)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier gudang |
| `store_id` | `UUID` | FK (`stores.id`), NOT NULL | Relasi ke toko kepemilikan (Pemisahan inventaris antar toko) |
| `code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Kode referensi gudang (cth: `WH-DEPAN`, `WH-BEKU`) |
| `name` | `VARCHAR(100)`| NOT NULL | Nama deskriptif gudang penyimpanan |
| `is_active` | `BOOLEAN` | DEFAULT TRUE | Status gudang aktif menerima barang |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

---

### 2.3. Manajemen Produk dan Inventaris Katalog

Sistem ini didesain mendukung produk eceran maupun grosir (multi-satuan).

#### Tabel: `uoms` (Unit of Measurement / Satuan)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier UoM |
| `code` | `VARCHAR(10)` | UNIQUE, NOT NULL | Kode referensi teknis (cth: `PCS`, `KG`, `DUS`, `BALL`) |
| `name` | `VARCHAR(50)` | NOT NULL | Nama standar satuan (cth: Pieces, Kilogram, Karton) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `product_categories` (Kategori Produk)
Menggunakan pola *Adjacency List* (`parent_id`) untuk hierarki kategori yang tak terbatas.
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier kategori |
| `parent_id` | `UUID` | FK (`product_categories.id`), NULL| Kategori induk. Bernilai `NULL` jika kategori utama (*Root Category*) |
| `name` | `VARCHAR(100)`| NOT NULL | Contoh: "Sembako", "Minuman", "Perawatan Tubuh" |
| `slug` | `VARCHAR(120)`| UNIQUE, NOT NULL | SEO/URL Friendly identifier (cth: `minuman-ringan`) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `products` (Produk/Barang Utama)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier produk |
| `sku` | `VARCHAR(50)` | UNIQUE, NOT NULL | *Stock Keeping Unit* (Identifikasi internal pergudangan) |
| `barcode` | `VARCHAR(50)` | UNIQUE, NULL | Nomor Barcode EAN-13 / UPC (Untuk integrasi mesin *scanner* POS) |
| `name` | `VARCHAR(200)`| NOT NULL | Nama lengkap entitas produk |
| `category_id` | `UUID` | FK (`product_categories.id`) | Pemetaan pada kategori hierarkis |
| `base_uom_id` | `UUID` | FK (`uoms.id`), NOT NULL | Satuan terkecil yang digunakan untuk tracking kartu stok (cth: `PCS`) |
| `is_stockable` | `BOOLEAN` | DEFAULT TRUE | `TRUE` untuk barang fisik, `FALSE` untuk item tipe jasa (biaya antar, dll) |
| `min_stock_alert`| `DECIMAL(15,2)`| DEFAULT 0 | Ambang batas *Reorder Point* (Notifikasi untuk *Purchasing*) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `product_uom_conversions` (Konversi Multi-Satuan)
Sangat krusial untuk fitur Grosir (Contoh: Beli Pemasok 1 Dus, Jual 24 Pcs).
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier konversi |
| `product_id` | `UUID` | FK (`products.id`), NOT NULL | Relasi ke produk spesifik |
| `uom_id` | `UUID` | FK (`uoms.id`), NOT NULL | Satuan grosir/alternatif (cth: `DUS`, `RENCENG`) |
| `conversion_rate`| `DECIMAL(15,4)`| NOT NULL | *Multiplier*. Contoh jika uom = DUS dan base_uom = PCS, nilainya `24.0000` |
| `barcode` | `VARCHAR(50)` | UNIQUE, NULL | Barcode spesifik untuk kemasan grosir ini |

#### Tabel: `price_lists` (Buku Harga / Daftar Harga)
Untuk mengakomodir harga yang berbeda berdasarkan cabang, jenis pelanggan (retail vs grosir), atau masa promo.
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier price list |
| `code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Kode buku harga (cth: `RETAIL-01`, `GROSIR-A`) |
| `name` | `VARCHAR(100)`| NOT NULL | Nama daftar harga (cth: Harga Eceran Standar, Harga Reseller) |
| `currency_code`| `VARCHAR(3)` | DEFAULT 'IDR'| Kode mata uang (standar ISO 4217) |
| `is_active` | `BOOLEAN` | DEFAULT TRUE | Status ketersediaan *price list* |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `product_prices` (Harga Produk per Daftar Harga)
Memisahkan harga jual dari tabel `products` agar lebih fleksibel (*Dynamic Pricing*).
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier harga produk |
| `price_list_id`| `UUID` | FK (`price_lists.id`), NOT NULL| Relasi ke daftar harga |
| `product_id` | `UUID` | FK (`products.id`), NOT NULL | Relasi ke entitas produk |
| `uom_id` | `UUID` | FK (`uoms.id`), NOT NULL | Satuan spesifik harga (Contoh: Harga 1 PCS berbeda margin dengan 1 DUS) |
| `sell_price` | `DECIMAL(19,4)`| NOT NULL | Harga Dasar Penjualan (Sebelum dipotong diskon oleh Promotion Engine) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `inventory_stocks` (Saldo Stok Gudang / Master Inventory)
Melacak saldo ketersediaan barang secara *real-time* di setiap lokasi fisik secara presisi (diukur dalam `base_uom`).
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier saldo inventaris |
| `warehouse_id` | `UUID` | FK (`warehouses.id`), NOT NULL| Relasi ke gudang penyimpan |
| `product_id` | `UUID` | FK (`products.id`), NOT NULL | Relasi ke produk yang disimpan |
| `quantity` | `DECIMAL(15,4)`| DEFAULT 0.0000 | Saldo ketersediaan fisik stok saat ini |
| `reserved_qty` | `DECIMAL(15,4)`| DEFAULT 0.0000 | Stok yang sudah di-*booking* transaksi penjualan tapi belum dikirim |
| `average_buy_price` | `DECIMAL(19,4)`| DEFAULT 0.0000 | HPP / Harga Pokok Pembelian Rata-rata Bergerak per gudang |
| `last_audit_at`| `TIMESTAMPTZ`| NULL | Waktu pencatatan terakhir melalui *Stock Opname* |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `store_product_assortments` (Katalog Jajaran Produk Cabang)
Digunakan untuk pengaturan eksplisit skala besar guna memblokir atau mengizinkan penjualan barang tertentu di cabang tertentu, terlepas dari ketersediaan harga dan stok.
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier pengaturan *assortment* |
| `store_id` | `UUID` | FK (`stores.id`), NOT NULL| Relasi ke cabang/toko spesifik |
| `product_id` | `UUID` | FK (`products.id`), NOT NULL| Relasi ke entitas produk |
| `status` | `ENUM` | NOT NULL | Status ketersediaan katalog: `ACTIVE`, `DISCONTINUED`, atau `NOT_ASSORTED` |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

---

### 2.4. Entitas Bisnis Eksternal

#### Tabel: `suppliers` (Pemasok/Vendor)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier pemasok |
| `code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Kode register vendor ERP |
| `name` | `VARCHAR(150)`| NOT NULL | Nama entitas bisnis distributor/pabrik |
| `contact_person` | `VARCHAR(100)`| NULL | Person In Charge (PIC) dari supplier |
| `phone_number` | `VARCHAR(20)` | NULL | Kontak seluler / kantor |
| `tax_id` | `VARCHAR(50)` | NULL | NPWP Supplier untuk keperluan faktur |
| `ap_account_id` | `UUID` | FK (`chart_of_accounts.id`), NULL| **[FINTECH]** Akun Utang Usaha spesifik untuk penjurnalan otomatis |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `customers` (Pelanggan / Member)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier pelanggan / member loyalitas |
| `code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Kode referensi pelanggan (Digunakan di kartu *member*) |
| `name` | `VARCHAR(150)`| NOT NULL | Nama individu / institusi |
| `phone_number` | `VARCHAR(20)` | UNIQUE, NULL | Nomor WhatsApp integrasi (*CRM/Loyalty Program*) |
| `point_balance` | `DECIMAL(15,2)`| DEFAULT 0 | Saldo poin akumulatif (*Customer Retention*) |
| `credit_limit` | `DECIMAL(19,4)`| DEFAULT 0 | Plafon piutang untuk transaksi grosir / kasbon |
| `ar_account_id` | `UUID` | FK (`chart_of_accounts.id`), NULL| **[FINTECH]** Akun Piutang Usaha spesifik (bila pelanggan korporat) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

---

### 2.5. Fondasi Akuntansi (Chart of Accounts & Keuangan)

Berfungsi sebagai motor penjurnalan ganda (*Double-Entry Accounting*) di balik layar operasional.

#### Tabel: `chart_of_accounts` (Bagan Akun - CoA)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier CoA |
| `account_code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Nomor akun konvensi akuntansi (cth: `1110.01`, `5100.00`) |
| `name` | `VARCHAR(100)`| NOT NULL | Nama entitas akun (cth: "Kas Laci Bawah", "HPP Grosir") |
| `account_type` | `ENUM` | NOT NULL | Enum: `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE` |
| `normal_balance` | `ENUM` | NOT NULL | Sifat penambahan saldo normal: `DEBIT` atau `CREDIT` |
| `is_active` | `BOOLEAN` | DEFAULT TRUE | Non-aktifkan jika akun tidak lagi digunakan untuk menjurnal |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `taxes` (Pajak Transaksional)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier setup pajak |
| `name` | `VARCHAR(50)` | NOT NULL | Label pajak di struk/invoice (cth: "PPN 11%", "PB1") |
| `rate_percentage`| `DECIMAL(6,4)` | NOT NULL | Besaran tarif. PPN 11% disimpan `11.0000` |
| `tax_account_id` | `UUID` | FK (`chart_of_accounts.id`), NOT NULL| **[FINTECH]** Mapping ke Akun Liabilitas Pajak (Hutang Pajak PPN Keluaran) |
| *Audit Trails* | `TIMESTAMPTZ`, `UUID`| | Standard audit trails |

#### Tabel: `payment_methods` (Metode Pembayaran Kasir)
| Kolom | Tipe Data | Constraint | Deskripsi |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK | Identifier *Payment Gateway/Method* |
| `code` | `VARCHAR(20)` | UNIQUE, NOT NULL | Identifier internal (cth: `CASH`, `QRIS_GOPAY`, `DEBIT_BCA`) |
| `name` | `VARCHAR(50)` | NOT NULL | Nama metode bayar yang tampil di layar POS |
| `mdr_percentage` | `DECIMAL(5,2)` | DEFAULT 0.00 | *Merchant Discount Rate* (Potongan admin bank payment gateway) |
| `deposit_account_id`| `UUID` | FK (`chart_of_accounts.id`), NOT NULL| **[FINTECH]** Akun penerima kas (Bank BCA, Kas Laci, dll) |
| `expense_account_id`| `UUID` | FK (`chart_of_accounts.id`), NULL| **[FINTECH]** Akun beban admin (untuk pencatatan potongan MDR otomatis) |
| `is_active` | `BOOLEAN` | DEFAULT TRUE | Ketersediaan di terminal kasir |

## 3. Catatan Arsitektur Engineer (Best Practices)
1. **Pemisahan Harga (Pricing)**: Anda mungkin menyadari tidak ada kolom harga beli/jual di tabel `products`. Ini disengaja. Di sistem ERP Pro, harga dikelola di tabel terpisah (misal: `price_lists` atau *Price Book*) karena harga berubah seiring waktu, lokasi toko, dan segmen pelanggan (Retail vs Grosir).
2. **Akuntansi Terintegrasi**: Master data finansial (`payment_methods`, `taxes`, `suppliers`) telah dipasangkan ke tabel CoA. Saat kasir menekan "Bayar QRIS", *backend* dapat dengan mulus mendebit akun *Deposit Account*, mendebit akun *MDR Expense*, dan mengkredit *Revenue Account* secara simultan dalam 1 transaksi ACID database.
3. **Konversi Satuan yang Presisi**: Dengan tabel `product_uom_conversions`, manajemen inventori selalu dicatat (*Stock Card*) dalam unit terkecil (`base_uom`), mencegah desimal persediaan yang tidak konsisten saat jual grosir namun eceran tidak laku.
