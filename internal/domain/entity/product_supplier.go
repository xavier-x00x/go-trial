package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductSupplier merepresentasikan katalog barang milik Supplier tertentu (Many-to-Many).
// Tabel ini sangat penting agar staf Purchasing tidak salah memesan barang ke supplier yang salah.
type ProductSupplier struct {
	BaseModel
	ProductID           uuid.UUID       `gorm:"type:char(36);not null;index;uniqueIndex:idx_product_supplier_store" json:"product_id"`
	Product             Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"` // Relasi ke Produk Internal
	SupplierID          uuid.UUID       `gorm:"type:char(36);not null;index;uniqueIndex:idx_product_supplier_store" json:"supplier_id"`
	Supplier            Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`                            // Relasi ke Vendor/Distributor
	StoreID             *uuid.UUID      `gorm:"type:char(36);index;uniqueIndex:idx_product_supplier_store" json:"store_id"` // Jika NULL = Kontrak Nasional. Jika terisi = Kontrak Harga Khusus Cabang tsb.
	Store               *Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`                                  // Relasi ke Cabang
	SupplierSKU         *string         `gorm:"type:varchar(50)" json:"supplier_sku"`                                       // Kode barang versi pabrik/distributor (Biasanya berbeda dengan SKU internal kita)
	IsPrimary           bool            `gorm:"default:false" json:"is_primary"`                                            // Flag penanda apakah ini distributor utama (prioritas pertama saat re-order)
	IsConsignment       bool            `gorm:"default:false" json:"is_consignment"`                                        // True = Barang Titipan (Konsinyasi). Hutang tercipta hanya jika barang laku terjual.
	IsReturnable        bool            `gorm:"default:true" json:"is_returnable"`                                          // True = Barang bisa di-retur (Purchase Return) ke supplier jika rusak/expired.
	DefaultLeadTimeDays int             `gorm:"default:0" json:"default_lead_time_days"`                                    // Estimasi waktu pengiriman default (hari) secara nasional
	OfferedPrice        decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"offered_price"`                          // Harga beli penawaran/kontrak dari supplier ini
	PurchaseUOMID       *uuid.UUID      `gorm:"type:char(36);index" json:"purchase_uom_id"`                                 // Jika NULL, gunakan BaseUOM Produk.
	PurchaseUOM         *UOM            `gorm:"foreignKey:PurchaseUOMID" json:"purchase_uom,omitempty"`
	MinOrderQty         decimal.Decimal `gorm:"type:decimal(15,3);default:1" json:"min_order_qty"`                          // Minimum kuantitas pesanan (MOQ) untuk mendapat harga tersebut
}

/*
	Saat staf Gudang menerima barang dari truk (Dokumen Goods Receipt), sistem Backend akan mengecek flag ini:

	Jika IsConsignment == false (Beli Putus): Sistem membuat jurnal Hutang (Account Payable) saat itu juga.

	Jika IsConsignment == true (Titip Jual): Sistem hanya menambah stok fisik di InventoryStock, namun TIDAK mencatat jurnal Hutang.

	Hutang untuk barang titipan ini baru akan otomatis digenerate oleh sistem di malam hari (saat Tutup Kasir), dengan menghitung khusus barang titipan apa saja yang sudah terjual ke konsumen hari itu.

	Dengan adanya flag ini di level Master Data, tim Akuntansi Anda tidak akan pernah lagi sakit kepala menghitung selisih laporan hutang barang titipan di akhir bulan.
*/

/*
	Cara Menghitung Safety Stock (Stok Pengaman):
	Safety Stock = (Max Lead Time dalam hari - Rata-rata Lead Time dalam hari) x Rata-rata Penjualan Harian

	Contoh Kasus:
	- Produk A: Rata-rata penjualan 10 pcs/hari, Lead Time 3 hari.
	- Produk B: Rata-rata penjualan 10 pcs/hari, Lead Time 5 hari (Supplier kadang telat).

	Perhitungan Safety Stock Produk B:
	Diasumsikan variasi lead time: (5 hari - 3 hari) x 10 pcs/hari = 20 pcs. (Safety Stock Bisa Manual, karena hanya asumsi)

	Reorder Point Produk B = (10 pcs/hari x 5 hari) + 20 pcs = 70 pcs.

	Logika Bisnis:
	Jika stok Produk B sudah mencapai 70 pcs, sistem harus segera mengirim Purchase Order (PO) ke supplier.
	Ini untuk memastikan barang tidak akan habis selama 2 hari ekstra (antisipasi keterlambatan pengiriman dari 3 hari normal menjadi 5 hari).

	Dengan adanya LeadTimeDays di tabel ini, modul Purchasing dan Inventory Anda bisa memprediksi kapan harus memesan barang agar toko tidak "Kehabisan Stok" (Stockout) di tengah hari.
*/
