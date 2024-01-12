package controllers

import (
	"api-go/initializers"
	"api-go/models"
	"api-go/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email           string
		Password        string
		ConfirmPassword string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	if utils.IsEmpty(body.Email) || !utils.IsEmailValid(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
		})

		return
	}

	if utils.IsEmpty(body.Password) || len(body.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must have at least 6 characters",
		})

		return
	}

	if body.Password != body.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password and password confirmation must be the same",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": struct {
			ID    uint
			Email string
		}{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	if utils.IsEmpty(body.Email) || utils.IsEmpty(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	var user models.User

	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func EditEmail(c *gin.Context) {

	currentUser, _ := c.Get("user")

	var body struct {
		NewEmail string `json:"newEmail"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	if utils.IsEmpty(body.NewEmail) || !utils.IsEmailValid(body.NewEmail) || currentUser.(models.User).Email == body.NewEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
		})

		return
	}

	result := initializers.DB.Model(&currentUser).Where("id = ?", currentUser.(models.User).ID).Update("Email", body.NewEmail)

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update email",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": struct {
			ID    uint
			Email string
		}{
			ID:    currentUser.(models.User).ID,
			Email: body.NewEmail,
		},
	})
}

func EditPassword(c *gin.Context) {
	currentUser, _ := c.Get("user")

	var body struct {
		NewPassword        string `json:"newPassword"`
		ConfirmNewPassword string `json:"confirmNewPassword"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	if utils.IsEmpty(body.NewPassword) || len(body.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must have at least 6 characters",
		})

		return
	}

	if body.NewPassword != body.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password and password confirmation must be the same",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	result := initializers.DB.Model(&currentUser).Where("id = ?", currentUser.(models.User).ID).Update("Password", string(hash))

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update password",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": struct {
			ID    uint
			Email string
		}{
			ID:    currentUser.(models.User).ID,
			Email: currentUser.(models.User).Email,
		},
	})
}

func Validate(c *gin.Context) {
	currentUser, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"user": struct {
			ID    uint
			Email string
		}{
			ID:    currentUser.(models.User).ID,
			Email: currentUser.(models.User).Email,
		},
	})
}
