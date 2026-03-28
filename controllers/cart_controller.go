package controllers

import (
	"github.com/gin-gonic/gin"
	"symphoney/config"
	"symphoney/models"
)

func AddToCart(c *gin.Context) {

	userID := c.GetInt("user_id")

	var cart models.Cart

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// 🔍 CHECK IF PRODUCT EXISTS
	var existingID int
	var existingQty int

	err := config.DB.QueryRow(
		`SELECT id, quantity FROM cart WHERE user_id=$1 AND product_id=$2`,
		userID,
		cart.ProductID,
	).Scan(&existingID, &existingQty)

	if err == nil {
		// ✅ UPDATE QUANTITY
		_, err = config.DB.Exec(
			`UPDATE cart SET quantity=$1 WHERE id=$2`,
			existingQty+cart.Quantity,
			existingID,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update cart"})
			return
		}

		c.JSON(200, gin.H{"message": "Cart updated"})
		return
	}

	// 🆕 INSERT NEW ITEM
	query := `
	INSERT INTO cart (user_id, product_id, quantity)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	err = config.DB.QueryRow(
		query,
		userID,
		cart.ProductID,
		cart.Quantity,
	).Scan(&cart.ID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Product added to cart",
		"data":    cart,
	})
}
func GetCart(c *gin.Context) {

	userID := c.GetInt("user_id")

	rows, err := config.DB.Query(`
	SELECT 
		c.id,
		c.product_id,
		c.quantity,
		p.name,
		p.price
	FROM cart c
	JOIN products p ON c.product_id = p.id
	WHERE c.user_id=$1
	`, userID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch cart"})
		return
	}

	defer rows.Close()

	var cartItems []gin.H

	for rows.Next() {

		var id, productID, quantity int
		var name string
		var price float64

		err := rows.Scan(&id, &productID, &quantity, &name, &price)
		if err != nil {
			continue
		}

		cartItems = append(cartItems, gin.H{
			"id":         id,
			"product_id": productID,
			"quantity":   quantity,
			"name":       name,
			"price":      price,
		})
	}

	c.JSON(200, gin.H{
		"data": cartItems,
	})
}

func UpdateCart(c *gin.Context) {
	id := c.Param("id")
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	query := `
	UPDATE cart
	SET quantity=$1
	WHERE id=$2
	`

	_, err := config.DB.Exec(query, cart.Quantity, id)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(200, gin.H{"message": "Cart updated"})
}

func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	query := `
	DELETE FROM cart
	WHERE id=$1
	`
	_, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete item"})
		return
	}
	c.JSON(200, gin.H{"message": "Item removed from cart"})
}
