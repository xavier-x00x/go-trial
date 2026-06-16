package row

import (
	"encoding/json"
	"time"
)

type MasterDataProposalListRow struct {
	ID              string     `json:"id"`
	ReferenceNumber string     `json:"reference_number"`
	EntityType      string     `json:"entity_type"`
	ActionType      string     `json:"action_type"`
	TotalItems      int        `json:"total_items"`
	Status          string     `json:"status"`
	ProposedByID    string     `json:"proposed_by_id"`
	ProposedByName  string     `json:"proposed_by_name,omitempty"`
	ReviewedByID    *string    `json:"reviewed_by_id,omitempty"`
	ReviewedByName  *string    `json:"reviewed_by_name,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	Reason          string     `json:"reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type MasterDataProposalDetailRow struct {
	MasterDataProposalListRow
	ReviewNotes *string                     `json:"review_notes,omitempty"`
	UpdatedAt   time.Time                   `json:"updated_at"`
	Items       []MasterDataProposalItemRow `json:"items" gorm:"-"`
}

type MasterDataProposalItemRow struct {
	ID           string          `json:"id"`
	ProposalID   string          `json:"proposal_id,omitempty"`
	SeqNo        int             `json:"seq_no"`
	EntityID     *string         `json:"entity_id,omitempty"`
	PayloadJSON  json.RawMessage `json:"payload_json"`
	SnapshotJSON json.RawMessage `json:"snapshot_json,omitempty"`
}
