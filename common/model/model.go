package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TimestampModel struct {
	Created string
	Updated string
}
type User struct {
	TimestampModel
	Id                   int    `gorm:"primaryKey;autoIncrement"`
	IsAdmin              bool   `gorm:"not null;default:false"`
	Birthday             string `gorm:"size:40;not null"`
	DeviceID             string `gorm:"size:40;not null" json:"device_id"`
	Gender               bool   `gorm:"not null"`
	FCMToken             string `gorm:"size:255;not null" json:"fcm_token"`
	IsFirst              bool   `gorm:"not null;default:true"  json:"is_first"`
	Name                 string `gorm:"size:40;not null"`
	PhoneNum             string `gorm:"size:40;not null"  json:"phone_num"`
	UseAutoLogin         bool   `gorm:"not null;default:false"  json:"use_auto_login"`
	UsePrivacyProtection bool   `gorm:"not null;default:false" json:"use_privacy_protection"`
	UseSleepTracking     bool   `gorm:"not null;default:false" json:"use_sleep_tracking"`
	UserType             string `gorm:"size:40;not null" json:"user_type"`
	Email                string `gorm:"size:40;not null"`
}

type Alarm struct {
	TimestampModel
	Id        int    `gorm:"primaryKey;autoIncrement"`
	Uid       int    `gorm:"not null"`
	Type      string `gorm:"size:255;not null"`
	Body      string `gorm:"type:text;not null"`
	StartAt   string `gorm:"size:255;not null" json:"start_at"`
	EndAt     string `gorm:"size:255;not null" json:"end_at"`
	Timestamp string `gorm:"size:255;not null"`
	Week      string `gorm:"size:255;not null"`
}

type Notification struct {
	TimestampModel
	Id     int    `gorm:"primaryKey;autoIncrement"`
	Uid    int    `gorm:"not null"`
	Type   string `gorm:"size:40;not null"`
	Body   string `gorm:"type:text;not null"`
	IsRead bool   `gorm:"not null;default:false" json:"is_read"`
}

type Inquire struct {
	TimestampModel
	Id      int
	Uid     int
	Email   string
	Title   string
	Content string
	Replies []InquireReply
}

type InquireReply struct {
	TimestampModel
	Id        int
	Uid       int
	InquireId int  `json:"inquire_id"`
	ReplyType bool `json:"reply_type"`
	Content   string
}

type DietPreset struct {
	TimestampModel
	Id    int
	Uid   int
	Name  string
	Foods FoodSlice `gorm:"type:json"`
}

type FoodSlice []string

// Scan - sql.Scanner 구현
func (f *FoodSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	result := []string{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}

	*f = result
	return nil
}

// Value - driver.Valuer 구현
func (f FoodSlice) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f)
}

func (tm *TimestampModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if tm.Created == "" {
		tm.Created = now
	}
	tm.Updated = now
	return
}
