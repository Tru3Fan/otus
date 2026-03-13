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

// CreateTask godoc
// @Summary Создать задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body TaskRequest true "Данные задачи"
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/task [post]
func CreateTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := repository.AddTask(model.Task{Title: req.Title})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add task"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// GetTasks godoc
// @Summary Получить все задачи
// @Tags tasks
// @Produce json
// @Success 200 {array} model.Task
// @Router /api/tasks [get]
func GetTasks(c *gin.Context) {
	all := repository.GetAllTasks()
	if all == nil {
		all = []model.Task{}
	}
	c.JSON(http.StatusOK, all)
}

// GetTask godoc
// @Summary Получить задачу по ID
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/task/{id} [get]
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

// UpdateTask godoc
// @Summary Обновить задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Param task body TaskRequest true "Новые данные"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/task/{id} [put]
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

// DeleteTask godoc
// @Summary Удалить задачу
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/task/{id} [delete]
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
