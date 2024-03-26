// /sleep-service/service/service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"sleep-service/common/model"
	"sleep-service/common/util"
	"sleep-service/dto"
	pb "sleep-service/proto"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type SleepService interface {
	SaveSleepAlarm(sleepRequest dto.SleepAlarmRequest) (string, error)
	GetSleepAlarms(id uint) ([]dto.SleepAlarmResponse, error)
	RemoveSleepAlarms(ids []uint, uid uint) (string, error)
	GetSleepTimes(id uint, startDate, endDate string) ([]dto.SleepTimeResponse, error)
	SaveSleepTime(sleepRequest dto.SleepTimeRequest) (string, error)
	RemoveSleepTime(id uint, uid uint) (string, error)
}

type sleepService struct {
	db          *gorm.DB
	alarmClient pb.AlarmServiceClient
}

func NewSleepService(db *gorm.DB, conn *grpc.ClientConn) SleepService {
	alarmClient := pb.NewAlarmServiceClient(conn)
	return &sleepService{db: db, alarmClient: alarmClient}
}

func (service *sleepService) SaveSleepAlarm(sleepRequest dto.SleepAlarmRequest) (string, error) {

	if err := validateSleep(sleepRequest); err != nil {
		return "", err
	}
	var weekdays []json.RawMessage
	service.db.Model(&model.SleepAlarm{}).Where("id != ? AND uid=?", sleepRequest.Id, sleepRequest.Uid).Pluck("weekdays", &weekdays)

	if err := isDuplicateWeek(sleepRequest.Weekdays, weekdays); err != nil {
		return "", err
	}

	var sleep model.SleepAlarm

	result := service.db.Where("id=? AND uid=?", sleepRequest.Id, sleepRequest.Uid).First(&model.SleepAlarm{})

	if err := util.CopyStruct(sleepRequest, &sleep); err != nil {
		return "", err
	}

	sleep.Uid = sleepRequest.Uid //  json: "-" 이라서

	newWeekdays, unique, err2 := validateWeek(sleep.Weekdays)
	if err2 != nil {
		return "", errors.New("week error")
	}
	sleep.Weekdays = newWeekdays

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		sleep.Id = 0

		if err := service.db.Create(&sleep).Error; err != nil {
			return "", err
		}
		ar := &pb.AlarmRequest{
			ParentId: int32(sleep.Id),
			Uid:      int32(sleep.Uid),
			Body:     "취침 할 시간입니다.",
			Type:     int32(util.SleepType),
			StartAt:  time.Now().Format("2006-01-02"),
			// EndAt:     "",
			Timestamp: sleep.AlarmTime,
			Week:      unique,
		}
		if sleep.IsActive {
			go sendAlarm(service, ar)
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트

		updateFields := make(map[string]interface{})

		userRequestValue := reflect.ValueOf(sleepRequest)
		userRequestType := userRequestValue.Type()
		for i := 0; i < userRequestValue.NumField(); i++ {
			field := userRequestValue.Field(i)
			fieldName := userRequestType.Field(i).Tag.Get("json")
			if fieldName == "-" {
				continue
			}
			if !field.IsZero() {
				if fieldName == "weekdays" {
					// weekdays 필드를 JSON 형식으로 변환
					updateFields[fieldName], _ = json.Marshal(field.Interface())
				} else {
					updateFields[fieldName] = field.Interface()
				}
			}
		}
		if err := service.db.Model(&sleep).Debug().Updates(updateFields).Error; err != nil {
			return "", err
		}
		ar := &pb.AlarmRequest{
			ParentId: int32(sleep.Id),
			Uid:      int32(sleep.Uid),
			Body:     "취침 할 시간입니다.",
			Type:     int32(util.SleepType),
			StartAt:  time.Now().Format("2006-01-02"),
			// EndAt:     "",
			Timestamp: sleep.AlarmTime,
			Week:      unique,
		}
		if sleep.IsActive {
			go updateAlarm(service, ar)
		} else {
			b := make([]int32, 1)

			b[0] = int32(sleep.Id)

			arr := &pb.AlarmRemoveRequest{
				ParentIds: b,
				Uid:       int32(sleep.Uid),
				Type:      int32(util.SleepType),
			}
			go removeAlarm(service, arr)
		}
	}

	return "200", nil
}

