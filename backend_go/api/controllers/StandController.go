package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/api/requests"
	"project/internal/initializers"
	"project/internal/models"
	"time"
)

// @Summary Crée un nouveau stand
// @Description Crée un nouveau stand avec les informations fournies
// @Tags Stand
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param stand body requests.StandRequest true "Stand à créer"
// @Success 201 {object} models.Stand
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /create-stand [post]
func CreateStand(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 && currentUser.Role != 3 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission to do that"})
		return
	}

	var standData requests.StandRequest
	stand := models.Stand{
		Name:         standData.Name,
		Type:         standData.Type,
		JetonsRequis: standData.JetonsRequis,
		UserID:       currentUser.ID,
	}
	if err := c.ShouldBind(&stand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&stand).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"stand": stand})
}

// @Summary Récupère tous les produits
// @Description Récupère la liste de tous les stands
// @Tags Stand
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Produce json
// @Success 200 {object} []models.Stand
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /stands [get]
func GetAllStands(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission to do that"})
		return
	}
	var stands []models.Stand
	if err := initializers.DB.Find(&stands).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stands": stands})
}

// @Summary Récupère un stand par ID
// @Description Récupère les informations d'un utilisateur spécifique
// @Tags Stand
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Success 200 {object} models.Stand
// @Failure 404 {object} gin.H "Stand non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /stands/{id} [get]
func GetStandById(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	standID := c.Param("id")
	var standRetrieved models.Stand
	if err := initializers.DB.First(&standRetrieved, standID).Preload("Stock").Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stand": standRetrieved})
}

// @Summary Met à jour un stand par ID
// @Description Met à jour les informations d'un stand spécifique
// @Tags Stand
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Param stand body models.Stand true "Stand à mettre à jours"
// true "Stand à mettre à jour"
// @Success 200 {object} models.Stand
// @Failure 404 {object} gin.H "Stand non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /stands/{id}/update [put]
func UpdateStand(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	standID := c.Param("id")
	if err := initializers.DB.First("id = ?", standID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var standRetrieved models.Stand
	currentUser := user.(models.User)
	if currentUser.Role != 1 || currentUser.ID != standRetrieved.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission to do that"})
	}

	if err := initializers.DB.Save(&standRetrieved).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stand mis à jours",
		"stand":   standRetrieved})
}

// @Summary Supprime un stand par ID
// @Description Supprime un stand spécifique
// @Tags Stand
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de du stand"
// @Success 204 {object} nil
// @Failure 404 {object} gin.H "stand non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /stands/{id}/delete [delete]
func DeleteStand(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission to do that"})
		return
	}

	var stand models.Stand
	standID := c.Param("id")
	if err := initializers.DB.Delete(stand, standID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "stand not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "stand supprimé"})
}

// @Summary Interagir avec un stand
// @Description Permet d'interagir avec un stand en échange de jetons
// @Tags Stand
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "ID du stand"
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Success 200 {object} models.Stand
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /stands/{id}/interact [post]
func InteractWithStand(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	standID := c.Param("id")
	var stand models.Stand
	if err := initializers.DB.First(&stand, "id = ?", standID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Jetons <= stand.JetonsRequis {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have enough coins to do that"})
		return
	}

	stand.Conso = currentUser.Jetons - stand.JetonsRequis
	if err := initializers.DB.Save(&stand).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser.Jetons -= stand.JetonsRequis
	if err := initializers.DB.Save(&currentUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var historique = models.History{
		Date:      time.Now(),
		NbJetons:  stand.JetonsRequis,
		StandName: stand.Name,
		UserID:    currentUser.ID,
	}

	if err := initializers.DB.Create(&historique).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stand":       stand.Conso,
		"jetons user": currentUser.Jetons,
	})
}

// @Summary Achat d'un produit sur un stand
// @Description Permet à un utilisateur d'acheter un produit sur un stand avec ses jetons
// @Tags Stand
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param stand_id path uint true "ID du stand"
// @Param product_id path uint true "ID du produit"
// @Param quantity body requests.QuantityProductRequest true  "Quantité de produit à acheter"
// @Success 200 {object} gin.H "Success"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Router /stands/{stand_id}/products/{product_id}/buy [post]
func BuyProduct(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	standID := c.Param("stand_id")
	productID := c.Param("product_id")

	var stand models.Stand
	var product models.Product
	var quantity requests.QuantityProductRequest

	if err := c.BindJSON(&quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := initializers.DB.Preload("Stock").First(&stand, standID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stand not found"})
		return
	}
	if err := initializers.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if product.Nb_Products < uint64(quantity.Quantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
		return
	}

	totalJetons := product.JetonsRequis * quantity.Quantity

	if user.Jetons < totalJetons {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough tokens"})
		return
	}

	user.Jetons -= totalJetons
	stand.Conso += totalJetons
	product.Nb_Products -= uint64(quantity.Quantity)

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Save(&stand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Historiser la transaction
	historique := models.History{
		UserID:    user.ID,
		NbJetons:  totalJetons,
		StandName: stand.Name,
		Date:      time.Now(),
	}
	if err := initializers.DB.Create(&historique).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase successful"})
}

// @Summary Attribue des points à un utilisateur depuis un stand
// @Description Un stand peut attribuer des points à un utilisateur en fonction du type de stand
// @Tags Stand
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param stand_id path uint true "ID du stand"
// @Param user_id path uint true "ID de l'utilisateur"
// @Param points body requests.GivePointsRequest true "Nombre de points à attribuer"
// @Success 200 {object} gin.H "Success"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Router /stands/{stand_id}/users/{user_id}/points [post]
func GivePoints(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	standOwner := currentUser.(models.User)

	// Récupérer les paramètres du chemin
	standID := c.Param("stand_id")
	userID := c.Param("user_id")

	// Vérifier que l'utilisateur connecté est bien le propriétaire du stand
	var stand models.Stand
	if err := initializers.DB.First(&stand, standID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stand not found"})
		return
	}
	if stand.UserID != standOwner.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this stand"})
		return
	}

	var body requests.GivePointsRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.PtsAttribues += body.Points
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "points successfully given", "points_given": body.Points})
}
