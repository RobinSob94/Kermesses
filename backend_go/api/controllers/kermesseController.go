package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/api/requests"
	"project/internal/initializers"
	"project/internal/models"
)

// @Summary Créé une kermesse
// @Description Permet aux user de créé un groupe de groupeVoyage
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param kermesse body requests.KermeseRequest true "Données de la kermesse"
// @Success 201 {object} gin.H "Groupe créé"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Voyage non trouvé"
// @Failure 409 {object} gin.H "Conflict"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /create-kermesse [post]
func CreateKermesse(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	currentUser := user.(models.User)
	if currentUser.Role != 1 && currentUser.Role != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You do not have permission to do that"})
		return
	}

	var kermesseData requests.KermeseRequest
	kermesse := models.Kermesse{
		Name:   kermesseData.Name,
		UserID: currentUser.ID,
	}

	if err := c.ShouldBind(&kermesse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Create(&kermesse).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Kermesse créée avec succès",
		"kermesse": kermesse,
	})

}

// @Summary Get all Kermesses based on user role
// @Description Fetches kermesses for admin, organizers, and participants
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Success 200 {object} []models.Kermesse "List of kermesses"
// @Failure 401 {object} gin.H "User not logged"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /kermesses [get]
func GetAllKermesses(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	var kermesses []models.Kermesse
	currentUser := user.(models.User)

	if currentUser.Role == 1 {
		if err := initializers.DB.Preload("Organisateurs").Preload("Participants").Preload("Stands").Find(&kermesses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"kermesses": kermesses})
		return
	}

	if currentUser.Role == 2 {
		if err := initializers.DB.Preload("Organisateurs").Preload("Participants").Preload("Stands").
			Where("user_id = ?", currentUser.ID).Or("id IN (SELECT kermesse_id FROM kermesse_organisateurs WHERE user_id = ?)", currentUser.ID).
			Find(&kermesses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"kermesses": kermesses})
		return
	}

	if currentUser.Role >= 3 {
		if err := initializers.DB.
			Where("id IN (SELECT kermesse_id FROM kermesse_participants WHERE user_id = ?)", currentUser.ID).
			Find(&kermesses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(kermesses) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Aucune kermesse trouvée pour cet utilisateur"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"kermesses": kermesses})
}

// @Summary Get a Kermesse by its ID
// @Description Fetch a kermesse by its ID if the user has permission
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "Kermesse ID"
// @Success 200 {object} models.Kermesse "Kermesse found"
// @Failure 401 {object} gin.H "User not logged"
// @Failure 403 {object} gin.H "Forbidden"
// @Failure 404 {object} gin.H "Kermesse not found"
// @Router /kermesses/{id} [get]
func GetKermesseById(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	var kermesse models.Kermesse
	id := c.Param("id")

	if err := initializers.DB.Preload("Organisateurs").Preload("Participants").Preload("Stands").
		Where("id = ?", id).First(&kermesse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse not found"})
		return
	}

	currentUser := user.(models.User)

	// Vérifier les permissions
	if currentUser.Role == 1 || // Admin
		kermesse.UserID == currentUser.ID || // Créateur de la kermesse
		initializers.DB.Model(&kermesse).Where("id IN (SELECT kermesse_id FROM kermesse_organisateurs WHERE user_id = ?)", currentUser.ID).RowsAffected > 0 || // Organisateur
		initializers.DB.Model(&kermesse).Where("id IN (SELECT kermesse_id FROM kermesse_participants WHERE user_id = ?)", currentUser.ID).RowsAffected > 0 { // Participant
		c.JSON(http.StatusOK, gin.H{"kermesse": kermesse})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this kermesse"})
	}
}

// @Summary Ajouter des stands à la kermesse
// @Description Permet d'ajouter des users (partcicpants ou organisateurs)
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "Kermesse ID"
// @Param kermesse body requests.AddStandRequest true "Données du groupe"
// @Success 200 {object} gin.H "Stand ajouté"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "kermesse non trouvé"
// @Failure 409 {object} gin.H "Conflict"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /kermesses/{id}/add-stands [post]
func AddStand(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	kermesseID := c.Param("id")
	var kermesse models.Kermesse

	if err := initializers.DB.Preload("Organisateurs").First(&kermesse, "id = ?", kermesseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse not found"})
		return
	}

	currentUser := user.(models.User)

	isOrganisateur := false
	for _, organisateur := range kermesse.Organisateurs {
		if organisateur.ID == currentUser.ID {
			isOrganisateur = true
			break
		}
	}

	if currentUser.Role != 1 && currentUser.ID != kermesse.UserID && !isOrganisateur {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have the permission to do this"})
		return
	}

	var standReq requests.AddStandRequest
	if err := c.ShouldBindJSON(&standReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if len(standReq.StandIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "StandId is required"})
		return
	}

	var stands []models.Stand
	if err := initializers.DB.Where("id IN ?", standReq.StandIds).Find(&stands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des enfants"})
		return
	}

	if len(stands) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun enfant trouvé pour les IDs donnés"})
		return
	}

	if err := initializers.DB.Model(&kermesse).Association("Stands").Append(&stands); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding stand to kermesse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stand added successfully",
		"stand":   stands,
	})
}

