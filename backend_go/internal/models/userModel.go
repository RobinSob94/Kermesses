package models

type User struct {
	ID           uint   `gorm:"primary_key; not null; autoIncrement" json:"id"`
	Firstname    string `gorm:"size:64; not null" json:"firstname"`
	Lastname     string `gorm:"size:64; not null" json:"lastname"`
	Email        string `gorm:"size:100; not null; unique" json:"email"`
	Password     string `gorm:"size:100; not null" json:"password"`
	Picture      string `gorm:"size:100;" json:"picture"`
	Role         uint   `gorm:"size: 64; not null" json:"role"` /* 1 = ADMIN / 2 = ORGANISATEUR / 3 = TENEUR DE STAND / 4 = PARENT / 5 ELEVE  */
	Jetons       uint   `gorm:"size: 64; default:0; not null" json:"jetons"`
	PtsAttribues uint   `gorm:"size: 64; default: 0" json:"pts_attribues"`

	// Relations Many-to-Many pour Parents/Enfants
	Parents []User `gorm:"many2many:user_parents;" json:"parents"`
	Enfants []User `gorm:"many2many:user_enfants;" json:"enfants"`

	// Relations One-to-Many pour les Kermesses et Stands gérés par l'utilisateur
	Kermesses    []Kermesse    `gorm:"foreignKey:UserID" json:"kermesses"`
	Stands       []Stand       `gorm:"foreignKey:UserID" json:"stands"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions"`
	Historique   []History     `gorm:"foreignKey:UserID" json:"historique"`
}
