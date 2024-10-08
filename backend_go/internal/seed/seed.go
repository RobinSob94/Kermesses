package seed

import (
	"fmt"
	"log"
	"project/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleAdmin        = 1
	RoleOrganisateur = 2
	RoleTeneur       = 3
	RoleParent       = 4
	RoleEnfant       = 5
)

func CleanDatabase(DB *gorm.DB) {
	// Supprimer les entrées dans les tables sans générer d'erreur si elles n'existent pas
	_ = DB.Exec("DELETE FROM kermesses")
	_ = DB.Exec("DELETE FROM stands")
	_ = DB.Exec("DELETE FROM products")
	_ = DB.Exec("DELETE FROM user_parents")
	_ = DB.Exec("DELETE FROM user_enfants")
	_ = DB.Exec("DELETE FROM users")
	_ = DB.Exec("DELETE FROM jetons")

	fmt.Println("Base de données nettoyée.")
}

func SeedData(DB *gorm.DB) {
	// Nettoyer la base de données avant d'insérer les données
	CleanDatabase(DB)

	// Insertion des jetons
	jetons := []models.Jetons{
		{NbJetons: 15, Price: 10},
		{NbJetons: 32, Price: 18},
		{NbJetons: 65, Price: 25},
		{NbJetons: 80, Price: 37},
		{NbJetons: 110, Price: 50},
		{NbJetons: 150, Price: 70},
	}

	for _, jeton := range jetons {
		if err := DB.Create(&jeton).Error; err != nil {
			log.Fatalf("Erreur lors de l'insertion du jeton %v : %v", jeton, err)
		}
	}

	// Insertion des utilisateurs
	users := []models.User{
		{Firstname: "Admin", Lastname: "User", Email: "admin@example.com", Password: hashPassword("adminpass"), Role: RoleAdmin},
		{Firstname: "Organisateur1", Lastname: "User", Email: "org1@example.com", Password: hashPassword("orgpass"), Role: RoleOrganisateur},
		{Firstname: "Organisateur2", Lastname: "User", Email: "org2@example.com", Password: hashPassword("orgpass"), Role: RoleOrganisateur},
		{Firstname: "Teneur1", Lastname: "Stand", Email: "teneur1@example.com", Password: hashPassword("teneurpass"), Role: RoleTeneur},
		{Firstname: "Teneur2", Lastname: "Stand", Email: "teneur2@example.com", Password: hashPassword("teneurpass"), Role: RoleTeneur},
		{Firstname: "Parent1", Lastname: "User", Email: "parent1@example.com", Password: hashPassword("parentpass"), Role: RoleParent},
		{Firstname: "Parent2", Lastname: "User", Email: "parent2@example.com", Password: hashPassword("parentpass"), Role: RoleParent},
	}

	// Enregistrer les utilisateurs dans la base de données
	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Fatalf("Erreur lors de l'insertion de l'utilisateur %v : %v", user, err)
		}
	}

	// Création des enfants pour chaque parent
	enfants := []models.User{
		{Firstname: "Enfant1", Lastname: "Parent1", Email: "enfant1@example.com", Password: hashPassword("enfantpass"), Role: RoleEnfant},
		{Firstname: "Enfant2", Lastname: "Parent1", Email: "enfant2@example.com", Password: hashPassword("enfantpass"), Role: RoleEnfant},
		{Firstname: "Enfant1", Lastname: "Parent2", Email: "enfant3@example.com", Password: hashPassword("enfantpass"), Role: RoleEnfant},
	}

	for _, enfant := range enfants {
		if err := DB.Create(&enfant).Error; err != nil {
			log.Fatalf("Erreur lors de l'insertion de l'enfant %v : %v", enfant, err)
		}

		// Lier l'enfant au parent
		var parent models.User
		if enfant.Lastname == "Parent1" {
			DB.Where("email = ?", "parent1@example.com").First(&parent)
		} else {
			DB.Where("email = ?", "parent2@example.com").First(&parent)
		}

		// Associer l'enfant au parent
		if err := DB.Model(&parent).Association("Enfants").Append(&enfant); err != nil {
			log.Fatalf("Erreur lors de l'association de l'enfant %v avec le parent %v : %v", enfant, parent, err)
		}
	}

	// Insertion des kermesses
	kermesses := []models.Kermesse{
		{Name: "Kermesse de Printemps", Picture: "kermesse1.jpg", UserID: 1}, // Organisateur : Admin
		{Name: "Kermesse d'Été", Picture: "kermesse2.jpg", UserID: 2},        // Organisateur : Organisateur1
		{Name: "Kermesse d'Hiver", Picture: "kermesse3.jpg", UserID: 3},      // Organisateur : Organisateur2
	}

	for _, kermesse := range kermesses {
		if err := DB.Create(&kermesse).Error; err != nil {
			log.Fatalf("Erreur lors de l'insertion de la kermesse %v : %v", kermesse, err)
		}
	}

	// Insertion des stands
	stands := []models.Stand{
		{Name: "Stand de Nourriture", Type: "Nourriture", Pts_Donnees: 10, Conso: 5, JetonsRequis: 2, UserID: 4}, // Teneur1
		{Name: "Stand de Boissons", Type: "Boissons", Pts_Donnees: 5, Conso: 3, JetonsRequis: 1, UserID: 4},      // Teneur1
		{Name: "Stand de Jeux", Type: "Jeux", Pts_Donnees: 15, Conso: 8, JetonsRequis: 3, UserID: 5},             // Teneur2
	}

	for _, stand := range stands {
		if err := DB.Create(&stand).Error; err != nil {
			log.Fatalf("Erreur lors de l'insertion du stand %v : %v", stand, err)
		}
	}

	// Associer les stands aux kermesses
	for i := 0; i < len(kermesses); i++ {
		start := i * 2
		end := start + 2

		// Si nous avons moins de stands que prévu, ajuster la fin
		if end > len(stands) {
			end = len(stands)
		}

		// Si le début dépasse la longueur, nous avons fini d'associer des stands
		if start >= len(stands) {
			break
		}
		standSlice := stands[start:end] // Assurez-vous que cette partie ne dépasse pas la longueur du tableau

		if err := DB.Model(&kermesses[i]).Association("Stands").Append(standSlice); err != nil {
			log.Fatalf("Erreur lors de l'association des stands %v avec la kermesse %v : %v", standSlice, kermesses[i], err)
		}
	}

	// Insertion des produits
	products := []models.Product{
		{Name: "Frites", Picture: "frites.jpg", Type: "Nourriture", JetonsRequis: 2, Nb_Products: 100, StandID: 1},
		{Name: "Soda", Picture: "soda.jpg", Type: "Boissons", JetonsRequis: 1, Nb_Products: 200, StandID: 2},
		{Name: "Jeu de société", Picture: "jeu.jpg", Type: "Jeux", JetonsRequis: 3, Nb_Products: 50, StandID: 3},
	}

	for _, product := range products {
		if err := DB.Create(&product).Error; err != nil {
			log.Fatalf("Erreur lors de l'insertion du produit %v : %v", product, err)
		}
	}

	fmt.Println("Données insérées avec succès.")
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erreur lors du hachage du mot de passe : %v", err)
	}
	return string(bytes)
}
