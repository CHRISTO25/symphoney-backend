package controllers

import (
	"fmt"
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

	fmt.Println(userID)

	query := `
	INSERT INTO cart (user_id,product_id,quantity)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	cart.UserID = userID
	err := config.DB.QueryRow(
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
	SELECT id,user_id,product_id,quantity
	FROM cart
	WHERE user_id=$1
	`, userID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch cart"})
		return
	}

	defer rows.Close()

	var cartItems []models.Cart

	for rows.Next() {

		var item models.Cart

		rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
		)

		cartItems = append(cartItems, item)
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
