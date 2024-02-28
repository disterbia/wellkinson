// /alarm-service/service/grpc-service.go
package service

import (
	"alarm-service/common/model"
	"alarm-service/common/util"
	pb "alarm-service/proto"
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type AlarmServer struct {
	pb.UnimplementedAlarmServiceServer
	Db *gorm.DB
}
type TempAlarm struct {
	Id        uint            `gorm:"primaryKey;autoIncrement"`
	Uid       uint            `gorm:"not null"`
	Type      string          `gorm:"size:255;not null"`
	Body      string          `gorm:"type:text;not null"`
	StartAt   string          `gorm:"size:255;not null" json:"start_at"`
	EndAt     string          `gorm:"size:255;not null" json:"end_at"`
	Timestamp string          `gorm:"size:255;not null"`
	Week      json.RawMessage `gorm:"type:json"`
}

func (s *AlarmServer) SetAlarm(ctx context.Context, req *pb.AlarmRequest) (*pb.AlarmResponse, error) {

	var alarm model.Alarm

	if err := util.CopyStruct(req, &alarm); err != nil {
		return nil, err
	}
	if err := s.Db.Create(&alarm).Error; err != nil {
		return nil, errors.New("db error")
	}

	return &pb.AlarmResponse{Status: "Success"}, nil

}

func (s *AlarmServer) RemoveAlarm(ctx context.Context, req *pb.AlarmRemoveRequest) (*pb.AlarmResponse, error) {

	result := s.Db.Where("parent_id IN ? AND uid= ? AND type=?", req.ParentIds, req.Uid, req.Type).Delete(&model.Alarm{})

	if result.Error != nil {
		return nil, errors.New("db error")
	}

	return &pb.AlarmResponse{Status: "Success"}, nil
}

func (s *AlarmServer) UpdateAlarm(ctx context.Context, req *pb.AlarmRequest) (*pb.AlarmResponse, error) {

	var alarm model.Alarm

	result := s.Db.Where("parent_id = ? AND uid= ? AND type=?", req.ParentId, req.Uid, req.Type).Delete(&model.Alarm{})
	if result.Error != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(req, &alarm); err != nil {
		return nil, err
	}
	if err := s.Db.Create(&alarm).Error; err != nil {
		return nil, errors.New("db error")
	}

	return &pb.AlarmResponse{Status: "Success"}, nil

}

func (s *AlarmServer) MultiSetAlarm(ctx context.Context, req *pb.MultiAlarmRequest) (*pb.AlarmResponse, error) {

	var alarms []model.Alarm

	if err := util.CopyStruct(req.AlarmRequests, &alarms); err != nil {
		return nil, err
	}
	if err := s.Db.Create(&alarms).Error; err != nil {
		return nil, errors.New("db error")
	}

	return &pb.AlarmResponse{Status: "Success"}, nil

}

func (s *AlarmServer) MultiUpdateAlarm(ctx context.Context, req *pb.MultiAlarmRequest) (*pb.AlarmResponse, error) {

	var alarms []model.Alarm
	var ids []int32
	for _, v := range req.AlarmRequests {
		ids = append(ids, v.ParentId)
	}
	result := s.Db.Where("parent_id IN ? AND uid= ? AND type=?", ids, req.AlarmRequests[0].Uid, req.AlarmRequests[0].Type).Delete(&model.Alarm{})
	if result.Error != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(req.AlarmRequests, &alarms); err != nil {
		return nil, err
	}
	if err := s.Db.Create(&alarms).Error; err != nil {
		return nil, errors.New("db error")
	}

	return &pb.AlarmResponse{Status: "Success"}, nil

}
