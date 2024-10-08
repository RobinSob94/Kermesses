package models

type Tombola struct {
	ID    uint    `gorm:"primary_key; autoIncreme,t; not null" json:"id"`
	Price float32 `gorm:"not null" json:"price"`
}
