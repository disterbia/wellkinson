// /face-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"errors"
	"face-service/dto"

	"gorm.io/gorm"
)

type FaceService interface {
	SaveFaceScores(faceScoreRequests []dto.FaceScoreRequest) (string, error)
	GetFaceScores(id int, startDate, endDate string) ([]dto.FaceScoreResponse, error)
	GetFaceExams() ([]dto.FaceExamResponse, error)
}

type faceService struct {
	db *gorm.DB
}

func NewFaceService(db *gorm.DB) FaceService {
	return &faceService{db: db}
}
func (service *faceService) GetFaceExams() ([]dto.FaceExamResponse, error) {
	var faceExams []model.FaceExams
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

func (service *faceService) GetFaceScores(id int, startDate, endDate string) ([]dto.FaceScoreResponse, error) {
	var faceScores []model.FaceScores
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

	var faceScores []model.FaceScores

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
