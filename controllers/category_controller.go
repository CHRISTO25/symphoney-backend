package controllers

import (
	"net/http"
	"symphoney/config"
	"symphoney/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {

	rows, err := config.DB.Query(`
	SELECT id,name,description,image_url
	FROM categories
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	defer rows.Close()

	var categories []models.Category

	for rows.Next() {

		var category models.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.ImageURL,
		)

		if err != nil {
			continue
		}

		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}
