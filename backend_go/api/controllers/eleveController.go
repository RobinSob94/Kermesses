package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/initializers"
	"project/internal/models"
)

// @Summary Récupère tous les utilisateurs avec le rôle d'élève
// @Description Récupère la liste de tous les utilisateurs
// @Tags Student
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Produce json
// @Success 200 {object} []models.User
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /students [get]
func GetStudents(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	var students []models.User

	if err := initializers.DB.Where("role = ?", 5).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des élèves"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"students": students})
}