// @Summary Ajouter des users à la kermesse
// @Description Permet aux user de créé un groupe de groupeVoyage
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "Kermesse ID"
// @Param kermesse body requests.AddUserRequest true "Données du groupe"
// @Success 200 {object} gin.H "User ajouté(s)"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "kermesse non trouvé"
// @Failure 409 {object} gin.H "Conflict"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /kermesses/{id}/add-users [post]
func AddParticipantAndOrga(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	kermesseID := c.Param("id")
	var kermesse models.Kermesse

	if err := initializers.DB.Preload("Organisateurs").First(&kermesse, "id = ?", kermesseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse not found"})
		return
	}

	var users []models.User
	currentUser := user.(models.User)

	isOrganisateur := false
	for _, organisateur := range kermesse.Organisateurs {
		if organisateur.ID == currentUser.ID {
			isOrganisateur = true
			break
		}
	}

	if currentUser.Role != 1 && currentUser.ID != kermesse.UserID && !isOrganisateur {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have the permission to do this"})
		return
	}

	var userReq requests.AddUserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if userReq.Type == "participants" {
		if len(userReq.UserIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserIds are required"})
			return
		}

		if err := initializers.DB.Where("id IN ?", userReq.UserIds).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des enfants"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun enfant trouvé pour les IDs donnés"})
			return
		}

		if err := initializers.DB.Model(&kermesse).Association("Participants").Append(&users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding a participant to kermesse"})
			return
		}
	} else if userReq.Type == "organisateurs" {
		if len(userReq.UserIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserIds are required"})
			return
		}

		if err := initializers.DB.Where("id IN ?", userReq.UserIds).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des enfants"})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun enfant trouvé pour les IDs donnés"})
			return
		}

		if err := initializers.DB.Model(&kermesse).Association("Participants").Append(&users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding a participant to kermesse"})
			return
		}

		if err := initializers.DB.Save(&kermesse).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding a participant to kermesse"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users added successfully"})
}

// @Summary Update a Kermesse
// @Description Allows an admin or the creator to update a Kermesse
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "Kermesse ID"
// @Param kermesse body models.Kermesse true "Kermesse data"
// @Success 200 {object} models.Kermesse "Kermesse updated"
// @Failure 401 {object} gin.H "User not logged"
// @Failure 403 {object} gin.H "Forbidden"
// @Failure 404 {object} gin.H "Kermesse not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /kermesses/{id}/update [put]
func UpdateKermesse(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	currentUser := user.(models.User)
	kermesseID := c.Param("id")

	// Rechercher la kermesse à mettre à jour
	var kermesse models.Kermesse
	if err := initializers.DB.Where("id = ?", kermesseID).First(&kermesse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse not found"})
		return
	}

	if currentUser.Role != 1 && kermesse.UserID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this kermesse"})
		return
	}

	if err := c.ShouldBindJSON(&kermesse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Save(&kermesse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update kermesse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"kermesse": kermesse})
}

// @Summary Delete a Kermesse
// @Description Allows an admin or the creator to update a Kermesse
// @Tags Kermesse
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer Add access token here)
// @Param id path int true "Kermesse ID"
// @Success 200 {object} gin.H "Kermesse delete"
// @Failure 401 {object} gin.H "User not logged"
// @Failure 403 {object} gin.H "Forbidden"
// @Failure 404 {object} gin.H "Kermesse not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /kermesses/{id}/delete [delete]
func DeleteKermesse(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged"})
		return
	}

	currentUser := user.(models.User)
	kermesseID := c.Param("id")

	var kermesse models.Kermesse
	if err := initializers.DB.First(&kermesse, "id = ?", kermesseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kermesse not found"})
		return
	}

	if currentUser.Role != 1 && kermesse.UserID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this kermesse"})
		return
	}

	if err := initializers.DB.Delete(&kermesse, kermesseID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete kermesse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kermesse supprimé avec succès"})
}
