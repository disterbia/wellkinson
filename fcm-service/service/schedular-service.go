package service

import (
	"context"
	"encoding/json"
	"fcm-service/common/model"
	"fcm-service/common/util"
	"log"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

var firebaseClient *messaging.Client

func StartCentralCronScheduler(db *gorm.DB) {
	initializeFirebase()

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 * * * * *", func() {
		sendPendingNotifications(db)
	})
	if err != nil {
		log.Fatalf("Failed to create cron job: %v", err)
	}

	c.Start()
}

func initializeFirebase() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("./firebase-adminkey.json"))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
		return
	}
	firebaseClient = client
}

func sendPendingNotifications(db *gorm.DB) {
	now := time.Now()

	var alarms []model.Alarm
	if err := db.Preload("User").Where("(start_at = '' OR start_at <= ?) AND (end_at = '' OR end_at >= ?)",
		now.Format("2006-01-02"), now.Format("2006-01-02")).Find(&alarms).Error; err != nil {
		return
	}

	// 2. 사용자별 미확인 알림 카운트 조회
	notificationCounts := getUnreadNotificationCounts(db)

	// 3. 메시지와 저장할 알림을 준비
	var messages []*messaging.Message
	var newNotifications []model.Notification

	for _, alarm := range alarms {

		if shouldSendNotification(now, alarm) {
			notificationCount := notificationCounts[alarm.Uid]

			// FCM 메시지 생성
			message := &messaging.Message{
				Data: map[string]string{
					"uid":                strconv.FormatUint(uint64(alarm.Uid), 10),
					"type":               strconv.FormatUint(uint64(alarm.Type), 10),
					"notification_count": strconv.FormatUint(uint64(notificationCount), 10),
					"timestamp":          time.Now().Format(time.RFC3339),
					"parent_id":          strconv.FormatUint(uint64(alarm.ParentId), 10),
				},
				Notification: &messaging.Notification{
					Title: getNotificationTitle(alarm.Type),
					Body:  alarm.Body,
				},
				Token: alarm.User.FCMToken,
			}
			messages = append(messages, message)

			// DB에 저장할 알림 준비
			newNotifications = append(newNotifications, model.Notification{
				Uid:       alarm.Uid,
				Type:      alarm.Type,
				Body:      alarm.Body,
				ParentId:  alarm.ParentId,
				IsRead:    false,
				Timestamp: alarm.Timestamp,
			})
		}
	}

	// 4. FCM 메시지 일괄 전송
	if result := sendBatchFCMMessages(messages); result {
		// 5. 새 알림을 DB에 일괄 저장
		if len(newNotifications) > 0 {
			if err := db.Create(&newNotifications).Error; err != nil {
				log.Printf("error creating notifications: %v\n", err)
			}
		}
	}

}

// 사용자별 읽지 않은 알림 수 조회
func getUnreadNotificationCounts(db *gorm.DB) map[uint]uint {
	var results []struct {
		Uid   uint
		Count uint
	}
	db.Model(&model.Notification{}).
		Select("uid, COUNT(*) as count").
		Where("is_read = ?", false).
		Group("uid").
		Scan(&results)

	counts := make(map[uint]uint)
	for _, result := range results {
		counts[result.Uid] = result.Count
	}
	return counts
}

// FCM 메시지를 일괄 전송
func sendBatchFCMMessages(messages []*messaging.Message) bool {
	if len(messages) == 0 {
		return false
	}
	ctx := context.Background()
	br, err := firebaseClient.SendAll(ctx, messages)
	if err != nil {
		log.Printf("error sending batch messages: %v\n", err)
		return false
	} else {
		log.Printf("Successfully sent %d messages\n", br.SuccessCount)
		return true
	}
}

// 알람 유형에 따른 제목 생성
func getNotificationTitle(notificationType uint) string {
	switch notificationType {
	case uint(util.ExerciseType):
		return "운동 시간"
	case uint(util.MedicineType):
		return "약물 복용"
	case uint(util.SleepType):
		return "수면 시간"
	default:
		return "알림"
	}
}

// 알람이 전송될 조건인지 확인
func shouldSendNotification(now time.Time, alarm model.Alarm) bool {
	currentWeekday := int(now.Weekday())

	var alarmWeekdays []int
	if err := json.Unmarshal(alarm.Week, &alarmWeekdays); err != nil {
		return false
	}

	alarmTime, _ := time.Parse("15:04", alarm.Timestamp)

	for _, weekday := range alarmWeekdays {
		if weekday == currentWeekday && now.Format("15:04") == alarmTime.Format("15:04") {
			return true
		}
	}
	return false
}
