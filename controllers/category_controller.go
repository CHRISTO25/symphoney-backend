package controllers

import (
	"net/http"
	"symphoney/config"
	"symphoney/models"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {

	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	query := `
	INSERT INTO categories (name,description,image_url)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	err := config.DB.QueryRow(
		query,
		category.Name,
		category.Description,
		category.ImageURL,
	).Scan(&category.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created",
		"data":    category,
	})
}

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

func UpdateCategory(c *gin.Context) {

	id := c.Param("id")

	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid input",
		})
		return
	}

	query := `
	UPDATE categories
	SET name=$1, description=$2, image_url=$3
	WHERE id=$4
	`

	_, err := config.DB.Exec(
		query,
		category.Name,
		category.Description,
		category.ImageURL,
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to update category",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Category updated successfully",
	})
}

func DeleteCategory(c *gin.Context) {

	id := c.Param("id")

	query := `
	DELETE FROM categories
	WHERE id=$1
	`

	_, err := config.DB.Exec(query, id)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to delete category",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Category deleted successfully",
	})
}
