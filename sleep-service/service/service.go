// /sleep-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"context"
	"encoding/json"
	"errors"
	"log"
	"sleep-service/dto"
	pb "sleep-service/proto"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type SleepService interface {
	SaveSleepAlarm(sleepRequest dto.SleepAlarmRequest) (string, error)
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

	var sleep model.SleepAlarm

	result := service.db.Where("id=? AND uid=?", sleepRequest.Id, sleepRequest.Uid).First(&model.SleepAlarm{})

	var weekdays []json.RawMessage
	err := service.db.Where("uid=?", sleepRequest.Id, sleepRequest.Uid).Pluck("weekdays", &weekdays)
	if err != nil {
		return "", errors.New("db error")
	}

	if err := isDuplicateWeek(sleepRequest.Weekdays, weekdays); err != nil {
		return "", err
	}

	if err := util.CopyStruct(sleepRequest, &sleep); err != nil {
		return "", err
	}

	// ////

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
		if err := service.db.Model(&sleep).Updates(sleep).Error; err != nil {
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
			}
			go removeAlarm(service, arr)
		}
	}

	return "200", nil
}

// func (service *exerciseService) GetExercises(id uint, startDateStr, endDateStr string) ([]dto.ExerciseDateInfo, error) {

// 	// 문자열을 time.Time 타입으로 변환
// 	startDate, err := time.Parse("2006-01-02", startDateStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	endDate, err := time.Parse("2006-01-02", endDateStr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	var exercises []model.Exercise
// 	var exerciseResponse []dto.ExerciseResponse
// 	err = service.db.Where("uid = ? AND plan_start_at <= ? AND plan_end_at >= ?", id, endDate, startDate).Find(&exercises).Error
// 	if err != nil {
// 		return nil, errors.New("db error")
// 	}

// 	if err := util.CopyStruct(exercises, &exerciseResponse); err != nil {
// 		return nil, err
// 	}

// 	// 해당 운동 실행내역 조회
// 	var exerciseIDs []uint
// 	for _, exercise := range exerciseResponse {
// 		exerciseIDs = append(exerciseIDs, exercise.Id)
// 	}
// 	var performedExercises []model.ExerciseInfo
// 	err = service.db.Where("exercise_id IN (?) AND date_performed BETWEEN ? AND ?", exerciseIDs, startDate, endDate).Find(&performedExercises).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 운동 실행내역 응답형식으로 가공
// 	performedMap := make(map[uint]map[string]bool)
// 	for _, pe := range performedExercises {
// 		if performedMap[pe.ExerciseId] == nil {
// 			performedMap[pe.ExerciseId] = make(map[string]bool)
// 		}
// 		performedMap[pe.ExerciseId][pe.DatePerformed] = true
// 	}

// 	// var exerciseDates []dto.ExerciseDateInfo
// 	// log.Println(exerciseDates) // 출력: [] 이지만 실제로 nil임

// 	// 전체 날짜에서 실행한 날짜 체크
// 	exerciseDates := make([]dto.ExerciseDateInfo, 0)
// 	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
// 		var dailyExercises []dto.ExerciseDoneInfo
// 		for _, e := range exerciseResponse {
// 			planStartAt, _ := time.Parse("2006-01-02", e.PlanStartAt)
// 			planEndAt, _ := time.Parse("2006-01-02", e.PlanEndAt)

// 			if d.After(planStartAt) && d.Before(planEndAt) && isExerciseDay(e.Weekdays, d.Weekday()) {
// 				performed := performedMap[e.Id][d.Format("2006-01-02")]
// 				dailyExercises = append(dailyExercises, dto.ExerciseDoneInfo{Exercise: e, Done: performed})
// 			}
// 		}
// 		if len(dailyExercises) > 0 {
// 			exerciseDates = append(exerciseDates, dto.ExerciseDateInfo{Date: d.Format("2006-01-02"), Exercises: dailyExercises})
// 		}
// 	}
// 	log.Println(exerciseDates)
// 	return exerciseDates, nil
// }

// func isExerciseDay(weekdays []uint, day time.Weekday) bool {
// 	for _, d := range weekdays {
// 		if uint(day) == d {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (service *exerciseService) DoExercise(exerciseDo dto.ExerciseDo) (string, error) {
// 	var info model.ExerciseInfo
// 	datePerformed, err := time.Parse("2006-01-02", exerciseDo.PerformedDate)
// 	if err != nil {
// 		return "", errors.New("must YYYY-MM-DD")
// 	}
// 	result := service.db.Where("exercise_id = ? AND uid=? AND date_performed = ?", exerciseDo.ExerciseId, exerciseDo.Uid, datePerformed).First(&info)

// 	if result.RowsAffected == 0 {
// 		// 기록이 없으면 새로운 기록 추가
// 		newInfo := model.ExerciseInfo{
// 			Uid:           exerciseDo.Uid,
// 			ExerciseId:    exerciseDo.ExerciseId,
// 			DatePerformed: exerciseDo.PerformedDate,
// 		}
// 		if err := service.db.Create(&newInfo).Error; err != nil {
// 			return "", errors.New("db error")
// 		}
// 	} else {
// 		result := service.db.Where("exercise_id = ? AND uid=? AND date_performed = ?", exerciseDo.ExerciseId, exerciseDo.Uid, exerciseDo.PerformedDate).Delete(&model.ExerciseInfo{})
// 		if result.Error != nil {
// 			return "", errors.New("db error2")
// 		}
// 	}
// 	return "200", nil
// }

// func (service *exerciseService) RemoveExercises(ids []uint, uid uint) (string, error) {
// 	result := service.db.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.Exercise{})

// 	if result.Error != nil {
// 		return "", errors.New("db error")
// 	}

// 	b := make([]int32, len(ids))

// 	for i, v := range ids {
// 		b[i] = int32(v)
// 	}
// 	arr := &pb.AlarmRemoveRequest{
// 		ParentIds: b,
// 		Uid:       int32(uid),
// 		Type:      int32(util.ExerciseType),
// 	}

// 	go removeAlarm(service, arr)
// 	return "200", nil
// }

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
