package models

import "github.com/google/uuid"

// SmartApp represents smartapp from list.
type SmartApp struct {
	AppID         string
	Enabled       bool
	ID            uuid.UUID
	Name          string
	Avatar        string
	AvatarPreview string
}

// SmartAppManifestWebLayout defines layout choice.
type SmartAppManifestWebLayout string

const (
	SmartAppWebMinimal SmartAppManifestWebLayout = "minimal"
	SmartAppWebHalf    SmartAppManifestWebLayout = "half"
	SmartAppWebFull    SmartAppManifestWebLayout = "full"
)

// SmartAppManifestWebParams represents web params.
type SmartAppManifestWebParams struct {
	DefaultLayout  SmartAppManifestWebLayout   `json:"default_layout"`
	ExpandedLayout SmartAppManifestWebLayout   `json:"expanded_layout"`
	AllowedLayouts []SmartAppManifestWebLayout `json:"allowed_layouts,omitempty"`
	AlwaysPinned   bool                        `json:"always_pinned"`
}

// SmartAppManifestMobileParams represents mobile params.
type SmartAppManifestMobileParams struct {
	FullscreenLayout bool `json:"fullscreen_layout"`
}

// SmartAppManifestUnreadCounterParams represents unread counter link params.
type SmartAppManifestUnreadCounterParams struct {
	UserHUID    []uuid.UUID `json:"user_huid"`
	GroupChatID []uuid.UUID `json:"group_chat_id"`
	AppID       []string    `json:"app_id"`
}

// SmartAppManifest represents manifest.
type SmartAppManifest struct {
	IOS               *SmartAppManifestMobileParams        `json:"ios,omitempty"`
	Android           *SmartAppManifestMobileParams        `json:"android,omitempty"`
	Aurora            *SmartAppManifestMobileParams        `json:"aurora,omitempty"`
	Web               *SmartAppManifestWebParams           `json:"web,omitempty"`
	UnreadCounterLink *SmartAppManifestUnreadCounterParams `json:"unread_counter_link,omitempty"`
}
