// /exercise-service/service/service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"exercise-service/common/model"
	"exercise-service/common/util"
	"exercise-service/dto"
	pb "exercise-service/proto"
	"log"
	"reflect"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type ExerciseService interface {
	SaveExercise(presetRequest dto.ExerciseRequest) (string, error)
	GetExercises(id uint, startDate, endDate string) ([]dto.ExerciseDateInfo, error)
	RemoveExercises(ids []uint, uid uint) (string, error)
	DoExercise(exerciseDo dto.ExerciseDo) (string, error)
	GetProjects() ([]dto.ProjectResponse, error)
	GetVideos(projectId string, page uint) ([]dto.VideoResponse, error)
}

type exerciseService struct {
	db          *gorm.DB
	alarmClient pb.AlarmServiceClient
}

func NewExerciseService(db *gorm.DB, conn *grpc.ClientConn) ExerciseService {
	alarmClient := pb.NewAlarmServiceClient(conn)
	return &exerciseService{db: db, alarmClient: alarmClient}
}

func (service *exerciseService) GetExercises(id uint, startDateStr, endDateStr string) ([]dto.ExerciseDateInfo, error) {

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
	err = service.db.Debug().Where("uid = ? AND plan_start_at <= ? AND plan_end_at >= ?", id, endDate.Format("2006-01-02"), startDate.Format("2006-01-02")).Find(&exercises).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(exercises, &exerciseResponse); err != nil {
		return nil, err
	}

	// 해당 운동 실행내역 조회
	var exerciseIDs []uint
	for _, exercise := range exerciseResponse {
		exerciseIDs = append(exerciseIDs, exercise.Id)
	}
	var performedExercises []model.ExerciseInfo
	err = service.db.Where("exercise_id IN (?) AND date_performed BETWEEN ? AND ?", exerciseIDs, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&performedExercises).Error
	if err != nil {
		return nil, err
	}

	// 운동 실행내역 응답형식으로 가공
	performedMap := make(map[uint]map[string]bool)
	for _, pe := range performedExercises {
		if performedMap[pe.ExerciseId] == nil {
			performedMap[pe.ExerciseId] = make(map[string]bool)
		}
		performedMap[pe.ExerciseId][pe.DatePerformed] = true
	}

	// var exerciseDates []dto.ExerciseDateInfo
	// log.Println(exerciseDates) // 출력: [] 이지만 실제로 nil임

	// 전체 날짜에서 실행한 날짜 체크
	exerciseDates := make([]dto.ExerciseDateInfo, 0)
	repeatMap := make(map[uint]uint)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		var dailyExercises []dto.ExerciseDoneInfo
		for _, e := range exerciseResponse {

			planStartAt, _ := time.Parse("2006-01-02", e.PlanStartAt)
			planEndAt, _ := time.Parse("2006-01-02", e.PlanEndAt)

			if !d.Before(planStartAt) && d.Before(planEndAt.AddDate(0, 0, 1)) && isExerciseDay(e.Weekdays, d.Weekday()) {
				repeatMap[e.Id] += 1
				e.Repeat = repeatMap[e.Id]
				performed := performedMap[e.Id][d.Format("2006-01-02")]
				dailyExercises = append(dailyExercises, dto.ExerciseDoneInfo{Exercise: e, Done: performed})
			}
		}
		if len(dailyExercises) > 0 {
			exerciseDates = append(exerciseDates, dto.ExerciseDateInfo{Date: d.Format("2006-01-02"), Exercises: dailyExercises})
		}
	}
	return exerciseDates, nil
}

