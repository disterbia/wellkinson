// /fcm-service/service/schedular-service.go

package service

import (
	"context"
	"encoding/json"
	"fcm-service/common/model"
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
	db.Where("(start_at ='' OR start_at <= ?) AND (end_at ='' OR end_at >= ?)", now.Format("2006-01-02"), now.Format("2006-01-02")).Find(&alarms)

	for _, alarm := range alarms {
		if shouldSendNotification(now, alarm) {
			go sendMedicationReminder(context.Background(), alarm, db)
		}
	}
}

func shouldSendNotification(now time.Time, alarm model.Alarm) bool {
	// sunday = 0, ...
	currentWeekday := int(now.Weekday())

	// alarm.Week를 []int로 언마샬링
	var alarmWeekdaysInts []int
	err := json.Unmarshal(alarm.Week, &alarmWeekdaysInts)
	if err != nil {
		// 언마샬링 에러 처리
		return false
	}

	// []int를 []string으로 변환
	alarmWeekdays := make([]string, len(alarmWeekdaysInts))
	for i, w := range alarmWeekdaysInts {
		alarmWeekdays[i] = strconv.Itoa(w)
	}

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
	log.Println(user.FCMToken)
	message := &messaging.Message{
		Data: map[string]string{
			"start_at":           alarm.StartAt,
			"end_at":             alarm.EndAt,
			"uid":                strconv.FormatUint(uint64(alarm.Uid), 10),
			"type":               strconv.FormatUint(uint64(alarm.Type), 10),
			"notification_count": strconv.FormatUint(uint64(notification_count), 10),
			"timestamp":          strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
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
		Uid:      alarm.Uid,
		Type:     alarm.Type,
		Body:     alarm.Body,
		ParentId: alarm.ParentId,
		IsRead:   false,
	}
	db.Create(&newNotification)

	var notifications []model.Notification
	db.Where("uid = ?", alarm.Uid).Find(&notifications)
	if len(notifications) > 100 {
		db.Delete(&notifications[0])
	}
}

func calculateNotificationCount(db *gorm.DB, uid uint) uint {
	var count int64
	db.Model(&model.Notification{}).Where("uid = ? AND is_read = ?", uid, false).Count(&count)
	return uint(count)
}
