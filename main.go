package main

import (
	"go-postgres/controllers/auth_controller"
	"go-postgres/controllers/product_controller"
	"go-postgres/controllers/user_controller"
	"go-postgres/middlewares"
	"go-postgres/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.New()
	r.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)

	models.ConnectDatabase()

	public := r
	// public
	public.POST("/register", auth_controller.Register)
	public.POST("/login", auth_controller.Login)
	public.GET("/logout", auth_controller.Logout)

	protected := r.Group("/api/v1")
	protected.Use(middlewares.JwtAuthMiddleware())

	// protocted
	protected.GET("/", product_controller.Index)
	protected.GET("/products", product_controller.FindProducts)
	protected.GET("/product/:id", product_controller.FindProductById)
	protected.POST("/product", product_controller.CreateProduct)
	protected.PUT("/product/:id", product_controller.UpdateProduct)
	protected.DELETE("/product/:id", product_controller.DeleteProduct)

	protected.GET("/users", user_controller.FindUsers)
	protected.GET("/user/:id", user_controller.FindUserById)
	protected.POST("/user", user_controller.CreateUser)
	protected.PUT("/user/:id", user_controller.UpdateUser)
	protected.DELETE("/user/:id", user_controller.DeleteUser)

	protected.GET("/user", auth_controller.CurrentUser)

	r.Run(":9898")
}
