package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Авторизация
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login and password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Login != os.Getenv("LOGIN") || req.Password != os.Getenv("PASSWORD") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return

	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": req.Login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
