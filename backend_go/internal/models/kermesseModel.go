package models

type Kermesse struct {
	ID      uint   `gorm:"primary_key; not null; autoIncrement" json:"id"`
	Name    string `gorm:"size:64; not null" json:"name"`
	Picture string `gorm:"size:64" json:"picture"`

	Stands []Stand `gorm:"many2many:kermesse_stands;" json:"stands"`

	// Relations Many-to-Many : Organisateurs et participants de la kermesse
	Organisateurs []User `gorm:"many2many:kermesse_organisateurs;" json:"organisateurs"`
	Participants  []User `gorm:"many2many:kermesse_participants;" json:"participants"`

	// Relation Many-to-One : L'utilisateur qui crée la kermesse
	UserID uint `gorm:"not null" json:"user_id"`
}
