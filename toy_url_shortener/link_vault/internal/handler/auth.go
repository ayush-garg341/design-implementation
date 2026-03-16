package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkvault/internal/service"
)

type UserCreateBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(s *service.UserService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var p UserCreateBody

	// Try to decode the request body into the struct
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.service.Create(c.Request.Context(), p.Name, p.Email, p.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, *u)
}

func (h *AuthHandler) Login(c *gin.Context) {

	var p UserLoginBody

	// Try to decode the request body into the struct
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), p.Email, p.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}

	c.JSON(http.StatusOK, resp)
}
