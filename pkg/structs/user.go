package structs

import (
	"gosquash/api/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name     string    `json:"name"`
	Email    string    `json:"-" gorm:"unique"`
	Password string    `json:"-"`

	Players      []Player      `json:"-"`
	GroupMembers []GroupMember `json:"-"`
}

func GetUser(id any) *User {

	// If string parse to uuid
	if idStr, ok := id.(string); ok {
		id, _ = uuid.Parse(idStr)
	}

	var user User

	if result := db.DB.First(&user, "id = ?", id); result.Error != nil {
		return nil
	}

	return &user
}

func (u *User) GetGroups() *[]Group {

	var groups []Group

	db.DB.Model(&u).Association("Groups").Find(&groups)

	return &groups
}

type UserClaims = jwt.RegisteredClaims

func GetUserByToken(token string) *User {

	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil
	}

	claims, ok := jwtToken.Claims.(*UserClaims)

	if !ok || !jwtToken.Valid {
		return nil
	}

	var user User

	if result := db.DB.Where("id = ?", claims.Subject).First(&user); result.Error != nil {
		return nil
	}

	return &user
}
