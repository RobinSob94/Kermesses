package models

type Jetons struct {
	ID       uint    `gorm:"primary_key; auto_increment" json:"id"`
	NbJetons uint    `json:"nb_jetons" json:"nb_jetons"`
	Price    float64 `json:"price" json:"price"`
}
