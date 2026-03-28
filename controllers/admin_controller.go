package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"symphoney/config"
	"symphoney/models"
)

// products ---------------------------------------------

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	query := `
	INSERT INTO products (name,description,price,category_id,stock)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`
	err := config.DB.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryID,
		product.Stock,
	).Scan(&product.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(201, gin.H{
		"message": "Product created",
		"data":    product,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	query := `
	UPDATE products
	SET name=$1,description=$2,price=$3,category_id=$4,stock=$5
	WHERE id=$6
	`
	_, err := config.DB.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryID,
		product.Stock,
		id,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(200, gin.H{"message": "Product updated"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	query := `
	DELETE FROM products
	WHERE id=$1
	`
	_, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted"})
}

//category---------------------------

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

func GetAllUsers(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, email, role, is_blocked
		FROM users
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.IsBlocked,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading users",
			})
			return
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

func BlockUser(c *gin.Context) {
	id := c.Param("id")

	query := `UPDATE users SET is_blocked = TRUE WHERE id = $1`

	_, err := config.DB.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to block user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User blocked successfully",
	})
}

func UnblockUser(c *gin.Context) {
	id := c.Param("id")

	query := `UPDATE users SET is_blocked = FALSE WHERE id = $1`

	_, err := config.DB.Exec(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unblock user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User unblocked successfully",
	})
}
