package services

import (
	"gorm.io/gorm"
	"project/internal/models"
)

type StandService struct {
	db *gorm.DB
}

func NewStandService(db *gorm.DB) *StandService {
	return &StandService{db: db}
}

func (s *StandService) CreateStand(stand *models.Stand) error {
	return s.db.Create(stand).Error
}

func (s *StandService) GetStandByID(id uint) (*models.Stand, error) {
	var stand models.Stand
	if err := s.db.Where("id = ?", id).First(&stand).Error; err != nil {
		return nil, err
	}
	return &stand, nil
}

func (s *StandService) UpdateStand(stand *models.Stand) error {
	return s.db.Save(stand).Error
}

func (s *StandService) DeleteStand(id uint) error {
	return s.db.Delete(&models.Stand{}, id).Error
}
