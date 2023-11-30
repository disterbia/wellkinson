package model

import (
	"time"
)

type User struct {
	Id                   int       `gorm:"primaryKey;autoIncrement"`
	isAdmin              bool      `gorm:"not null;default:false" json:"-"`
	Birthday             string    `gorm:"size:40;not null"`
	DeviceID             string    `gorm:"size:40;not null"`
	Gender               bool      `gorm:"not null"`
	FCMToken             string    `gorm:"size:255;not null"`
	IsFirst              bool      `gorm:"not null;default:true"`
	Name                 string    `gorm:"size:40;not null"`
	PhoneNum             string    `gorm:"size:40;not null"`
	UseAutoLogin         bool      `gorm:"not null;default:false"`
	UsePrivacyProtection bool      `gorm:"not null;default:false"`
	UseSleepTracking     bool      `gorm:"not null;default:false"`
	UserType             string    `gorm:"size:40;not null"`
	Email                string    `gorm:"size:40;not null"`
	Created              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	Updated              time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

func (u *User) IsAdmin() bool { // json: "-" 하면됨
	return u.isAdmin
}

type Alarm struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	Uid       int       `gorm:"not null"`
	Type      string    `gorm:"size:255;not null"`
	Body      string    `gorm:"type:text;not null"`
	StartAt   string    `gorm:"size:255;not null"`
	EndAt     string    `gorm:"size:255;not null"`
	Timestamp string    `gorm:"size:255;not null"`
	Week      string    `gorm:"size:255;not null"`
	Created   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	Updated   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

type Notification struct {
	Id      int       `gorm:"primaryKey;autoIncrement"`
	Uid     int       `gorm:"not null"`
	Type    string    `gorm:"size:40;not null"`
	Body    string    `gorm:"type:text;not null"`
	IsRead  bool      `gorm:"not null;default:false"`
	Created time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	Updated time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

type Inquire struct {
	Id      int `gorm:"primaryKey;autoIncrement"`
	Uid     int //user 아이디
	Email   string
	Title   string
	Content string
	Created time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	Updated time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

type InquireReply struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	Uid       int       `gorm:"not null"`
	InquireId int       `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	replyType bool      `gorm:"not null"`
	Created   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"-"`
	Updated   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"-"`
}

type BasicResponse struct {
	Code string `json:"code"`
}

type LoginRequest struct {
	IdToken string `json:"idToken"`
	User    User   `json:"user"`
}

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}

type AutoLoginRequest struct {
	Jwt string `json:"jwt"`
}
