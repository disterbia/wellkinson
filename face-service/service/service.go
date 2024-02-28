// /face-service/service/service.go
package service

import (
	"errors"
	"face-service/common/model"
	"face-service/common/util"
	"face-service/dto"

	"gorm.io/gorm"
)

type FaceService interface {
	SaveFaceScores(faceScoreRequests []dto.FaceScoreRequest) (string, error)
	GetFaceScores(id uint, startDate, endDate string) ([]dto.FaceScoreResponse, error)
	GetFaceExams() ([]dto.FaceExamResponse, error)
	GetFaceExercises() ([]dto.FaceExerciseResponse, error)
}

type faceService struct {
	db *gorm.DB
}

func NewFaceService(db *gorm.DB) FaceService {
	return &faceService{db: db}
}
func (service *faceService) GetFaceExams() ([]dto.FaceExamResponse, error) {
	var faceExams []model.FaceExam
	var faceExamResponses []dto.FaceExamResponse

	err := service.db.Find(&faceExams).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(faceExams, &faceExamResponses); err != nil {
		return nil, err
	}

	return faceExamResponses, nil
}

// exercise-servie의 video는 따로 가져옴 둘중 좋은방법 선택하기

func (service *faceService) GetFaceExercises() ([]dto.FaceExerciseResponse, error) {

	var faceExercises []model.FaceExercise
	faceExercisResponses := make([]dto.FaceExerciseResponse, 0)

	err := service.db.Model(&model.FaceExercise{}).
		Select("type, " +
			"CASE " +
			"WHEN type = 1 THEN '기쁨' " +
			"WHEN type = 2 THEN '슬픔' " +
			"WHEN type = 3 THEN '놀람' " +
			"WHEN type = 4 THEN '분노' " +
			"ELSE '기타' END as title, " +
			"count(*) as count").
		Group("type").Scan(&faceExercisResponses).Error

	if err != nil {
		return nil, err
	}

	err = service.db.Find(&faceExercises).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	faceMap := make(map[uint][]model.FaceExercise)
	for _, faceExercise := range faceExercises {
		faceMap[faceExercise.Type] = append(faceMap[faceExercise.Type], faceExercise)
	}

	for i := range faceExercisResponses {
		faceExercisResponses[i].FaceExercise = faceMap[faceExercisResponses[i].Type]
	}

	return faceExercisResponses, nil
}

func (service *faceService) GetFaceScores(id uint, startDate, endDate string) ([]dto.FaceScoreResponse, error) {
	var faceScores []model.FaceScore
	var faceScoreResponses []dto.FaceScoreResponse

	query := service.db.Where("uid = ?", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")

	err := query.Find(&faceScores).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(faceScores, &faceScoreResponses); err != nil {
		return nil, err
	}

	return faceScoreResponses, nil
}

func (service *faceService) SaveFaceScores(faceScoreRequests []dto.FaceScoreRequest) (string, error) {

	var faceScores []model.FaceScore

	if err := util.CopyStruct(faceScoreRequests, &faceScores); err != nil {
		return "", err
	}
	var len = len(faceScoreRequests)
	for i := 0; i < len; i++ {
		faceScores[i].Uid = faceScoreRequests[0].Uid
	}

	if err := service.db.Create(&faceScores).Error; err != nil {
		return "", err
	}

	return "200", nil
}
