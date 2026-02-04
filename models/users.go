package models

import (
	"time"

	"github.com/google/uuid"
)

// UserKind represents type of user account.
type UserKind string

const (
	UserKindUser         UserKind = "user"
	UserKindCTSUser      UserKind = "cts_user"
	UserKindBotX         UserKind = "botx"
	UserKindUnregistered UserKind = "unregistered"
	UserKindGuest        UserKind = "guest"
)

// SyncSourceType represents user sync source.
type SyncSourceType string

const (
	SyncSourceAD     SyncSourceType = "ad"
	SyncSourceAdmin  SyncSourceType = "admin"
	SyncSourceEmail  SyncSourceType = "email"
	SyncSourceOpenID SyncSourceType = "openid"
	SyncSourceBotX   SyncSourceType = "botx"
)

// UserFromSearch represents user profile from search endpoints.
type UserFromSearch struct {
	HUID            uuid.UUID
	ADLogin         string
	ADDomain        string
	Username        string
	Company         string
	CompanyPosition string
	Department      string
	Emails          []string
	OtherID         string
	UserKind        UserKind
	Active          *bool
	Description     string
	IPPhone         string
	Manager         string
	Office          string
	OtherIPPhone    string
	OtherPhone      string
	PublicName      string
	CTSID           *uuid.UUID
	RTSID           *uuid.UUID
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

// UserFromCSV represents user record from CSV export.
type UserFromCSV struct {
	HUID            uuid.UUID
	ADLogin         string
	ADDomain        string
	Username        string
	SyncSource      SyncSourceType
	Active          bool
	UserKind        UserKind
	Email           string
	Company         string
	Department      string
	Position        string
	Avatar          string
	AvatarPreview   string
	Office          string
	Manager         string
	ManagerHUID     *uuid.UUID
	Description     string
	Phone           string
	OtherPhone      string
	IPPhone         string
	OtherIPPhone    string
	PersonnelNumber string
}
