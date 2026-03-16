package handler

import (
	"errors"
	"net/http"
	"otus/internal/model"
	"otus/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserResult struct {
	Username string `json:"username" binding:"required"`
}

// POST /api/user
func CreateUser(c *gin.Context) {
	var req UserResult
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := model.User{Username: req.Username}
	if err := repository.AddUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sace user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u})
}

// GET /api/users
func GetUsers(c *gin.Context) {
	all, err := repository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if all == nil {
		all = []model.User{}
	}
	c.JSON(http.StatusOK, all)
}

// GET /api/user/:id
func GetUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	u, err := repository.GetUserByID(id)
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

// PUT /api/user/:id
func UpdateUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	var req UserResult
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := repository.UpdateUser(id, model.User{Username: req.Username})
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

// DELETE /api/user/:id
func DeleteUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	if err := repository.DeleteUser(id); err != nil {
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
