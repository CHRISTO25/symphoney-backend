package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"symphoney/config"
	"symphoney/models"
	"symphoney/utils"
)

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendOTPEmail(to string, otp string) error {
	from := "christovarghese555@gmail.com"
	password := "yvfbjfxuqubpydmg"

	msg := []byte("Subject: Your OTP\n\nYour OTP is: " + otp)

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{to},
		msg,
	)
}

// func RegisterUser(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": "something bad occours"})
// 		return
// 	}
// 	if user.Role == "" {
// 		user.Role = "user"
// 	}
// 	var hashedPassword, err = bcrypt.GenerateFromPassword(
// 		[]byte(user.Password),
// 		bcrypt.DefaultCost,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Password hashing failed",
// 		})
// 		return
// 	}

// 	query := `INSERT INTO users(name, email , password , role)
// 	          VALUES($1,$2,$3,$4)
// 			  RETURNING id`

// 	err = config.DB.QueryRow(
// 		query,
// 		user.Name,
// 		user.Email,
// 		string(hashedPassword),
// 		user.Role,
// 	).Scan(&user.Id)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": "data not inserted in data base",
// 			"realError": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "data inserted",
// 		"status":  "success",
// 		"data":    user,
// 	})
// }

func LoginUser(c *gin.Context) {

	var loginData models.User

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var user models.User

	if user.IsBlocked {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Your account is blocked",
		})
		return
	}

	query := `
	SELECT id,name,email,password,role
	FROM users
	WHERE email=$1
	`

	err := config.DB.QueryRow(
		query,
		loginData.Email,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email",
		})
		return
	}

	// COMPARE PASSWORD
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(loginData.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}

	token, err := utils.GenerateToken(user.Id, user.Email, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token generation failed",
		})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
		"role":   user.Role,
		"user":   user,
	})
}

func Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Message": "Items working properly"})
}
