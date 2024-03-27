package structs

import (
	"gosquash/api/internal/db"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Game is a struct that represents a game
type Game struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Players []Player  `json:"players"`

	AddedById uuid.UUID `json:"-"`
	AddedBy   User      `json:"added_by" gorm:"foreignKey:AddedById"`

	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"autoDeleteTime"`
}

// NewGame creates a new game
func GetGame(id any) *Game {

	// If string parse to uuid
	if idStr, ok := id.(string); ok {
		id, _ = uuid.Parse(idStr)
	}

	var game Game

	if result := db.DB.Preload("Players").Preload("Players.User").Preload("AddedBy").First(&game, "id = ?", id); result.Error != nil {
		return nil
	}

	return &game
}

func (g *Game) CanEdit(user *User) bool {
	return g.AddedById == user.Id
}

func (g *Game) CanSee(user *User) bool {
	for _, player := range g.Players {
		if player.UserId == user.Id {
			return true
		}
	}

	return g.CanEdit(user)
}

// getWinner returns the player with the most points
func (g *Game) getWinner() Player {

	winner := g.Players[0]

	for _, player := range g.Players {
		if player.Points > winner.Points {
			winner = player
		}
	}

	return winner
}

// Player is a struct that represents a player in the game
type Player struct {
	Id uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	GameId uuid.UUID `json:"-" gorm:"type:uuid"`
	UserId uuid.UUID `json:"-" gorm:"type:uuid"`
	User   User      `json:"user" gorm:"foreignKey:UserId"`

	Points uint8 `json:"points" gorm:"default:0"`
}
