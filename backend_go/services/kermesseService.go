package services

import (
	"gorm.io/gorm"
	"project/internal/models"
)

type KermesseService struct {
	db *gorm.DB
}

func NewKermesseService(db *gorm.DB) *KermesseService {
	return &KermesseService{db: db}
}
func (s *KermesseService) CreateKermesse(kermesse *models.Kermesse) error {
	return s.db.Create(kermesse).Error
}

func (s *KermesseService) GetKermesseById(id uint) (*models.Kermesse, error) {
	var kermesse models.Kermesse
	if err := s.db.First(&kermesse, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &kermesse, nil
}

func (s *KermesseService) UpdateKermesse(kermesse *models.Kermesse) error {
	return s.db.Save(kermesse).Error
}

func (s *KermesseService) DeleteKermesse(id uint) error {
	return s.db.Delete(&models.Kermesse{}, "id = ?", id).Error
}
