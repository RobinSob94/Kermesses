package routes

import (
	"github.com/gin-gonic/gin"
	"project/api/controllers"
	"project/api/middlewares"
)

// Authentifications
func AuthRoutes(r *gin.Engine) {
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.GET("/profile", middlewares.CheckAuth, controllers.UserProfile)
	r.PUT("/profile/update", middlewares.CheckAuth, controllers.UpdateProfile)
}

func UserRoutes(r *gin.Engine) {
	r.POST("/api/users", middlewares.CheckAuth, controllers.CreateUser)
	r.GET("/api/users", middlewares.CheckAuth, controllers.GetAllUsers)
	r.GET("/api/users/:id", middlewares.CheckAuth, controllers.GetUser)
	r.PUT("/api/users/:id", middlewares.CheckAuth, controllers.UpdateUser)
	r.DELETE("/api/users/:id", middlewares.CheckAuth, controllers.DeleteUser)
}

func KermesseRoutes(r *gin.Engine) {
	r.POST("/create-kermesse", middlewares.CheckAuth, controllers.CreateKermesse)
	r.GET("/kermesses", middlewares.CheckAuth, controllers.GetAllKermesses)
	r.GET("/kermesses/:id", middlewares.CheckAuth, controllers.GetKermesseById)
	r.PUT("/kermesses/:id/update", middlewares.CheckAuth, controllers.UpdateKermesse)
	r.DELETE("/kermesses/:id/delete", middlewares.CheckAuth, controllers.DeleteKermesse)
	r.POST("/kermesses/:id/add-stands", middlewares.CheckAuth, controllers.AddStand)
	r.POST("/kermesses/:id/add-users", middlewares.CheckAuth, controllers.AddParticipantAndOrga)
}

func StandRoutes(r *gin.Engine) {
	r.POST("/create-stand", middlewares.CheckAuth, controllers.CreateStand)
	r.POST("/stands/:id/interact", middlewares.CheckAuth, controllers.InteractWithStand)
	r.GET("/stands", middlewares.CheckAuth, controllers.GetAllStands)
	r.GET("/stands/:id", middlewares.CheckAuth, controllers.GetStandById)
	r.PUT("/stands/:id/update", middlewares.CheckAuth, controllers.UpdateStand)
	r.DELETE("/stands/:id/delete", middlewares.CheckAuth, controllers.DeleteStand)
	r.POST("/stands/:id/products/products/:product_id/buy", middlewares.CheckAuth, controllers.BuyProduct)
	r.POST("/stands/:id/users/:user_id/points", middlewares.CheckAuth, controllers.GivePoints)
}

func ProductRoutes(r *gin.Engine) {
	r.POST("/create-product", middlewares.CheckAuth, controllers.CreateProduct)
	r.GET("/products", middlewares.CheckAuth, controllers.GetProducts)
	r.PUT("/products/:id/update", middlewares.CheckAuth, controllers.UpdateProduct)
	r.DELETE("/products/:id/delete", middlewares.CheckAuth, controllers.DeleteProduct)
}

func JetonsRoutes(r *gin.Engine) {
	r.POST("/create-jeton", middlewares.CheckAuth, controllers.CreateJetons)
	r.GET("/jetons", controllers.GetJetons)
	r.PUT("/jetons/:id/update", middlewares.CheckAuth, controllers.UpdateJeton)
	r.DELETE("/jetons/:id/delete", middlewares.CheckAuth, controllers.DeleteJeton)
}

func PaymentRoutes(r *gin.Engine) {
	r.POST("/payment", middlewares.CheckAuth, controllers.Payment)
}

func TransactionsRoutes(r *gin.Engine) {
	r.GET("/transactions", middlewares.CheckAuth, controllers.GetTransactions)
}

func ParentRoutes(r *gin.Engine) {
	r.POST("/add-children", middlewares.CheckAuth, controllers.AddChildren)
	r.POST("/api/users/:id/give-coins", middlewares.CheckAuth, controllers.GiveCoins)
}

func ElevesRoutes(r *gin.Engine) {
	r.GET("/students", middlewares.CheckAuth, controllers.GetStudents)
}
