package handler

import (
	"errors"
	"net/http"
	"otus/internal/model"
	"otus/internal/repository"

	"github.com/gin-gonic/gin"
)

type TaskRequest struct {
	Title string `json:"title" binding:"required"`
}

// POST /api/task
func CreateTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t := model.Task{Title: req.Title}
	if err := repository.AddTask(t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add task"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// GET /api/tasks
func GetTasks(c *gin.Context) {
	all := repository.GetAllTasks()
	if all == nil {
		all = []model.Task{}
	}
	c.JSON(http.StatusOK, all)
}

// GET /api/task/:id
func GetTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	t, err := repository.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// PUT /api/task/:id
func UpdateTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := repository.UpdateTask(id, model.Task{Title: req.Title})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DELETE /api/task/:id
func DeleteTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}

	if err := repository.DeleteTask(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