func (service *sleepService) GetSleepAlarms(id uint) ([]dto.SleepAlarmResponse, error) {
	var sleepAlarms []model.SleepAlarm
	var alarmsResponses []dto.SleepAlarmResponse
	err := service.db.Where("uid = ? ", id).Order("end_time").Find(&sleepAlarms).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(sleepAlarms, &alarmsResponses); err != nil {
		return nil, err
	}
	return alarmsResponses, nil
}

func (service *sleepService) RemoveSleepAlarms(ids []uint, uid uint) (string, error) {
	result := service.db.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.SleepAlarm{})

	if result.Error != nil {
		return "", errors.New("db error")
	}

	b := make([]int32, len(ids))

	for i, v := range ids {
		b[i] = int32(v)
	}

	arr := &pb.AlarmRemoveRequest{
		ParentIds: b,
		Uid:       int32(uid),
		Type:      int32(util.SleepType),
	}

	go removeAlarm(service, arr)
	return "200", nil
}

func (service *sleepService) GetSleepTimes(id uint, startDateStr, endDateStr string) ([]dto.SleepTimeResponse, error) {

	// 문자열을 time.Time 타입으로 변환
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)

	if err != nil {
		return nil, err
	}

	var sleepTimes []model.SleepTime
	var sleepTimeResponses []dto.SleepTimeResponse
	err = service.db.Where("date_sleep BETWEEN ? AND ?", startDate, endDate).Find(&sleepTimes).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(sleepTimes, &sleepTimeResponses); err != nil {
		return nil, err
	}

	return sleepTimeResponses, nil
}

func (service *sleepService) RemoveSleepTime(id uint, uid uint) (string, error) {
	result := service.db.Where("id = ? AND uid= ?", id, uid).Delete(&model.SleepTime{})
	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *sleepService) SaveSleepTime(sleepRequest dto.SleepTimeRequest) (string, error) {
	if err := validateSleepTime(sleepRequest); err != nil {
		return "", err
	}

	var sleepTime model.SleepTime

	result := service.db.Where("date_sleep = ? AND uid=?", sleepRequest.DateSleep, sleepRequest.Uid).First(&sleepTime)

	if err := util.CopyStruct(sleepRequest, &sleepTime); err != nil {
		return "", err
	}

	sleepTime.Uid = sleepRequest.Uid
	if result.RowsAffected == 0 {

		if err := service.db.Create(&sleepTime).Error; err != nil {
			return "", errors.New("db error")
		}
	} else {
		result := service.db.Model(&sleepTime).Updates(sleepRequest)
		if result.Error != nil {
			return "", errors.New("db error2")
		}
	}
	return "200", nil
}

func sendAlarm(service *sleepService, ar *pb.AlarmRequest) {
	reponse, err := service.alarmClient.SetAlarm(context.Background(), ar)
	if err != nil {
		log.Printf("Failed to set Alarm: %v", err)
	}
	log.Printf("set Alarm: %v", reponse)
}

func removeAlarm(service *sleepService, arr *pb.AlarmRemoveRequest) {
	reponse, err := service.alarmClient.RemoveAlarm(context.Background(), arr)
	if err != nil {
		log.Printf("Failed to remove Alarm: %v", err)
	}
	log.Printf("remove Alarm: %v", reponse)
}

func updateAlarm(service *sleepService, ar *pb.AlarmRequest) {
	reponse, err := service.alarmClient.UpdateAlarm(context.Background(), ar)
	if err != nil {
		log.Printf("Failed to update Alarm: %v", err)
	}
	log.Printf("update Alarm: %v", reponse)
}
