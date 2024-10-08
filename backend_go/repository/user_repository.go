package repository

import "project/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindById(id int64) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
}
