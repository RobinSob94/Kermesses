package services

import (
	"gorm.io/gorm"
	"project/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Créatiion d'un nouveau User
func (s *UserService) CreateUser(user models.User) error {
	return s.db.Create(user).Error
}

// Récupérer un user par ID
func (s *UserService) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Mise à jour d'un user
func (s *UserService) UpdateUser(user models.User) error {
	return s.db.Save(&user).Error
}

// Suppression d'un utilisateur
func (s *UserService) DeleteUserById(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}
