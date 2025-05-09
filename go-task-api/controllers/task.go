package controllers

import (
	"fmt"
	"go-task-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
    var input struct {
        Title string `json:"title"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, _ := c.Get("user")
    task := models.Task{Title: input.Title, UserID: user.(models.User).ID}

    if err := models.DB.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"task": task})
}

func ListTasks(c *gin.Context) {
    user, _ := c.Get("user")

    pageStr := c.DefaultQuery("page", "1")
    sizeStr := c.DefaultQuery("size", "10")

    page, _ := strconv.Atoi(pageStr)
    size, _ := strconv.Atoi(sizeStr)
    offset := (page - 1) * size

		start := time.Now()

		var tasks []models.TaskPreview
		var total int64

		models.DB.Model(&models.Task{}).
        Where("user_id = ? AND deleted_at IS NULL", user.(models.User).ID).
        Count(&total)
				
    if err := models.DB.Model(&models.Task{}).
        Where("user_id = ?", user.(models.User).ID).
				Order("id ASC").
        Limit(size).
        Offset(offset).
        Find(&tasks).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks"})
        return
    }

		fmt.Println("查詢時間：", time.Since(start))

    c.JSON(http.StatusOK, gin.H{
        "page":  page,
        "size":  size,
        "tasks": tasks,
        "total": total,
    })
}

func GetTask(c *gin.Context) {
	user, _ := c.Get("user")
	var task models.Task

	if err := models.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.(models.User).ID).
			First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// PUT /api/tasks/:id
func UpdateTask(c *gin.Context) {
	user, _ := c.Get("user")
	var task models.Task

	if err := models.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.(models.User).ID).First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
	}

	var input struct {
			Title     *string `json:"title"`
			Completed *bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	if input.Title != nil {
			task.Title = *input.Title
	}
	if input.Completed != nil {
			task.Completed = *input.Completed
	}

	models.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// DELETE /api/tasks/:id
func DeleteTask(c *gin.Context) {
	user, _ := c.Get("user")
	var task models.Task

	if err := models.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.(models.User).ID).First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
	}

	models.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func ViewTrashedTasks(c *gin.Context) {
	user, _ := c.Get("user")
	var tasks []models.Task

	models.DB.
			Unscoped().
			Where("user_id = ? AND deleted_at IS NOT NULL", user.(models.User).ID).
			Order("id ASC").
			Find(&tasks)

	c.JSON(http.StatusOK, gin.H{"trashed_tasks": tasks})
}

func RestoreTask(c *gin.Context) {
	user, _ := c.Get("user")
	var task models.Task

	if err := models.DB.
		Unscoped().
		Where("id = ? AND user_id = ?", c.Param("id"), user.(models.User).ID).
		First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := models.DB.Unscoped().
		Model(&task).
		Updates(map[string]interface{}{"deleted_at": nil}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task restored", "task": task})
}

func HardDeleteTask(c *gin.Context) {
	user, _ := c.Get("user")
	var task models.Task

	if err := models.DB.
		Unscoped().
		Where("id = ? AND user_id = ?", c.Param("id"), user.(models.User).ID).
		First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	models.DB.Unscoped().Delete(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task permanently deleted"})
}