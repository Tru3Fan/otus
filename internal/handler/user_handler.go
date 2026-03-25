package handler

import (
	"errors"
	"net/http"
	"otus/internal/model"
	"otus/internal/repository"
	"otus/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
}

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// CreateUser godoc
// @Summary Создать пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body internal_handler.UserRequest true "Данные пользователя"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.svc.CreateUser(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u})
}

// GetUsers godoc
// @Summary Получить всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Router /api/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	all, err := h.svc.GetUsers()
	if err != nil {

		return
	}
	if all == nil {
		all = []model.User{}
	}
	c.JSON(http.StatusOK, all)
}

// GetUser godoc
// @Summary Получить пользователя по ID
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/user/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	u, err := h.svc.GetUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

// UpdateUser godoc
// @Summary Обновить пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Param user body internal_handler.UserRequest true "Новые данные"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	var req UserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.svc.UpdateUser(id, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	if err := h.svc.DeleteUser(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func parseID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, err
	}
	return id, nil
}
