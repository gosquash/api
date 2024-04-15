package groups

import (
	"net/http"

	"github.com/gosquash/api/internal/db"
	"github.com/gosquash/api/pkg/middleware"
	"github.com/gosquash/api/pkg/structs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func Init(g *echo.Group) {
	g.Use(middleware.Authentication)

	g.POST("", createGroup)

	g.GET("/:id", getGroup)

	g.GET("/:id/games", getGroupGames)
	g.POST("/:id/games", createGroupGame)

	g.GET("/:id/members", getGroupMembers)
	g.POST("/:id/members", addGroupMember)
}

func createGroup(c echo.Context) error {

	user := c.Get("user").(*structs.User)

	var body struct {
		Name string `json:"name" validate:"required"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	group := structs.Group{
		Name:      body.Name,
		CreatorId: user.Id,
		Members: []structs.GroupMember{
			{UserId: user.Id, State: structs.GroupMemberStateAdmin},
		},
	}

	// Create group
	db.DB.Create(&group)

	return c.JSON(200, echo.Map{
		"group_id": group.Id,
	})
}

func getGroup(c echo.Context) error {

	user := c.Get("user").(*structs.User)

	group := structs.GetGroup(c.Param("id"))

	if group == nil || !group.IsMember(user) {
		return echo.NewHTTPError(404, "Group not found")
	}

	return c.JSON(200, group)
}

func getGroupGames(c echo.Context) error {

	user := c.Get("user").(*structs.User)

	group := structs.GetGroup(c.Param("id"))

	if group == nil || !group.IsMember(user) {
		return c.JSON(404, echo.Map{
			"message": "Group not found",
		})
	}

	type response struct {
		Games *[]structs.Game `json:"games"`

		Total uint `json:"total"`
		Page  uint `json:"page"`
	}

	games, total := group.GetGames()

	return c.JSON(200, response{
		Games: games,
		Total: total,
	})
}

func createGroupGame(c echo.Context) error {

	user := c.Get("user").(*structs.User)

	group := structs.GetGroup(c.Param("id"))

	if group == nil || !group.IsMember(user) {
		return c.JSON(404, echo.Map{
			"message": "Group not found",
		})
	}

	type player struct {
		UserId uuid.UUID `json:"user_id"`
		Points uint8     `json:"points"`
	}

	var body struct {
		Players []player `json:"players"`
	}

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
		GroupId:   group.Id,
		Players:   players,
		AddedById: user.Id,
	}

	db.DB.Create(&game)

	type response struct {
		GameId string `json:"game_id"`
	}

	return c.JSON(http.StatusCreated, response{
		GameId: game.Id.String(),
	})
}

func getGroupMembers(c echo.Context) error {

	user := c.Get("user").(*structs.User)
	group := structs.GetGroup(c.Param("id"))

	if group == nil || !group.IsMember(user) {
		return c.JSON(404, echo.Map{
			"message": "Group not found",
		})
	}

	type response struct {
		Members *[]structs.GroupMember `json:"members"`

		Total uint `json:"total"`
		Page  uint `json:"page"`
	}

	return c.JSON(200, response{
		Members: group.GetMembers(),
	})
}

// Add a member to a group
func addGroupMember(c echo.Context) error {

	user := c.Get("user").(*structs.User)
	group := structs.GetGroup(c.Param("id"))

	if group == nil || !group.CanEdit(user) {
		return c.JSON(400, echo.Map{
			"message": "You can not add members to this group",
		})
	}

	var body struct {
		UserId uuid.UUID `json:"user_id" validate:"required"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	group.AddMember(body.UserId)

	return c.String(201, "")
}
