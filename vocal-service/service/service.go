// /vocal-service/service/service.go
package service

import (
	"errors"
	"vocal-service/common/model"
	"vocal-service/common/util"
	"vocal-service/dto"

	"gorm.io/gorm"
)

type VocalService interface {
	GetVoiceTables() ([]dto.VoiceWordResponse, error)
	GetVocalScores(id uint, startDate, endDate string) ([]dto.VocalScoreResponse, error)
	SaveVocalScores(vocalScoreRequest []dto.VocalScoreRequest) (string, error)
}

type vocalService struct {
	db *gorm.DB
}

func NewVocalService(db *gorm.DB) VocalService {
	return &vocalService{db: db}
}
func (service *vocalService) GetVoiceTables() ([]dto.VoiceWordResponse, error) {
	var voiceWords []model.VocalWord
	var voiceResponses []dto.VoiceWordResponse

	err := service.db.Find(&voiceWords).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(voiceWords, &voiceResponses); err != nil {
		return nil, err
	}

	return voiceResponses, nil
}

// exercise-servie의 video는 따로 가져옴 둘중 좋은방법 선택하기

func (service *vocalService) GetVocalScores(id uint, startDate, endDate string) ([]dto.VocalScoreResponse, error) {
	var vocalScores []model.VocalScore
	var vocalScoreResponse []dto.VocalScoreResponse

	query := service.db.Where("uid = ?", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")

	err := query.Find(&vocalScores).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(vocalScores, &vocalScoreResponse); err != nil {
		return nil, err
	}

	return vocalScoreResponse, nil
}

func (service *vocalService) SaveVocalScores(vocalScoreRequest []dto.VocalScoreRequest) (string, error) {

	var vocalScores []model.VocalScore

	if err := util.CopyStruct(vocalScoreRequest, &vocalScores); err != nil {
		return "", err
	}
	var len = len(vocalScoreRequest)
	for i := 0; i < len; i++ {
		vocalScores[i].Uid = vocalScoreRequest[0].Uid
	}

	if err := service.db.Create(&vocalScores).Error; err != nil {
		return "", err
	}

	return "200", nil
}
