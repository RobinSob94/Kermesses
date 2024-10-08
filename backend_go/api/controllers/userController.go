package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/api/requests"
	"project/internal/initializers"
	"project/internal/models"
)

// @Summary Crée un nouvel utilisateur
// @Description Crée un nouvel utilisateur avec les informations fournies
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param user body requests.SignupRequest true "Utilisateur à créer"
// @Success 201 {object} models.User
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /api/users  [post]
func CreateUser(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := user.(models.User)
	//Si le user n'est pas admin
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	var createdUser requests.SignupRequest
	newUser := models.User{
		Firstname: createdUser.Firstname,
		Lastname:  createdUser.Lastname,
		Email:     createdUser.Email,
		Password:  createdUser.Password,
		Picture:   createdUser.Picture,
		Role:      createdUser.Role,
	}
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": newUser})
}

// GetUsers - Récupère tous les utilisateurs
// @Summary Récupère tous les utilisateurs
// @Description Récupère la liste de tous les utilisateurs
// @Tags User
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Produce json
// @Success 200 {object} []models.User
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /api/users [get]
func GetAllUsers(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
	}

	var userRetrieved []models.User
	if err := initializers.DB.Find(&userRetrieved).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userRetrieved)
}

// @Summary Récupère un utilisateur par ID
// @Description Récupère les informations d'un utilisateur spécifique
// @Tags User
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {object} models.User
// @Failure 404 {object} gin.H "Utilisateur non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /api/users/{id} [get]
func GetUser(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	userID := c.Param("id")
	var userRetrieved models.User
	if err := initializers.DB.First(&userRetrieved, "id = ?", userID).
		Preload("Parents").
		Preload("Enfants").
		Preload("Kermesses").
		Preload("Stands").
		Preload("Transactions").
		Preload("Historique").
		First(&userRetrieved, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userRetrieved)
}

// @Summary Met à jour un utilisateur par ID
// @Description Met à jour les informations d'un utilisateur spécifique
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Param user body models.User true "Utilisateur à mettre à jour"
// @Success 200 {object} models.User
// @Failure 404 {object} gin.H "Utilisateur non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	userID := c.Param("id")
	var userRetrieved models.User
	if err := initializers.DB.First(&userRetrieved, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Save(&updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// @Summary Supprime un utilisateur par ID
// @Description Supprime un utilisateur spécifique
// @Tags User
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Success 204 {object} nil
// @Failure 404 {object} gin.H "Utilisateur non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := c.Param("id")
	var userRetrieved models.User
	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	if err := initializers.DB.First(&userRetrieved, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find this user"})
		return
	}

	if err := initializers.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": "Utilisateur supprimé"})
}
