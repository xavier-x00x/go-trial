package dto

// UOM (Unit of Measure) DTOs
type CreateUOMRequest struct {
	Code string `json:"code" validate:"required,min=1,max=10"`
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type UpdateUOMRequest struct {
	Code *string `json:"code,omitempty" validate:"omitempty,min=1,max=10"`
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
}

type UOMResponse struct {
	ID        string `json:"id"`
	Code     string `json:"code"`
	Name    string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}