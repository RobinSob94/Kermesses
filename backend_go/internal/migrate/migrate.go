package migrate

import (
	"project/internal/models"

	"gorm.io/gorm"
)

// MigrateDB effectue les migrations de la base de données.
func MigrateDB(DB *gorm.DB) error {
	// Exécute les migrations pour créer les tables
	return DB.AutoMigrate(
		&models.User{},
		&models.Kermesse{},
		&models.Stand{},
		&models.Product{},
		&models.Transaction{},
		&models.Jetons{},
		&models.History{},
	)
}
