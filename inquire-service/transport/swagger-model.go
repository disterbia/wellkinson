package transport

import (
	"time"
)

type User struct {
	ID                   int       `gorm:"primaryKey;autoIncrement"  example:"1"`
	Birthday             string    `gorm:"size:40;default:''" example:"yyyy-mm-dd"`
	DeviceID             string    `gorm:"size:40;default:''"`
	Gender               bool      `gorm:"size:10;default:''"`
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

type InquireReply struct {
	ID        int  `gorm:"primaryKey;autoIncrement"`
	UID       int  //user 아이디
	InquireId int  // Inquire 아이디
	Type      bool //답변 또는 추가 문의
	Content   string
	Date      time.Time
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}

type ErrorResponse struct {
	Err string `json:"err"` // wwwwww
}

type BasicResponse struct {
	Code string `json:"code"`
}
