package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"net/http"
	"os"
	"project/api/requests"
	"project/internal/initializers"
	"project/internal/models"
	"time"
)

// @Summary Crée une intention de paiement pour les jetons ou les tickets de tombola
// @Description Crée une intention de paiement pour acheter des jetons ou des tickets de tombola
// @Tags Payment
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param payment body requests.PaymentRequest true "Paiement des jetons ou tombola"
// @Success 201 {object} models.Transaction
// @Failure 500 {object} gin.H "Erreur serveur interne"
// @Router /payment [post]
func Payment(c *gin.Context) {
	// Vérification de l'utilisateur connecté
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non connecté"})
		return
	}
	currentUser := user.(models.User)

	var paymentReq requests.PaymentRequest
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Lier la requête à la structure PaymentRequest
	if err := c.ShouldBind(&paymentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Créer l'intention de paiement Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(paymentReq.Price * 100)),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gestion des paiements selon le type (jetons ou tombola)
	var transaction models.Transaction
	if paymentReq.Type == "jetons" {
		// Mettre à jour le nombre de jetons de l'utilisateur
		currentUser.Jetons += paymentReq.Quantity
		if err := initializers.DB.Save(&currentUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Créer une transaction pour les jetons
		transaction = models.Transaction{
			DateTransaction: time.Now(),
			Price:           float32(paymentReq.Price),
			Quantity:        paymentReq.Quantity,
			UserID:          currentUser.ID,
		}
	} else if paymentReq.Type == "tombola" {
		// Créer une transaction pour le ticket de tombola
		transaction = models.Transaction{
			DateTransaction: time.Now(),
			Price:           float32(paymentReq.Price),
			Quantity:        paymentReq.Quantity,
			UserID:          currentUser.ID,
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type de paiement invalide"})
		return
	}

	// Enregistrer la transaction dans la base de données
	if err := initializers.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de la création de la transaction"})
		return
	}

	// Répondre avec les détails de l'intention de paiement et de la transaction
	c.JSON(http.StatusOK, gin.H{
		"paymentIntent": pi,
		"transaction":   transaction,
	})
}
