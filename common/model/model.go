package model

import (
	"time"
)

type User struct {
	Id                   int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IsAdmin              bool      `gorm:"not null;default:false" json:"-"`
	Birthday             string    `gorm:"size:40;not null" json:"birthday"`
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
	Email                string    `gorm:"size:40;not null" json:"email"`
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

type Notification struct {
	Id      int       `gorm:"primaryKey;autoIncrement"`
	Uid     int       `gorm:"not null"`
	Type    string    `gorm:"size:40;not null"`
	Body    string    `gorm:"type:text;not null"`
	IsRead  bool      `gorm:"not null;default:false" json:"is_read"`
	Created time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	Updated time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

type Inquire struct {
	Id      int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Uid     int            `gorm:"not null" json:"uid"` //user 아이디
	Email   string         `gorm:"size:40;not null" json:"email"`
	Title   string         `gorm:"size:255;not null" json:"title"`
	Content string         `gorm:"type:text;not null" json:"content"`
	Created time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created"`
	Updated time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated"`
	Replies []InquireReply `gorm:"foreignKey:InquireId" json:"replies"`
}

type InquireReply struct {
	Id        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Uid       int       `gorm:"not null" json:"uid"`
	InquireId int       `gorm:"not null" json:"inquire_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ReplyType bool      `gorm:"not null" json:"reply_type"`
	Created   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created"`
	Updated   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated"`
}

type BasicResponse struct {
	Code string `json:"code"`
}

type LoginRequest struct {
	IdToken string `json:"id_token"`
	User    User   `json:"user"`
}

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}

type AutoLoginRequest struct {
	Jwt string `json:"jwt"`
}
