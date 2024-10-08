package controllers

import (
	"net/http"
	"os"
	"project/api/requests"
	"project/internal/initializers"
	"project/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Allow you to register as a new User
// @Description Create a new user with the provided information
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body requests.SignupRequest true "User data"
// @Success 201 {object} requests.SignupRequest "User created"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Bad request"
// @Failure 409 {object} gin.H "Conflict"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /signup [post]
func Signup(c *gin.Context) {
	var signupReq requests.SignupRequest

	if err := c.ShouldBindJSON(&signupReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("email = ?", signupReq.Email).Find(&userFound)
	if userFound.ID != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email already used"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Firstname: signupReq.Firstname,
		Lastname:  signupReq.Lastname,
		Password:  string(passwordHash),
		Email:     signupReq.Email,
		Picture:   signupReq.Picture,
	}
	/*mailer2.SendGoMail(user.Email, "Inscription", "./pkg/mailer/templates/registry.html", user)*/
	initializers.DB.Create(&user)
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// @Summary Allow you to log and have an JWT Token
// @Description login to the app
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body requests.LoginRequest true "User data"
// @Success 200 {object} gin.H "Connexion réussie"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Bad request"
// @Failure 409 {object} gin.H "Conflict"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginReq requests.LoginRequest

	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", loginReq.Email).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

// @Summary Logout
// @Description Inform the client to delete the token
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)//
// @Success 200 {object} gin.H "Déconnexion réussie"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Bad request"
// @Failure 409 {object} gin.H "Conflict"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /logout [post]
func Logout(c *gin.Context) {
	// Aucune action particulière nécessaire côté serveur pour les JWT
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// @Summary Récupère le profil de l'utilisateur actuellement connecté
// @Description Retourne les informations du profil de l'utilisateur connecté
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} gin.H "Success"
// @Failure 401 {object} gin.H "Unauthorized"
// @Router /profile [get]
func UserProfile(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)

	var userProfile models.User

	if err := initializers.DB.Preload("Parents").
		Preload("Enfants").
		Preload("Kermesses").
		Preload("Stands").
		Preload("Transactions").
		Preload("Historique").
		First(&userProfile, "id = ?", user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userProfile})
}

// @Summary Mise à jour du profil
// @Description Mettre à jour les champs du profil de l'utilisateur authentifié
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Insérez votre jeton d'accès" default(Bearer <Ajouter le jeton d'accès ici>)
// @Param User body requests.SignupRequest true "Les données du profil à mettre à jour"
// @Success 200 {object} gin.H "Profil mis à jour avec succès"
// @Failure 400 {object} gin.H "Erreur de validation"
// @Failure 401 {object} gin.H "Non autorisé"
// @Failure 500 {object} gin.H "Erreur du serveur"
// @Router /profile/update [put]
func UpdateProfile(c *gin.Context) {
	var signupReq requests.SignupRequest
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&signupReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := currentUser.(models.User)

	if signupReq.Firstname != "" {
		user.Firstname = signupReq.Firstname
	}
	if signupReq.Lastname != "" {
		user.Lastname = signupReq.Lastname
	}
	if signupReq.Email != "" {
		user.Email = signupReq.Email
	}
	if signupReq.Picture != "" {
		user.Picture = signupReq.Picture
	}
	if signupReq.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Password = string(passwordHash)
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}
