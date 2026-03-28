package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
	"symphoney/config"
	"symphoney/models"
	"time"
)

// ✅ GET ALL PRODUCTS
func GetProducts(c *gin.Context) {

	search := c.Query("search")
	categoryID := c.Query("category_id")
	page := c.DefaultQuery("page", "1")

	limit := 10
	pageInt, _ := strconv.Atoi(page)
	offset := (pageInt - 1) * limit

	query := `
	SELECT id,name,description,price,category_id,stock
	FROM products
	WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	if search != "" {
		query += " AND name ILIKE $" + strconv.Itoa(argID)
		args = append(args, "%"+search+"%")
		argID++
	}

	if categoryID != "" {
		query += " AND category_id = $" + strconv.Itoa(argID)
		args = append(args, categoryID)
		argID++
	}

	query += " ORDER BY id LIMIT $" + strconv.Itoa(argID) + " OFFSET $" + strconv.Itoa(argID+1)

	args = append(args, limit, offset)

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CategoryID,
			&product.Stock,
		)

		if err != nil {
			continue
		}

		// 🔥 LOAD IMAGES
		imageRows, _ := config.DB.Query(`
			SELECT image_url FROM product_images WHERE product_id=$1
		`, product.ID)

		var images []string

		for imageRows.Next() {
			var url string
			imageRows.Scan(&url)
			images = append(images, url)
		}

		product.Images = images
		products = append(products, product)
	}

	c.JSON(200, gin.H{
		"data": products,
	})
}

// ✅ GET PRODUCT BY ID
func GetProductByID(c *gin.Context) {

	id := c.Param("id")

	var product models.Product

	query := `
	SELECT id,name,description,price,category_id,stock
	FROM products
	WHERE id=$1
	`

	err := config.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryID,
		&product.Stock,
	)

	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	// 🔥 LOAD IMAGES
	imageRows, _ := config.DB.Query(`
		SELECT image_url FROM product_images WHERE product_id=$1
	`, product.ID)

	var images []string

	for imageRows.Next() {
		var url string
		imageRows.Scan(&url)
		images = append(images, url)
	}

	product.Images = images

	c.JSON(200, gin.H{
		"data": product,
	})
}

// ✅ ADD IMAGE
func AddProductImage(c *gin.Context) {

	idParam := c.Param("id")

	productID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product id"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Image required"})
		return
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	filePath := "uploads/" + fileName

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save image"})
		return
	}

	imageURL := "/uploads/" + fileName

	var imageID int

	query := `
	INSERT INTO product_images (product_id,image_url)
	VALUES ($1,$2)
	RETURNING id
	`

	err = config.DB.QueryRow(query, productID, imageURL).Scan(&imageID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Image uploaded",
		"image":   imageURL,
	})
}
