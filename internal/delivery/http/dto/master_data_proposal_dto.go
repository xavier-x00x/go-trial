package dto

import (
	"time"

	"github.com/google/uuid"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request DTOs
// ──────────────────────────────────────────────────────────────────────────────

// CreateMasterDataProposalRequest digunakan oleh Maker (Staf) untuk mengajukan
// satu dokumen usulan. Items bisa berisi 1 (single entry) atau banyak (bulk entry).
//
// Contoh Single Entry (Update 1 Supplier):
//
//	{ "entity_type": "SUPPLIER", "action_type": "UPDATE", "reason": "Update alamat",
//	  "items": [{ "entity_id": "uuid-supplier", "payload_json": "{...}" }] }
//
// Contoh Bulk Entry (Create 5 Product):
//
//	{ "entity_type": "PRODUCT", "action_type": "CREATE", "reason": "Barang baru Q2",
//	  "items": [{ "payload_json": "{...}" }, { "payload_json": "{...}" }, ...] }
type CreateMasterDataProposalRequest struct {
	EntityType string                              `json:"entity_type" validate:"required,oneof=PRODUCT PRODUCT_PRICE PRODUCT_UOM_CONVERSION SUPPLIER PRODUCT_SUPPLIER CHART_OF_ACCOUNT TAX"`
	ActionType string                              `json:"action_type" validate:"required,oneof=CREATE UPDATE DELETE"`
	Reason     string                              `json:"reason" validate:"omitempty"`
	Items      []CreateMasterDataProposalItemInput `json:"items" validate:"required,min=1,dive"`
}

// UpdateMasterDataProposalRequest digunakan untuk mengubah proposal yang masih PENDING.
type UpdateMasterDataProposalRequest struct {
	Reason string                              `json:"reason" validate:"omitempty"`
	Items  []CreateMasterDataProposalItemInput `json:"items" validate:"required,min=1,dive"`
}

// CreateMasterDataProposalItemInput adalah input per baris item dalam satu dokumen proposal.
type CreateMasterDataProposalItemInput struct {
	EntityID    *uuid.UUID `json:"entity_id"`                        // Wajib diisi jika ActionType = UPDATE atau DELETE
	PayloadJSON string     `json:"payload_json" validate:"required"` // JSON string berisi data request DTO
}

// ReviewMasterDataProposalRequest digunakan oleh Checker (Manajer) untuk
// menyetujui atau menolak dokumen usulan secara keseluruhan.
type ReviewMasterDataProposalRequest struct {
	Action      string  `json:"action" validate:"required,oneof=APPROVE REJECT"`
	ReviewNotes *string `json:"review_notes" validate:"required_if=Action REJECT"` // Wajib diisi jika menolak
}

// ──────────────────────────────────────────────────────────────────────────────
// Response DTOs
// ──────────────────────────────────────────────────────────────────────────────

// MasterDataProposalListResponse digunakan untuk halaman List/Index.
// Hanya menampilkan data header tanpa detail item (ringan & cepat).
type MasterDataProposalListResponse struct {
	ID              uuid.UUID  `json:"id"`
	ReferenceNumber string     `json:"reference_number"`
	EntityType      string     `json:"entity_type"`
	ActionType      string     `json:"action_type"`
	TotalItems      int        `json:"total_items"`
	Status          string     `json:"status"`
	ProposedByID    uuid.UUID  `json:"proposed_by_id"`
	ReviewedByID    *uuid.UUID `json:"reviewed_by_id,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	Reason          string     `json:"reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// MasterDataProposalDetailResponse digunakan untuk halaman Detail dokumen.
// Menampilkan header lengkap beserta seluruh item di dalamnya.
type MasterDataProposalDetailResponse struct {
	ID              uuid.UUID                        `json:"id"`
	ReferenceNumber string                           `json:"reference_number"`
	EntityType      string                           `json:"entity_type"`
	ActionType      string                           `json:"action_type"`
	TotalItems      int                              `json:"total_items"`
	Status          string                           `json:"status"`
	ProposedByID    uuid.UUID                        `json:"proposed_by_id"`
	ReviewedByID    *uuid.UUID                       `json:"reviewed_by_id,omitempty"`
	ReviewedAt      *time.Time                       `json:"reviewed_at,omitempty"`
	Reason          string                           `json:"reason,omitempty"`
	ReviewNotes     *string                          `json:"review_notes,omitempty"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       time.Time                        `json:"updated_at"`
	Items           []MasterDataProposalItemResponse `json:"items"`
}

// MasterDataProposalItemResponse adalah response per baris item detail.
type MasterDataProposalItemResponse struct {
	ID           uuid.UUID  `json:"id"`
	SeqNo        int        `json:"seq_no"`
	EntityID     *uuid.UUID `json:"entity_id,omitempty"`
	PayloadJSON  string     `json:"payload_json"`
	SnapshotJSON *string    `json:"snapshot_json,omitempty"`
}
