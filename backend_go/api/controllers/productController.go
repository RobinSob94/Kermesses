package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/initializers"
	"project/internal/models"
)

// @Summary Crée un nouveau produit
// @Description Crée un nouvel produit avec les informations fournies
// @Tags Product
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param product body models.Product true "Produit à créer"
// @Success 201 {object} models.Product
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /create-product  [post]
func CreateProduct(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to create a product"})
		return
	}

	var product models.Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "produit créé avec succès",
		"product": product})
}

// @Summary Récupère tous les produits
// @Description Récupère la liste de tous les produits
// @Tags Product
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Produce json
// @Success 200 {object} []models.Product
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /products [get]
func GetProducts(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
	}

	var products []models.Product
	if err := initializers.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// @Summary Met à jour un produit par son ID
// @Description Met à jour les informations d'un produit spécifique
// @Tags Product
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de l'utilisateur"
// @Param product body models.Product true "Produit à mettre à jour"
// @Success 200 {object} models.Product
// @Failure 404 {object} gin.H "Prodduit non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /products/{id}/update [put]
func UpdateProduct(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not logged"})
		return
	}
	currentUser := user.(models.User)

	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	productID := c.Param("id")
	if err := initializers.DB.First("id = ?", productID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := initializers.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated",
		"product": product})
}

// @Summary Supprime un produit par ID
// @Description Supprime un utilisateur spécifique
// @Tags Product
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "ID de du produit"
// @Success 204 {object} nil
// @Failure 404 {object} gin.H "produit non trouvé"
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /products/{id}/delete [delete]
func DeleteProduct(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user"})
		return
	}
	id := c.Param("id")
	var productFound models.Product
	if err := initializers.DB.First(&productFound, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produdct not found"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	if err := initializers.DB.Delete(&productFound, id).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