func isExerciseDay(weekdays []uint, day time.Weekday) bool {
	for _, d := range weekdays {
		if uint(day) == d {
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
	result := service.db.Debug().Where("exercise_id = ? AND uid=? AND date_performed = ?", exerciseDo.ExerciseId, exerciseDo.Uid, datePerformed.Format("2006-01-02")).First(&info)

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
		result := service.db.Debug().Where("exercise_id = ? AND uid=? AND date_performed = ?", exerciseDo.ExerciseId, exerciseDo.Uid, exerciseDo.PerformedDate).Delete(&model.ExerciseInfo{})
		if result.Error != nil {
			return "", errors.New("db error2")
		}
	}
	return "200", nil
}

func (service *exerciseService) SaveExercise(exerciseRequest dto.ExerciseRequest) (string, error) {

	if exerciseRequest.PlanEndAt == "" {
		exerciseRequest.PlanEndAt = "2099-01-01"
	}

	if err := validateExercise(exerciseRequest); err != nil {
		return "", err
	}

	var exercise model.Exercise

	result := service.db.Where("id=? AND uid=?", exerciseRequest.Id, exerciseRequest.Uid).First(&model.Exercise{})

	if err := util.CopyStruct(exerciseRequest, &exercise); err != nil {
		return "", err
	}

	exercise.Uid = exerciseRequest.Uid //  json: "-" 이라서

	newWeekdays, unique, err := validateWeek(exercise.Weekdays)
	if err != nil {
		return "", err
	}
	exercise.Weekdays = newWeekdays

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		exercise.Id = 0

		if err := service.db.Create(&exercise).Error; err != nil {
			return "", err
		}
		ar := &pb.AlarmRequest{
			ParentId:  int32(exercise.Id),
			Uid:       int32(exercise.Uid),
			Body:      "운동 할 시간입니다.",
			Type:      int32(util.ExerciseType),
			StartAt:   exercise.PlanStartAt,
			EndAt:     exercise.PlanEndAt,
			Timestamp: exercise.ExerciseStartAt,
			Week:      unique,
		}
		if exercise.UseAlarm {
			go sendAlarm(service, ar)
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		updateFields := make(map[string]interface{})

		userRequestValue := reflect.ValueOf(exerciseRequest)
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

		if err := service.db.Model(&exercise).Updates(updateFields).Error; err != nil {
			return "", err
		}
		ar := &pb.AlarmRequest{
			ParentId:  int32(exercise.Id),
			Uid:       int32(exercise.Uid),
			Body:      "운동 할 시간입니다.",
			Type:      int32(util.ExerciseType),
			StartAt:   exercise.PlanStartAt,
			EndAt:     exercise.PlanEndAt,
			Timestamp: exercise.ExerciseStartAt,
			Week:      unique,
		}
		if exercise.UseAlarm {
			go updateAlarm(service, ar)
		} else {
			b := make([]int32, 1)

			b[0] = int32(exercise.Id)

			arr := &pb.AlarmRemoveRequest{
				ParentIds: b,
				Uid:       int32(exercise.Uid),
			}
			go removeAlarm(service, arr)
		}
	}

	return "200", nil
}

func (service *exerciseService) RemoveExercises(ids []uint, uid uint) (string, error) {
	result := service.db.Model(&model.Exercise{}).Where("id IN (?) AND uid= ?", ids, uid).Select("is_delete").Updates(map[string]interface{}{"is_delete": true})

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
		Type:      int32(util.ExerciseType),
	}

	go removeAlarm(service, arr)
	return "200", nil
}

func (service *exerciseService) GetProjects() ([]dto.ProjectResponse, error) {

	var projects []dto.ProjectResponse
	err := service.db.Model(&model.Video{}).
		Select("project_id, project_name as name, count(*) as count").
		Group("project_id").Scan(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil
}

// face-service 의 face-exercise 쪽에는 한번에 다가져옴

func (service *exerciseService) GetVideos(projectId string, page uint) ([]dto.VideoResponse, error) {
	pageSize := 20
	var videos []model.Video
	offset := page * uint(pageSize)

	query := service.db.Where("project_id = ?", projectId)

	query = query.Order("id DESC")
	result := query.Offset(int(offset)).Limit(pageSize).Find(&videos)

	if result.Error != nil {
		return nil, result.Error
	}

	var VideoResponses []dto.VideoResponse
	if err := util.CopyStruct(videos, &VideoResponses); err != nil {
		return nil, err
	}

	return VideoResponses, nil
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

func updateAlarm(service *exerciseService, ar *pb.AlarmRequest) {
	reponse, err := service.alarmClient.UpdateAlarm(context.Background(), ar)
	if err != nil {
		log.Printf("Failed to update Alarm: %v", err)
	}
	log.Printf("update Alarm: %v", reponse)
}
