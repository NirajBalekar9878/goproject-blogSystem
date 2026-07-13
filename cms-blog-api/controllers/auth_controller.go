package controllers

import (
	"cms-blog-api/models"
	"cms-blog-api/services"
	"cms-blog-api/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// Register handles POST /auth/register
func (ac *AuthController) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	user, err := ac.service.Register(input)
	if err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondCreated(c, "User registered successfully", user)
}

// Login handles POST /auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	tokenString, user, err := ac.service.Login(input)
	if err != nil {
		utils.RespondUnauthorized(c, err.Error())
		return
	}

	utils.RespondOK(c, "Login successful", gin.H{
		"token": tokenString,
		"user":  user,
	})
}

// GetProfile handles GET /auth/profile (protected endpoint)
func (ac *AuthController) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.RespondUnauthorized(c, "Authentication required")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		utils.RespondUnauthorized(c, "Invalid token claims")
		return
	}

	user, err := ac.service.GetUserByID(userID)
	if err != nil {
		utils.RespondNotFound(c, "User profile not found")
		return
	}

	utils.RespondOK(c, "Profile retrieved successfully", user)
}
