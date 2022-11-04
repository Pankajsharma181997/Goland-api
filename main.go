package main

import (
	"examples/go-crud/controller"
	"examples/go-crud/initializers"
	"examples/go-crud/repository"
	"examples/go-crud/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

var (
	db             *gorm.DB                  = initializers.ConnectToDB()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService     services.JWTService       = services.NewJWTService()
	authService    services.AuthService      = services.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
)

func main() {
	defer initializers.CloseDBConnection(db)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}
	r.Run(":3000") // listen and serve on 0.0.0.0:8080
}
