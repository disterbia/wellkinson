// /inquire-service/service/inquire-service.go

package service

import (
	"common/model"
	"common/util"
	"context"
	"errors"
	"log"

	"inquire-service/dto"
	pb "inquire-service/proto"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type InquireService interface {
	AnswerInquire(answer dto.InquireReplyRequest) (string, error)
	SendInquire(inquire dto.InquireRequest) (string, error)
	GetMyInquires(id int, page int, startDate, endDate string) ([]dto.InquireResponse, error)
	GetAllInquires(id int, page int, startDate, endDate string) ([]dto.InquireResponse, error)
	RemoveInquire(id int, uid int) (string, error)
	RemoveReply(id int, uid int) (string, error)
}

type inquireService struct {
	db          *gorm.DB
	emailClient pb.EmailServiceClient
}

func NewInquireService(db *gorm.DB, conn *grpc.ClientConn) InquireService {
	emailClient := pb.NewEmailServiceClient(conn)
	return &inquireService{
		db:          db,
		emailClient: emailClient,
	}
}
func (service *inquireService) GetMyInquires(id int, page int, startDate, endDate string) ([]dto.InquireResponse, error) {

	if startDate != "" {
		if err := util.ValidateDate(startDate); err != nil {
			return nil, err
		}
	}
	if endDate != "" {
		if err := util.ValidateDate(endDate); err != nil {
			return nil, err
		}
	}

	pageSize := 10
	var inquires []model.Inquire
	offset := page * pageSize

	query := service.db.Where("uid = ? AND level = 0", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")
	result := query.Offset(offset).Limit(pageSize).Preload("Replies").Find(&inquires)

	if result.Error != nil {
		return nil, result.Error
	}

	var inquireResponses []dto.InquireResponse

	if err := util.CopyStruct(inquires, &inquireResponses); err != nil {
		return []dto.InquireResponse{}, err
	}

	return inquireResponses, nil
}

func (service *inquireService) GetAllInquires(id int, page int, startDate, endDate string) ([]dto.InquireResponse, error) {

	if startDate != "" {
		if err := util.ValidateDate(startDate); err != nil {
			return nil, err
		}
	}
	if endDate != "" {
		if err := util.ValidateDate(endDate); err != nil {
			return nil, err
		}
	}

	pageSize := 10
	var inquires []model.Inquire
	offset := page * pageSize

	var user model.User
	result := service.db.First(&user, id)

	if result.Error != nil {
		return nil, errors.New("db error")
	}

	if !user.IsAdmin {
		return nil, errors.New("unauthorized: user is not an admin")
	}

	query := service.db.Model(&model.Inquire{})

	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate)
	}
	query = query.Order("id DESC")
	result = query.Offset(offset).Limit(pageSize).Preload("Replies", "level= 0").Find(&inquires)

	if result.Error != nil {
		return nil, result.Error
	}

	var inquireResponses []dto.InquireResponse
	if err := util.CopyStruct(inquires, &inquireResponses); err != nil {
		return []dto.InquireResponse{}, err
	}

	return inquireResponses, nil
}

func (service *inquireService) SendInquire(inquireRequest dto.InquireRequest) (string, error) {
	var inquire model.Inquire
	if err := util.CopyStruct(inquireRequest, &inquire); err != nil {
		return "", err
	}

	inquire.Uid = inquireRequest.Uid
	result := service.db.Create(&inquire)

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *inquireService) AnswerInquire(inquireReplyRequest dto.InquireReplyRequest) (string, error) {
	var user model.User
	var inquire model.Inquire
	var inquireReply model.InquireReply

	result := service.db.First(&user, inquireReplyRequest.Uid)

	if result.Error != nil {
		return "", errors.New("db error")
	}

	result2 := service.db.First(&inquire, inquireReplyRequest.InquireId)

	if result2.Error != nil {
		return "", errors.New("db error2")
	}

	if inquireReplyRequest.ReplyType { // true = 답변
		if !user.IsAdmin {
			return "", errors.New("unauthorized: user is not an admin")
		}
	} else { // 추가문의
		if user.Id != inquire.Uid {
			return "", errors.New("unauthorized: illegal user")
		}
	}

	if err := util.CopyStruct(inquireReplyRequest, &inquireReply); err != nil {
		return "", err
	}
	inquireReply.Uid = inquireReplyRequest.Uid
	result = service.db.Create(&inquireReply)

	if result.Error != nil {
		return "", errors.New("db error")
	}

	go func() {
		reponse, err := service.emailClient.SendEmail(context.Background(), &pb.EmailRequest{
			Email:        inquire.Email,        // 받는 사람의 이메일
			Created:      inquire.Created,      // 문의 생성 날짜
			Title:        inquire.Title,        // 이메일 제목
			Content:      inquire.Content,      // 문의 내용
			ReplyContent: inquireReply.Content, // 답변 내용
			ReplyCreated: inquireReply.Created, // 답변 생성 날짜
		})
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}
		log.Printf(" send email: %v", reponse)
	}()

	return "200", nil
}

func (service *inquireService) RemoveInquire(id int, uid int) (string, error) {

	var inquire model.Inquire
	result := service.db.Model(&inquire).Where("id = ? AND uid = ?", id, uid).Select("level").Updates(map[string]interface{}{"level": 10})
	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *inquireService) RemoveReply(id int, uid int) (string, error) {

	var inquireReply model.InquireReply
	result := service.db.Model(&inquireReply).Where("id = ? AND uid = ?", id, uid).Select("level").Updates(map[string]interface{}{"level": 10})

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}
