package entity

import (
	"time"

	"github.com/google/uuid"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Enum untuk MasterDataProposal
// ──────────────────────────────────────────────────────────────────────────────

// ProposalEntityType mendefinisikan jenis entity master data yang memerlukan persetujuan.
// Entity yang TIDAK ada di daftar ini (UOM, Category, Store, Role, dll) menggunakan CRUD biasa.
const (
	ProposalEntityProduct         = "PRODUCT"
	ProposalEntityProductPrice    = "PRODUCT_PRICE"
	ProposalEntityProductUOM      = "PRODUCT_UOM_CONVERSION"
	ProposalEntitySupplier        = "SUPPLIER"
	ProposalEntityProductSupplier = "PRODUCT_SUPPLIER"
	ProposalEntityChartOfAccount  = "CHART_OF_ACCOUNT"
	ProposalEntityTax             = "TAX"
)

// ProposalActionType mendefinisikan jenis operasi yang diusulkan.
const (
	ProposalActionCreate = "CREATE"
	ProposalActionUpdate = "UPDATE"
	ProposalActionDelete = "DELETE"
)

// ProposalStatus mendefinisikan status siklus hidup proposal.
const (
	ProposalStatusPending  = "PENDING"
	ProposalStatusApproved = "APPROVED"
	ProposalStatusRejected = "REJECTED"
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: Header (Dokumen Induk)
// ──────────────────────────────────────────────────────────────────────────────

// MasterDataProposal adalah dokumen induk (Header) yang merangkum satu usulan
// perubahan data master. Satu dokumen bisa memiliki banyak item detail.
//
// Alur Kerja (Maker-Checker):
//  1. Staf (Maker) mengajukan perubahan → Header + Items masuk dengan Status = PENDING.
//  2. Manajer (Checker) membuka dokumen, meninjau daftar item di dalamnya.
//  3. Jika Approve → Backend memproses setiap item, mengeksekusi perubahan ke tabel asli.
//  4. Jika Reject → Status berubah menjadi REJECTED beserta ReviewNotes sebagai alasan.
//
// Halaman List/Index cukup query tabel ini saja (tanpa GROUP BY).
type MasterDataProposal struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	ReferenceNumber string `gorm:"type:varchar(50);uniqueIndex;not null" json:"reference_number"` // Nomor Bukti unik per dokumen (cth: PRD/2026/04/0001)
	EntityType      string `gorm:"type:varchar(30);not null;index" json:"entity_type"`            // Jenis entity target (cth: PRODUCT, SUPPLIER)
	ActionType      string `gorm:"type:varchar(10);not null" json:"action_type"`                  // Operasi yang diusulkan: CREATE, UPDATE, DELETE
	TotalItems      int    `gorm:"default:0" json:"total_items"`                                  // Jumlah item dalam dokumen (denormalisasi untuk performa list)

	// ── Status & Tracking ────────────────────────────────────────────────
	Status       string     `gorm:"type:varchar(10);not null;default:'PENDING';index" json:"status"` // Siklus hidup: PENDING → APPROVED / REJECTED
	ProposedByID uuid.UUID  `gorm:"type:char(36);not null;index" json:"proposed_by_id"`              // ID User yang mengajukan usulan (Maker)
	ProposedBy   User       `gorm:"foreignKey:ProposedByID" json:"proposed_by,omitempty"`            // Relasi ke tabel User
	ReviewedByID *uuid.UUID `gorm:"type:char(36);index" json:"reviewed_by_id"`                       // ID User yang mereview (Checker). NULL selama PENDING.
	ReviewedBy   *User      `gorm:"foreignKey:ReviewedByID" json:"reviewed_by,omitempty"`            // Relasi ke tabel User
	ReviewedAt   *time.Time `json:"reviewed_at"`                                                     // Waktu keputusan review diambil

	// ── Catatan ──────────────────────────────────────────────────────────
	Reason      string  `gorm:"type:text" json:"reason"`       // Alasan/justifikasi pengajuan oleh Maker
	ReviewNotes *string `gorm:"type:text" json:"review_notes"` // Catatan dari Checker (terutama saat REJECTED, wajib diisi alasan penolakan)

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []MasterDataProposalItem `gorm:"foreignKey:ProposalID" json:"items,omitempty"` // Daftar item detail dalam dokumen ini
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: Detail (Item per Baris)
// ──────────────────────────────────────────────────────────────────────────────

// MasterDataProposalItem adalah baris detail dari dokumen proposal.
// Setiap item merepresentasikan satu record entity yang akan di-Create/Update/Delete.
//
// Contoh:
//   - Bulk Create 5 Product → 1 Header + 5 Item (masing-masing berisi PayloadJSON dari CreateProductRequest)
//   - Update 1 Supplier     → 1 Header + 1 Item (berisi PayloadJSON + SnapshotJSON untuk diff view)
type MasterDataProposalItem struct {
	BaseModel
	ProposalID   uuid.UUID          `gorm:"type:char(36);not null;index" json:"proposal_id"` // FK ke Header
	Proposal     MasterDataProposal `gorm:"foreignKey:ProposalID" json:"proposal,omitempty"` // Relasi ke Header
	SeqNo        int                `gorm:"not null" json:"seq_no"`                          // Nomor urut item dalam dokumen (1, 2, 3, ...)
	EntityID     *uuid.UUID         `gorm:"type:char(36);index" json:"entity_id"`            // ID record target. NULL jika ActionType = CREATE
	PayloadJSON  string             `gorm:"type:json;not null" json:"payload_json"`          // JSON berisi data request DTO yang diusulkan
	SnapshotJSON *string            `gorm:"type:json" json:"snapshot_json"`                  // JSON snapshot kondisi SEBELUM diubah (audit trail & diff view). NULL jika CREATE.
}
