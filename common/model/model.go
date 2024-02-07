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
	Id                   uint
	IsAdmin              bool
	Birthday             string
	DeviceID             string `json:"device_id"`
	Gender               bool
	FCMToken             string `json:"fcm_token"`
	IsFirst              bool   `json:"is_first"`
	Name                 string
	PhoneNum             string `json:"phone_num"`
	UseAutoLogin         bool   `json:"use_auto_login"`
	UsePrivacyProtection bool   `json:"use_privacy_protection"`
	UseSleepTracking     bool   `json:"use_sleep_tracking"`
	UserType             string `json:"user_type"`
	Email                string
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
	Type   string
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
	Images []Image
	Foods  json.RawMessage `gorm:"type:json"`
}

type Image struct {
	TimestampModel
	Id  uint
	Uid uint

	//부모 아이디
	DietId uint `json:"diet_id"`

	//부모아이디 끝

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
	Id                 uint
	Uid                uint
	Timestamp          json.RawMessage `gorm:"type:json"`
	Weekdays           json.RawMessage `gorm:"type:json"`
	CustomMedicineType string          `json:"custom_medicine_type"`
	Dose               float32
	IntervalType       uint8   `json:"interval_type"`
	IsActive           bool    `json:"is_active"`
	LeastStore         float32 `json:"least_store"`
	MedicineType       string  `json:"medicine_type"`
	Name               string
	Store              float32
	StartAt            string `json:"start_at"`
	EndAt              string `json:"end_at"`
	UseLeastStore      bool   `json:"use_least_store"`
	UsePrivacy         bool   `json:"use_privacy"`
}

type MedicineTakes struct {
	TimestampModel
	Id         uint
	Uid        uint
	DateTaken  string `json:"date_taken"`
	TimeTaken  string `json:"time_taken"`
	Dose       float32
	MedicineId uint `json:"medicine_id"`
}

func (tm *TimestampModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if tm.Created == "" {
		tm.Created = now
	}
	tm.Updated = now
	return
}
