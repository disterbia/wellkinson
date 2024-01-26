// /alarm-service/service/grpc-service.go
package service

import (
	pb "alarm-service/proto"
	"common/model"
	"common/util"
	"context"
	"encoding/json"
	"errors"
	"log"

	"gorm.io/gorm"
)

type AlarmServer struct {
	pb.UnimplementedAlarmServiceServer
	Db *gorm.DB
}
type TempAlarm struct {
	Id        int             `gorm:"primaryKey;autoIncrement"`
	Uid       int             `gorm:"not null"`
	Type      string          `gorm:"size:255;not null"`
	Body      string          `gorm:"type:text;not null"`
	StartAt   string          `gorm:"size:255;not null" json:"start_at"`
	EndAt     string          `gorm:"size:255;not null" json:"end_at"`
	Timestamp string          `gorm:"size:255;not null"`
	Week      json.RawMessage `gorm:"type:json"`
}

func (s *AlarmServer) SetAlarm(ctx context.Context, req *pb.AlarmRequest) (*pb.AlarmResponse, error) {

	var alarm model.Alarm
	var tempAlarm TempAlarm
	result := s.Db.Where("id=? AND uid=?", req.Id, req.Uid).First(&alarm)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		// 레코드가 존재하지 않으면 새 레코드 생성
		if err := util.CopyStruct(req, &alarm); err != nil {
			return nil, err
		}
		if err := s.Db.Create(&alarm).Error; err != nil {
			return nil, errors.New("db error")
		}
	} else if result.Error != nil {
		return nil, errors.New("db error2")
	} else {
		if err := util.CopyStruct(req, &tempAlarm); err != nil {
			return nil, err
		}
		log.Println(tempAlarm.Week[0])
		log.Println(tempAlarm.StartAt)
		if err := s.Db.Model(&alarm).Updates(tempAlarm).Error; err != nil {
			return nil, errors.New("db error3")
		}
	}

	return &pb.AlarmResponse{Status: "Success"}, nil

}

func (s *AlarmServer) RemoveAlarm(ctx context.Context, req *pb.AlarmRemoveRequest) (*pb.AlarmResponse, error) {

	result := s.Db.Where("id IN ? AND uid= ?", req.Ids, req.Uid).Delete(&model.Alarm{})

	if result.Error != nil {
		return nil, errors.New("db error")
	}

	return &pb.AlarmResponse{Status: "Success"}, nil
}
