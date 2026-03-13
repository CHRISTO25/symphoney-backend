package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"symphoney/config"
	"symphoney/models"
)

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

func GetProducts(c *gin.Context) {

	rows, err := config.DB.Query(`
	SELECT id,name,description,price,category_id,stock
	FROM products
	`)

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

		imageRows, _ := config.DB.Query(`
		SELECT image_url
		FROM product_images
		WHERE product_id=$1
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

	c.JSON(200, gin.H{"data": product})
}

func AddProductImage(c *gin.Context) {

	idParam := c.Param("id")

	productID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product id"})
		return
	}

	var image models.ProductImage

	if err := c.ShouldBindJSON(&image); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	query := `
	INSERT INTO product_images (product_id,image_url)
	VALUES ($1,$2)
	RETURNING id
	`

	err = config.DB.QueryRow(
		query,
		productID,
		image.ImageURL,
	).Scan(&image.ID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	image.ProductID = productID

	c.JSON(201, gin.H{
		"message": "Image added",
		"data":    image,
	})
}
