// /fcm-service/service/schedular-service.go

package service

import (
	"common/model"
	"context"
	"log"
	"strconv"
	"strings"
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
	db.Where("start_at <= ? AND end_at >= ?", now.Format("2006-01-02"), now.Format("2006-01-02")).Find(&alarms)

	for _, alarm := range alarms {
		if shouldSendNotification(now, alarm) {
			go sendMedicationReminder(context.Background(), alarm, db)
		}
	}
}

func shouldSendNotification(now time.Time, alarm model.Alarm) bool {
	currentWeekday := int(now.Weekday())
	alarmWeekdays := strings.Split(alarm.Week, ",")
	alarmTime, _ := time.Parse("15:04", alarm.Timestamp)

	for _, w := range alarmWeekdays {
		weekday, _ := strconv.Atoi(w)
		if weekday == currentWeekday && now.Format("15:04") == alarmTime.Format("15:04") {
			return true
		}
	}
	return false
}

func sendMedicationReminder(ctx context.Context, alarm model.Alarm, db *gorm.DB) {
	var user model.User
	if err := db.First(&user, "id = ?", alarm.Uid).Error; err != nil {
		log.Printf("Failed to find user: %v\n", err)
		return
	}

	notification_count := calculateNotificationCount(db, alarm.Uid)

	message := &messaging.Message{
		Data: map[string]string{
			"startAt":            alarm.StartAt,
			"endAt":              alarm.EndAt,
			"uid":                strconv.Itoa(alarm.Uid),
			"type":               alarm.Type,
			"notification_count": strconv.Itoa(notification_count),
		},
		Notification: &messaging.Notification{
			Title: "알림 제목",
			Body:  alarm.Body,
		},
		Token: user.FCMToken,
	}

	response, err := firebaseClient.Send(ctx, message)
	if err != nil {
		log.Printf("error sending message: %v\n", err)
		return
	}
	log.Printf("Successfully sent message: %s\n", response)

	newNotification := model.Notification{
		Uid:    alarm.Uid,
		Type:   alarm.Type,
		Body:   alarm.Body,
		IsRead: false,
	}
	db.Create(&newNotification)

	var notifications []model.Notification
	db.Where("uid = ?", alarm.Uid).Find(&notifications)
	if len(notifications) > 100 {
		db.Delete(&notifications[0])
	}
}

func calculateNotificationCount(db *gorm.DB, uid int) int {
	var count int64
	db.Model(&model.Notification{}).Where("uid = ? AND is_read = ?", uid, false).Count(&count)
	return int(count)
}
