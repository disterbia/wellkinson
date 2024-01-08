package transport

import (
	"time"
)

type User struct {
	ID                   int       `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Birthday             string    `gorm:"size:40;default:''" json:"birthday" example:"yyyy-mm-dd"`
	DeviceID             string    `gorm:"size:40;not null" json:"device_id"`
	Gender               bool      `gorm:"not null" json:"gender"`
	FCMToken             string    `gorm:"size:255;not null" json:"fcm_token"`
	IsFirst              bool      `gorm:"not null;default:true" json:"is_first"`
	Name                 string    `gorm:"size:40;not null" json:"name"`
	PhoneNum             string    `gorm:"size:40;not null" json:"phone_num"`
	UseAutoLogin         bool      `gorm:"not null;default:false" json:"use_auto_login"`
	UsePrivacyProtection bool      `gorm:"not null;default:false" json:"user_privacy_protection"`
	UseSleepTracking     bool      `gorm:"not null;default:false" json:"use_sleep_tracking"`
	UserType             string    `gorm:"size:40;not null" json:"user_type"`
	Email                string    `gorm:"size:40;default:''" json:"email"`
	Created              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created"`
	Updated              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated"`
}

type Alarm struct {
	Id        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Uid       int       `gorm:"not null" json:"uid"`
	Type      string    `gorm:"size:255;not null" json:"type"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	StartAt   string    `gorm:"size:255;not null" json:"start_at"`
	EndAt     string    `gorm:"size:255;not null" json:"end_at"`
	Timestamp string    `gorm:"size:255;not null" json:"timestamp"`
	Week      string    `gorm:"size:255;not null" json:"week"`
	Created   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created"`
	Updated   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated"`
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
