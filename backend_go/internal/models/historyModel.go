package models

import "time"

type History struct {
	ID        uint      `gorm:"primary_key; autoIncrement" json:"id"`
	Date      time.Time `gorm:"not null" json:"date"`
	NbJetons  uint      `gorm:"not null" json:"nb_jetons"`
	StandName string    `gorm:"not null" json:"stand_name"`
	UserID    uint      `gorm:"not null" json:"user_id"`
}
