// /exercise-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"context"
	"encoding/json"
	"errors"
	"exercise-service/dto"
	pb "exercise-service/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type ExerciseService interface {
	SaveExercise(presetRequest dto.ExerciseRequest) (string, error)
	GetExercises(id int, startDate, endDate string) ([]dto.ExerciseDateInfo, error)
	RemoveExercises(ids []int, uid int) (string, error)
	DoExercise(exerciseDo dto.ExerciseDo) (string, error)
}

type exerciseService struct {
	db          *gorm.DB
	alarmClient pb.AlarmServiceClient
}

func NewExerciseService(db *gorm.DB, conn *grpc.ClientConn) ExerciseService {
	alarmClient := pb.NewAlarmServiceClient(conn)
	return &exerciseService{db: db, alarmClient: alarmClient}
}

func (service *exerciseService) GetExercises(id int, startDateStr, endDateStr string) ([]dto.ExerciseDateInfo, error) {

	// 문자열을 time.Time 타입으로 변환
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)

	if err != nil {
		return nil, err
	}

	var exercises []model.Exercise
	var exerciseResponse []dto.ExerciseResponse
	err = service.db.Where("uid = ? AND plan_start_at <= ? AND plan_end_at >= ?", id, endDate, startDate).Find(&exercises).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(exercises, &exerciseResponse); err != nil {
		return nil, err
	}

	var exerciseIDs []int
	for _, exercise := range exerciseResponse {
		exerciseIDs = append(exerciseIDs, exercise.Id)
	}

	var performedExercises []model.ExerciseInfo
	err = service.db.Where("exercise_id IN (?) AND date_performed BETWEEN ? AND ?", exerciseIDs, startDate, endDate).Find(&performedExercises).Error
	if err != nil {
		return nil, err
	}

	performedMap := make(map[int]map[string]bool)
	for _, pe := range performedExercises {
		if performedMap[pe.ExerciseId] == nil {
			performedMap[pe.ExerciseId] = make(map[string]bool)
		}
		performedMap[pe.ExerciseId][pe.DatePerformed] = true
	}

	// var exerciseDates []dto.ExerciseDateInfo
	// log.Println(exerciseDates) // 출력: [] 이지만 실제로 nil임
	exerciseDates := make([]dto.ExerciseDateInfo, 0)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		var dailyExercises []dto.ExerciseDoneInfo
		for _, e := range exerciseResponse {
			planStartAt, _ := time.Parse("2006-01-02", e.PlanStartAt)
			planEndAt, _ := time.Parse("2006-01-02", e.PlanEndAt)

			if d.After(planStartAt) && d.Before(planEndAt) && isExerciseDay(e.Weekdays, d.Weekday()) {
				performed := performedMap[e.Id][d.Format("2006-01-02")]
				dailyExercises = append(dailyExercises, dto.ExerciseDoneInfo{Exercise: e, Done: performed})
			}
		}
		if len(dailyExercises) > 0 {
			exerciseDates = append(exerciseDates, dto.ExerciseDateInfo{Date: d.Format("2006-01-02"), Exercises: dailyExercises})
		}
	}
	log.Println(exerciseDates)
	return exerciseDates, nil
}

func isExerciseDay(weekdays []int, day time.Weekday) bool {
	for _, d := range weekdays {
		if int(day) == d {
			return true
		}
	}
	return false
}

func (service *exerciseService) DoExercise(exerciseDo dto.ExerciseDo) (string, error) {
	var info model.ExerciseInfo
	datePerformed, err := time.Parse("2006-01-02", exerciseDo.PerformedDate)
	if err != nil {
		return "", errors.New("must YYYY-MM-DD")
	}
	result := service.db.Where("exercise_id = ? AND uid=? AND date_performed = ?", exerciseDo.ExerciseId, exerciseDo.Uid, datePerformed).First(&info)

	if result.RowsAffected == 0 {
		// 기록이 없으면 새로운 기록 추가
		newInfo := model.ExerciseInfo{
			Uid:           exerciseDo.Uid,
			ExerciseId:    exerciseDo.ExerciseId,
			DatePerformed: exerciseDo.PerformedDate,
		}
		if err := service.db.Create(&newInfo).Error; err != nil {
			return "", errors.New("db error")
		}
	} else {
		result := service.db.Where("exercise_id = ? AND uid=? AND date_performed = ?", exerciseDo.ExerciseId, exerciseDo.Uid, exerciseDo.PerformedDate).Delete(&model.ExerciseInfo{})
		if result.Error != nil {
			return "", errors.New("db error2")
		}
	}
	return "200", nil
}

