package entity

type Permission struct {
	BaseModel
	Path string `gorm:"unique" json:"path"`
	Name string `json:"name"`
}
