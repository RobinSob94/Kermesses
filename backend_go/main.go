package main

import (
	"log"
	"project/api/routes"
	"project/internal/initializers"
	"project/internal/migrate" // Assurez-vous que le chemin d'importation est correct
	"project/internal/seed"    // Assurez-vous que le chemin d'importation est correct

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	// Connexion à la base de données
	db := initializers.DB // Assurez-vous que `DB` est l'instance de votre base de données

	// Appeler la migration pour créer les tables
	if err := migrate.MigrateDB(db); err != nil {
		log.Fatalf("Erreur lors de la migration de la base de données : %v", err)
		return
	}

	// Appeler la seed pour insérer les données
	seed.SeedData(db) // Pas besoin de capturer une valeur de retour

	// Initialisation du serveur
	server := gin.Default()

	// Déclarer les routes
	routes.AuthRoutes(server)
	routes.UserRoutes(server)
	routes.KermesseRoutes(server)
	routes.StandRoutes(server)
	routes.ProductRoutes(server)
	routes.PaymentRoutes(server)
	routes.TransactionsRoutes(server)
	routes.JetonsRoutes(server)
	routes.ParentRoutes(server)
	routes.ElevesRoutes(server)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur : %v", err)
		return
	}
}
