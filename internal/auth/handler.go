package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/NotifyGo/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userService *user.UserService
}

func NewHandler(userService *user.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Login(c *gin.Context) {

	var req user.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	userFound, err := h.userService.GetByEmail(
		c.Request.Context(),
		req.Email,
	)

	if err != nil {
		fmt.Println("GET USER ERROR:", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// DEBUG
	//fmt.Printf("USER: %+v\n", userFound)

	fmt.Println("PASSWORD REQUEST:", req.PasswordHash)
	fmt.Println("PASSWORD HASH DB:", userFound.PasswordHash)

	err = bcrypt.CompareHashAndPassword(
		[]byte(userFound.PasswordHash),
		[]byte(req.PasswordHash),
	)

	// DEBUG
	fmt.Println("BCRYPT ERROR:", err)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	token, err := GenerateToken(
		userFound.ID.String(),
		string(userFound.Role),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
