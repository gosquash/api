package structs

import "github.com/google/uuid"

// Group Member
type GroupMember struct {
	// Permissions int `json:"permissions"`

	GroupId uuid.UUID `json:"-" gorm:"type:uuid;primaryKey"`
	Group   Group     `json:"-" gorm:"foreignKey:GroupId"`

	UserId uuid.UUID `json:"-" gorm:"type:uuid;primaryKey"`
	User   User      `json:"user" gorm:"foreignKey:UserId"`

	State GroupMemberState `json:"state" gorm:"default:0"`
}

type GroupMemberState int

const (
	GroupMemberStatePending GroupMemberState = iota
	GroupMemberStateAccepted
	GroupMemberStateDeclined
	GroupMemberStateAdmin
)
