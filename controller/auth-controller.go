package controller

import (
	"examples/go-crud/dto"
	"examples/go-crud/entity"
	"examples/go-crud/helper"
	"examples/go-crud/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)

	if errDTO != nil {
		response := helper.BulldErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	authResult := c.authService.VerifyCredentials(loginDTO.Email, loginDTO.Password)

	if v, ok := authResult.(entity.User); ok {
		generatedTokn := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedTokn
		response := helper.BuildResponse(true, " OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BulldErrorResponse("Please check your credentials again ", "Invalid Credentials", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO

	errDTO := ctx.ShouldBind(&registerDTO)

	if errDTO != nil {
		response := helper.BulldErrorResponse("Failed to Process the request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BulldErrorResponse("Failed to Process the request", "Duplicate Email", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusConflict, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}

}
