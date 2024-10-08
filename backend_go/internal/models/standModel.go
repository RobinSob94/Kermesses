package models

type Stand struct {
	ID           uint      `gorm:"primary_key; not null; autoIncrement " json:"id"`
	Name         string    `gorm:"size:64; not null" json:"name"`
	Type         string    `gorm:"size:64; not null" json:"type"`
	Stock        []Product `gorm:"foreignKey:StandID" json:"stocks"` // Clé étrangère vers Product
	Pts_Donnees  uint      `gorm:"not null" json:"pts_donnees"`
	Conso        uint      `gorm:"not null" json:"conso"`
	JetonsRequis uint      `gorm:"not null" json:"jetons_requis"`

	Kermesses []Kermesse `gorm:"many2many:kermesse_stands;" json:"kermesses"`

	UserID uint `gorm:"not null" json:"user_id"`
}
