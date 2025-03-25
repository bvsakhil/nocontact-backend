// controllers/avoid_controller.go
package controllers

import (
	"avoids-backend/database"
	"avoids-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateAvoid(c *gin.Context) {
	var avoid models.Avoid
	if err := c.BindJSON(&avoid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	avoid.UserID = userID.(uint)

	avoid.StartDate = time.Now()
	avoid.IsActive = true

	result := database.DB.Create(&avoid)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, avoid)
}

func CheckInAvoid(c *gin.Context) {
	var dailyCheck models.DailyCheck
	if err := c.BindJSON(&dailyCheck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var avoid models.Avoid
	if err := database.DB.First(&avoid, dailyCheck.AvoidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avoid not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if avoid.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to check in this avoid"})
		return
	}

	if time.Since(avoid.StartDate).Hours() > float64(avoid.Duration*24) {
		avoid.IsActive = false
		database.DB.Save(&avoid)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avoid duration has expired"})
		return
	}

	dailyCheck.CheckedDate = time.Now()
	dailyCheck.AvoidID = avoid.ID

	result := database.DB.Create(&dailyCheck)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	avoid.LastCheckedIn = dailyCheck.CheckedDate
	database.DB.Save(&avoid)

	c.JSON(http.StatusOK, dailyCheck)
}

func GetUserAvoids(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var avoids []models.Avoid
	result := database.DB.Where("user_id = ?", userID).Find(&avoids)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, avoids)
}

func GetAvoidDetails(c *gin.Context) {
	avoidID := c.Param("id")

	var avoid models.Avoid
	if err := database.DB.Preload("DailyChecks").First(&avoid, avoidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avoid not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if avoid.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this avoid"})
		return
	}

	c.JSON(http.StatusOK, avoid)
}
