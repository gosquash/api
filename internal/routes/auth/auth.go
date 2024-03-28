package auth

import (
	"errors"
	"gosquash/api/internal/db"
	"gosquash/api/pkg/structs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Init(g *echo.Group) {

	g.POST("/login", login)
	g.POST("/register", register)
}

func login(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	var user structs.User

	// Find user based on email
	if result := db.DB.Where("email = ?", body.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(400, echo.Map{
				"message": "Invalid credentials",
			})
		}
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.JSON(400, echo.Map{
			"message": "Invalid credentials",
		})
	}

	// Create JWT token
	token := createToken(user.Id.String())

	return c.JSON(200, echo.Map{
		"token": token,
	})
}

func register(c echo.Context) error {
	var body struct {
		Name     string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	if err := c.Validate(&body); err != nil {
		return err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	newUser := &structs.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashPassword),
	}

	if result := db.DB.Create(newUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return c.JSON(400, echo.Map{
				"message": "Email already exists",
			})
		}

		return c.JSON(400, echo.Map{
			"message": "Error creating user",
			"error":   result.Error.Error(),
		})
	}

	return c.String(204, "")
}

// Create JWT token
func createToken(userId string) string {

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &structs.UserClaims{
		Issuer:    "https://auth.gosquash.gg",
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
	})

	token, err := jwtToken.SignedString([]byte("secret"))

	if err != nil {
		panic(err)
	}

	return token
}
