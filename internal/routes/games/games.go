package games

import (
	"net/http"

	"github.com/gosquash/api/internal/db"
	"github.com/gosquash/api/pkg/middleware"
	"github.com/gosquash/api/pkg/structs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var games = []structs.Game{}

func Init(g *echo.Group) {

	g.Use(middleware.Authentication)

	g.GET("", getGamesByUser)
	g.POST("", createGame)

	g.GET("/:id", getGameById)
	// g.PATCH("/:id", updateGameById)
	g.DELETE("/:id", deleteGame)
}

// Get all games.
func getGamesByUser(c echo.Context) error {

	user := c.Get("user").(*structs.User)

	var games []structs.Game

	if result := db.DB.Model(&structs.Game{}).Preload("Players").Preload("Players.User").Preload("AddedBy").Joins("JOIN players ON players.user_id = ?", user.Id).Where("players.game_id = games.id").Find(&games); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Error fetching games",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"games": &games,
	})
}

// Create a game.
func createGame(c echo.Context) error {
	type Player struct {
		UserId uuid.UUID `json:"user_id"`
		Points uint8     `json:"points"`
	}

	var body struct {
		Players []Player `json:"players"`
	}

	user := c.Get("user").(*structs.User)

	if err := c.Bind(&body); err != nil {
		return err
	}

	var players []structs.Player = []structs.Player{}

	for _, player := range body.Players {
		players = append(players, structs.Player{
			UserId: player.UserId,
			Points: player.Points,
		})
	}

	game := structs.Game{
		AddedById: user.Id,
		Players:   players,
	}

	db.DB.Create(&game)

	type Response struct {
		GameId string `json:"game_id"`
	}

	return c.JSON(http.StatusCreated, Response{
		GameId: game.Id.String(),
	})
}

// Get a game by id.
func getGameById(c echo.Context) error {

	game := structs.GetGame(c.Param("id"))
	user := c.Get("user").(*structs.User)

	if game.CanSee(user) == false {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Game not found",
		})
	}

	if game == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Error fetching game",
		})
	}

	return c.JSON(http.StatusOK, game)
}

// // Update a game by id.
// func updateGameById(c echo.Context) error {

// 	game := structs.GetGame(c.Param("id"))
// }

// Delete a game by id.
func deleteGame(c echo.Context) error {

	game := structs.GetGame(c.Param("id"))
	user := c.Get("user").(*structs.User)

	if game.CanEdit(user) == false {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Game not found",
		})
	}

	if result := db.DB.Delete(&game, game.Id); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Error deleting game",
		})
	}

	return c.String(http.StatusNoContent, "")
}
