package structs

import (
	"time"

	"github.com/gosquash/api/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	Id   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string    `json:"name"`

	Members []GroupMember `json:"-"`

	CreatorId uuid.UUID `json:"-"`
	Creator   User      `json:"creator" gorm:"foreignKey:CreatorId"`

	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"autoDeleteTime"`

	Games []Game `json:"-"`
}

// Get group by id
func GetGroup(id any) *Group {

	// If string parse to uuid
	if idStr, ok := id.(string); ok {
		id, _ = uuid.Parse(idStr)
	}

	var group Group

	if result := db.DB.
		Preload("Members").
		Preload("Members.User").
		Preload("Creator").
		First(&group, "id = ?", id); result.Error != nil {
		return nil
	}

	return &group
}

// Check if user can edit group
func (g *Group) CanEdit(user *User) bool {
	return g.CreatorId == user.Id
}

// Check if user is a member of the group
func (g *Group) IsMember(user *User) bool {
	for _, m := range g.Members {
		if m.User.Id == user.Id {
			return true
		}
	}

	return false
}

// Get group members
func (g *Group) GetMembers() *[]GroupMember {
	return &g.Members
}

// Get group games
func (g *Group) GetGames() (*[]Game, uint) {

	var games []Game

	if result := db.DB.
		Preload("Players").
		Preload("Players.User").
		Preload("AddedBy").
		Order("created_at desc").
		Find(&games, "group_id = ?", g.Id); result.Error != nil {
		return nil, 0
	}

	g.Games = games

	return &games, uint(len(games))
}

// Add member to group
func (g *Group) AddMember(userId uuid.UUID) error {
	member := GroupMember{
		UserId: userId,
		State:  GroupMemberStatePending,
	}

	if result := db.DB.Create(&member); result.Error != nil {
		return result.Error
	}

	g.Members = append(g.Members, member)

	g.Update()

	return nil
}

// Update group in DB
func (g *Group) Update() error {
	if result := db.DB.Save(g); result.Error != nil {
		return result.Error
	}

	return nil
}
