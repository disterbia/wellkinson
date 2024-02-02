package model

import (
	"encoding/json"
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
	Id        int             `gorm:"primaryKey;autoIncrement"`
	Uid       int             `gorm:"not null"`
	Type      string          `gorm:"size:255;not null"`
	Body      string          `gorm:"type:text;not null"`
	StartAt   string          `gorm:"size:255;not null" json:"start_at"`
	EndAt     string          `gorm:"size:255;not null" json:"end_at"`
	Timestamp string          `gorm:"size:255;not null"`
	Week      json.RawMessage `gorm:"type:json"`
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
	Foods json.RawMessage `gorm:"type:json"`
}

type Diet struct {
	TimestampModel
	Id     int
	Uid    int
	Name   string
	Time   string
	Type   int
	Images []Image
	Foods  json.RawMessage `gorm:"type:json"`
}

type Image struct {
	TimestampModel
	Id  int
	Uid int

	//부모 아이디
	DietId int `json:"diet_id"`

	//부모아이디 끝

	Url          string
	ThumbnailUrl string `json:"thumbnail_url"`
}

type Emotion struct {
	TimestampModel
	Id      int
	Uid     int
	Emotion string
	State   string
}

type Exercise struct {
	TimestampModel
	Id              int
	Uid             int
	Title           string          `json:"title"`
	ExerciseStartAt string          `json:"exercise_start_at"`
	ExerciseEndAt   string          `json:"exercise_end_at"`
	PlanStartAt     string          `json:"plan_start_at"`
	PlanEndAt       string          `json:"plan_end_at"`
	UseAlarm        bool            `json:"use_alarm"`
	Weekdays        json.RawMessage `gorm:"type:json"`
}

type ExerciseInfo struct {
	TimestampModel
	Id            int
	Uid           int
	DatePerformed string `json:"date_performed"`
	ExerciseId    int    `json:"exercise_id"`
}

type FaceScores struct {
	TimestampModel
	Id    int
	Uid   int
	Score int
	Type  int
}

type FaceExams struct {
	TimestampModel
	Id      int
	Type    int
	Title   string
	VideoId string `json:"video_id"`
}

type Videos struct {
	TimestampModel
	Id           int
	ProjectName  string `json:"project_name"`
	Name         string
	Duration     int
	ProjectId    string `json:"project_id"`
	VideoId      string `json:"video_id"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

func (tm *TimestampModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if tm.Created == "" {
		tm.Created = now
	}
	tm.Updated = now
	return
}