func (service *exerciseService) SaveExercise(exerciseRequest dto.ExerciseRequest) (string, error) {

	if err := validateExercise(exerciseRequest); err != nil {
		return "", err
	}
	var exercise model.Exercise

	result := service.db.Where("id=? AND uid=?", exerciseRequest.Id, exerciseRequest.Uid).First(&model.Exercise{})

	if err := util.CopyStruct(exerciseRequest, &exercise); err != nil {
		return "", err
	}

	exercise.Uid = exerciseRequest.Uid //  json: "-" 이라서

	// JSON 배열을 Go 슬라이스로 변환
	var weekdaySlice []int32
	err := json.Unmarshal(exercise.Weekdays, &weekdaySlice)
	if err != nil {
		return "", err
	}

	seen := make(map[int32]bool)
	unique := []int32{}

	for _, v := range weekdaySlice {
		// 숫자가 0과 6 사이인지 확인
		if v >= 0 && v <= 6 {
			// 중복되지 않은 경우, 결과 슬라이스에 추가
			if !seen[v] {
				seen[v] = true
				unique = append(unique, v)
			}
		}
	}

	newWeekdays, err := json.Marshal(unique)
	if err != nil {
		return "", err
	}
	exercise.Weekdays = newWeekdays

	ar := &pb.AlarmRequest{
		Id:        int32(exercise.Id),
		Uid:       int32(exercise.Uid),
		Body:      "운동 할 시간입니다.",
		Type:      "exercise",
		StartAt:   exercise.PlanStartAt,
		EndAt:     exercise.PlanEndAt,
		Timestamp: exercise.ExerciseStartAt,
		Week:      unique,
		UseAlarm:  exercise.UseAlarm,
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		exercise.Id = 0

		if err := service.db.Create(&exercise).Error; err != nil {
			return "", err
		}

	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&exercise).Updates(exercise).Error; err != nil {
			return "", err
		}
	}
	if exercise.UseAlarm {
		sendAlarm(service, ar)
	} else {
		b := make([]int32, 1)

		b[0] = int32(exercise.Id)

		arr := &pb.AlarmRemoveRequest{
			Ids: b,
			Uid: int32(exercise.Uid),
		}
		removeAlarm(service, arr)
	}

	return "200", nil
}

func (service *exerciseService) RemoveExercises(ids []int, uid int) (string, error) {
	result := service.db.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.Exercise{})

	if result.Error != nil {
		return "", errors.New("db error")
	}

	b := make([]int32, len(ids))

	for i, v := range ids {
		b[i] = int32(v)
	}
	arr := &pb.AlarmRemoveRequest{
		Ids: b,
		Uid: int32(uid),
	}

	removeAlarm(service, arr)
	return "200", nil
}

func sendAlarm(service *exerciseService, ar *pb.AlarmRequest) {
	reponse, err := service.alarmClient.SetAlarm(context.Background(), ar)
	if err != nil {
		log.Printf("Failed to set Alarm: %v", err)
	}
	log.Printf("set Alarm: %v", reponse)
}

func removeAlarm(service *exerciseService, arr *pb.AlarmRemoveRequest) {
	reponse, err := service.alarmClient.RemoveAlarm(context.Background(), arr)
	if err != nil {
		log.Printf("Failed to remove Alarm: %v", err)
	}
	log.Printf("remove Alarm: %v", reponse)
}
