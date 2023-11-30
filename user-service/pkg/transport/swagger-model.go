package transport

import (
	"time"
)

type User struct {
	ID                   int       `gorm:"primaryKey;autoIncrement"`
	Birthday             time.Time `gorm:"size:40;default:''"`
	DeviceID             string    `gorm:"size:40;default:''"`
	Gender               bool
	FCMToken             string    `gorm:"size:255;default:''" json:"fcm_token"`
	IsFirst              bool      `gorm:"default:false"`
	Name                 string    `gorm:"size:40;default:''"`
	PhoneNum             string    `gorm:"size:40;default:''" json:"phone_num"`
	UseAutoLogin         bool      `gorm:"default:false"`
	UsePrivacyProtection bool      `gorm:"default:false"`
	UseSleepTracking     bool      `gorm:"default:false"`
	UserType             string    `gorm:"size:40;default:''"`
	Email                string    `gorm:"size:40;default:''"`
	Created              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	Updated              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type LoginRequest struct {
	IdToken string `json:"idToken"`
	User    User   `json:"user"`
}

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}
type SuccessResponse struct {
	Jwt string `json:"jwt"`
}

// ErrorResponse represents an error response for the API.
type ErrorResponse struct {
	Err string `json:"err" example:"account name"` // wwwwww
}

type BasicResponse struct {
	Code string `json:"code"`
}
