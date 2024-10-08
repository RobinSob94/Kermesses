package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/initializers"
	"project/internal/models"
)

// @Summary Crée un nouveau jeton
// @Description Crée un nouveau jeton
// @Tags Jeton
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param user body models.Jetons true "Jeton à créer"
// @Success 201 {object} models.Jetons
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /create-jeton  [post]
func CreateJetons(c *gin.Context) {
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
	var jetons models.Jetons
	if err := c.BindJSON(&jetons); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Create(&jetons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": jetons})
}

// @Summary Récupère tous les jetons
// @Description Récupère la liste de tous les jetons
// @Tags Jeton
// @Produce json
// @Success 200 {object} []models.Jetons
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /jetons [get]
func GetJetons(c *gin.Context) {
	/*	_, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}*/

	var jetons []models.Jetons
	if err := initializers.DB.Find(&jetons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jetons": jetons})
}

// @Summary Met à jour un jeton par ID
// @Description Met à jour les informations d'un jeton spécifique
// @Tags Jeton
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Param user body models.Jetons true "Utilisateur à mettre à jour"
// @Success 200 {object} models.Jetons
// @Failure 404 {object} gin.H "Jeton non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /jetons/{id}/update [put]
func UpdateJeton(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jetonID := c.Param("id")
	var jeton models.Jetons
	if err := initializers.DB.First(&jeton, jetonID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}
	var updatedJeton models.Jetons
	if err := c.BindJSON(&updatedJeton); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Save(&updatedJeton).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedJeton})
}

// @Summary Supprime un jeton par ID
// @Description Supprime un jeton spécifique
// @Tags Jeton
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID du jeton"
// @Success 204 {object} nil
// @Failure 404 {object} gin.H "Jeton non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /jetons/{id}/delete [delete]
func DeleteJeton(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jetonID := c.Param("id")
	var jeton models.Jetons
	if err := initializers.DB.First(&jeton, "id = ?", jetonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jeton not found"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	if err := initializers.DB.Delete(&jeton).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du jeton"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jeton supprimé avec succès"})
}
