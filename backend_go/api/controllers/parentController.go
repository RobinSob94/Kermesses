package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/api/requests"
	"project/internal/initializers"
	"project/internal/models"
	"strconv"
)

// @Summary Créer une relation parents/enfants
// @Description Permet de créer une relation entre les parents et les enfant
// @Tags Parent
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param parent body requests.AddChildrenRequest true "Ajouter un ou plusieurs enfants"
// @Success 200 {object} models.User
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /add-children [post]
func AddChildren(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
	}

	currentUser := user.(models.User)
	if currentUser.Role != 4 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to have the correct Role (Parent one)"})
		return
	}

	var req requests.AddChildrenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.ChildrenIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "children ids is empty"})
		return
	}

	var children []models.User
	if err := initializers.DB.Where("id IN ?", req.ChildrenIDs).Find(&children).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des enfants"})
		return
	}

	if len(children) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun enfant trouvé pour les IDs donnés"})
		return
	}

	if err := initializers.DB.Model(&currentUser).Association("Enfants").Append(children); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout des enfants"})
		return
	}
	if err := initializers.DB.Model(&children).Association("Parents").Append(currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errue lors de l'ajout de parent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Enfants ajoutés avec succès",
		"children": children,
	})
}

// @Summary Transférer des jetons aux enfants
// @Description Permet à un parent de transférer des jetons à ses enfants
// @Tags Parent
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path uint true "ID de l'enfant"
// @Param transaction body requests.GiveCoinRequest true "Détails du transfert de jetons (seulement la quantité de jetons)"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Mauvaise requête"
// @Failure 401 {object} gin.H "Non autorisé"
// @Failure 404 {object} gin.H "Enfant non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /api/users/{id}/give-coins [post]
func GiveCoins(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	currentUser := user.(models.User)

	if currentUser.Role > 4 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You do not have permission to give coins"})
		return
	}

	enfantIDParam := c.Param("id")
	enfantID, err := strconv.ParseUint(enfantIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child ID"})
		return
	}

	var req requests.GiveCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var enfant models.User
	if err := initializers.DB.Where("id = ?", enfantID).First(&enfant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	var children []models.User
	if err := initializers.DB.Model(&currentUser).Association("Enfants").Find(&children); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des enfants"})
		return
	}

	isChildLinked := false
	for _, child := range children {
		if child.ID == uint(enfantID) {
			isChildLinked = true
			break
		}
	}

	if !isChildLinked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to give coins to this child"})
		return
	}

	// Vérifier si le parent a suffisamment de jetons
	if currentUser.Jetons < req.NbJetons {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You do not have enough coins"})
		return
	}

	// Déduire les jetons du parent et les ajouter à l'enfant
	currentUser.Jetons -= req.NbJetons
	enfant.Jetons += req.NbJetons

	// Sauvegarder les modifications dans la base de données
	if err := initializers.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update parent's coins"})
		return
	}

	if err := initializers.DB.Save(&enfant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update child's coins"})
		return
	}

	// Retourner une réponse réussie
	c.JSON(http.StatusOK, gin.H{
		"message":      "Coins successfully transferred",
		"parent_coins": currentUser.Jetons,
		"child_coins":  enfant.Jetons,
	})
}
