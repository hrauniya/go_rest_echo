package controller

import (
	"net/http"

	"go_echo_rest/config"
	"go_echo_rest/dto"
	"go_echo_rest/middleware"
	"go_echo_rest/model"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	req := new(dto.RegisterDTO)
	db := config.DB()

	if err := c.Bind(req); err != nil {
		return FailResponse(c, http.StatusBadRequest, "Invalid Request body")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return FailResponse(c, http.StatusBadRequest, "Username, email, and password are required")
	}

	var existingUser model.User
	if result := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser); result.RowsAffected > 0 {
		return FailResponse(c, http.StatusConflict, "Username or email already exists")
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := db.Create(&user).Error; err != nil {
		return FailResponse(c, http.StatusInternalServerError, "Failed to register user")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user":    user.ToUserResponse(),
	})
}

func Login(c echo.Context) error {
	req := new(dto.LoginDTO)
	db := config.DB()

	if err := c.Bind(req); err != nil {
		return FailResponse(c, http.StatusBadRequest, "Invalid Request body")
	}

	var user model.User
	if result := db.Where("username = ?", req.Username).First(&user); result.Error != nil {
		return FailResponse(c, http.StatusUnauthorized, "Invalid Credentials")
	}

	if !user.CheckPassword(req.Password) {
		return FailResponse(c, http.StatusUnauthorized, "Incorrect Password. Please try again.")
	}

	accessToken, err := middleware.GenerateAccessToken(user)
	if err != nil {
		return FailResponse(c, http.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := middleware.GenerateRefreshToken(user)
	if err != nil {
		return FailResponse(c, http.StatusInternalServerError, "Failed to generate refresh token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Login successful",
		"user":         user.ToUserResponse(),
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func RefreshToken(c echo.Context) error {
	
	req := new(dto.RefreshDTO)

	if err := c.Bind(req); err != nil {
		return FailResponse(c, http.StatusBadRequest, "Invalid request body or missing refresh_token")
	}

	if req.RefreshToken == "" {
		return FailResponse(c, http.StatusBadRequest, "Refresh token is required")
	}

	refreshClaims, err := middleware.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return FailResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token")
	}

	var user model.User
	db := config.DB()
	if result := db.First(&user, refreshClaims.UserID); result.Error != nil {
		return FailResponse(c, http.StatusUnauthorized, "User associated with refresh token not found")
	}
	newAccessToken, err := middleware.GenerateAccessToken(user)
	if err != nil {
		return FailResponse(c, http.StatusInternalServerError, "Failed to generate new access token")
	}

	newRefreshToken, err := middleware.GenerateRefreshToken(user)
	if err != nil {
		return FailResponse(c, http.StatusInternalServerError, "Failed to generate new refresh token")
	}
	
	return SuccessResponse(c, map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})

}

func GetCurrentUser(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	db := config.DB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return FailResponse(c, http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user.ToUserResponse(),
	})


}

