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
	Id                    uint
	IsAdmin               bool
	Birthday              string
	DeviceID              string `json:"device_id"`
	Gender                bool
	IndemnificationClause bool   `json:"indemnification_clause"`
	FCMToken              string `json:"fcm_token"`
	IsFirst               bool   `json:"is_first"`
	Name                  string
	PhoneNum              string `json:"phone_num"`
	UseAutoLogin          bool   `json:"use_auto_login"`
	UsePrivacyProtection  bool   `json:"use_privacy_protection"`
	UseSleepTracking      bool   `json:"use_sleep_tracking"`
	UserType              uint   `json:"user_type"`
	Email                 string
	SnsType               uint          `json:"sns_type"`
	ProfileImage          Image         `json:"profile_image" gorm:"foreignkey:ParentId"`
	LinkedEmails          []LinkedEmail `json:"linked_emails" gorm:"foreignkey:Uid"`
}

type Alarm struct {
	TimestampModel
	Id        uint
	Uid       uint
	ParentId  uint `json:"parent_id"`
	Type      uint
	Body      string
	StartAt   string ` json:"start_at"`
	EndAt     string ` json:"end_at"`
	Timestamp string
	Week      json.RawMessage `gorm:"type:json"`
}

type Notification struct {
	TimestampModel
	Id     uint
	Uid    uint
	Type   uint
	Body   string
	IsRead bool `json:"is_read"`
}

type Inquire struct {
	TimestampModel
	Id      uint
	Uid     uint
	Email   string
	Title   string
	Content string
	Replies []InquireReply
}

type InquireReply struct {
	TimestampModel
	Id        uint
	Uid       uint
	InquireId uint `json:"inquire_id"`
	ReplyType bool `json:"reply_type"`
	Content   string
}

type DietPreset struct {
	TimestampModel
	Id    uint
	Uid   uint
	Name  string
	Foods json.RawMessage `gorm:"type:json"`
}

type Diet struct {
	TimestampModel
	Id     uint
	Uid    uint
	Name   string
	Time   string
	Type   uint
	Images []Image         `gorm:"foreignkey:ParentId"`
	Foods  json.RawMessage `gorm:"type:json"`
}

type Image struct {
	TimestampModel
	Id  uint
	Uid uint

	//부모 아이디
	ParentId uint `json:"parent_id"`
	Type     uint

	Url          string
	ThumbnailUrl string `json:"thumbnail_url"`
}

type Emotion struct {
	TimestampModel
	Id      uint
	Uid     uint
	Emotion string
	State   string
}

type Exercise struct {
	TimestampModel
	Id              uint
	Uid             uint
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
	Id            uint
	Uid           uint
	DatePerformed string `json:"date_performed"`
	ExerciseId    uint   `json:"exercise_id"`
}

type FaceScore struct {
	TimestampModel
	Id    uint
	Uid   uint
	Score uint
	Type  uint
}

type FaceExam struct {
	TimestampModel
	Id      uint
	Type    uint
	Title   string
	VideoId string `json:"video_id"`
}

type FaceExercise struct {
	TimestampModel
	Id      uint
	Type    uint
	Title   string
	VideoId string `json:"video_id"`
}

type Video struct {
	TimestampModel
	Id           uint
	ProjectName  string `json:"project_name"`
	Name         string
	Duration     uint
	ProjectId    string `json:"project_id"`
	VideoId      string `json:"video_id"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type Medicine struct {
	TimestampModel
	Id           uint
	Uid          uint
	Timestamp    json.RawMessage `gorm:"type:json"`
	Weekdays     json.RawMessage `gorm:"type:json"`
	Dose         float32
	IntervalType uint8   `json:"interval_type"`
	IsActive     bool    `json:"is_active"`
	LeastStore   float32 `json:"least_store"`
	MedicineType string  `json:"medicine_type"`
	Name         string
	Store        float32
	StartAt      string `json:"start_at"`
	EndAt        string `json:"end_at"`
	UsePrivacy   bool   `json:"use_privacy"`
}

type MedicineTake struct {
	TimestampModel
	Id         uint
	Uid        uint
	DateTaken  string `json:"date_taken"`
	TimeTaken  string `json:"time_taken"`
	Dose       float32
	MedicineId uint `json:"medicine_id"`
}

type MedicineSearch struct {
	TimestampModel
	Id   uint
	Name string
}

type SleepAlarm struct {
	TimestampModel
	Id        uint
	Uid       uint
	StartTime string          `json:"start_time"`
	AlarmTime string          `json:"alarm_time"`
	EndTime   string          `json:"end_time"`
	Weekdays  json.RawMessage `gorm:"type:json"`
	IsActive  bool            `json:"is_active"`
}

type SleepTime struct {
	TimestampModel
	Id        uint
	Uid       uint
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	DateSleep string `json:"date_sleep"`
}

type VocalWord struct {
	TimestampModel
	Id    uint
	Type  uint
	Title string
}

type VocalScore struct {
	TimestampModel
	Id    uint
	Uid   uint
	Score uint
	Type  uint
}

type MainService struct {
	TimestampModel
	Id    uint
	Title string
	Level uint
}

type UserService struct {
	TimestampModel
	Id        uint
	Uid       uint
	ServiceId uint `json:"service_id"`
	Title     string
}

type AuthCode struct {
	TimestampModel
	Id          uint
	PhoneNumber string
	Code        string
}

type VerifiedNumbers struct {
	TimestampModel
	Id          uint
	PhoneNumber string
}

type LinkedEmail struct {
	TimestampModel
	Id      uint
	Email   string
	Uid     uint
	SnsType uint `json:"sns_type"`
}

type AppVersion struct {
	TimestampModel
	Id            uint
	LatestVersion string `json:"latest_version"`
	AndroidLink   string `json:"android_link"`
	IosLink       string `json:"ios_link"`
}

type Polices struct {
	TimestampModel
	Id         uint
	Title      string `json:"title"`
	Body       string `json:"body"`
	PoliceType uint   `json:"police_type"`
}

func (tm *TimestampModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if tm.Created == "" {
		tm.Created = now
	}
	tm.Updated = now
	return
}
