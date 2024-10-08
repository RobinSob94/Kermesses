package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/initializers"
	"project/internal/models"
)

// @Summary Récupère toutes les transactions faites par le user
// @Description Récupère la liste de toutes les transactions
// @Tags Transactions
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Produce json
// @Success 200 {object} []models.Transaction
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /transactions [get]
func GetTransactions(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not loggeg"})
		return
	}

	currentUser := user.(models.User)

	var transactions []models.Transaction
	if err := initializers.DB.Find(&transactions, currentUser.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
