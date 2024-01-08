// /inquire-service/service/inquire-service.go

package service

import (
	"common/model"
	"common/util"
	"context"
	"errors"
	"log"

	pb "inquire-service/proto"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type InquireService interface {
	AnswerInquire(answer model.InquireReply) (string, error)
	SendInquire(inquire model.Inquire) (string, error)
	GetMyInquires(id int, page int, startDate, endDate string) ([]model.Inquire, error)
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
func (service *inquireService) GetMyInquires(id int, page int, startDate, endDate string) ([]model.Inquire, error) {

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

	query := service.db.Where("uid = ?", id)
	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate)
	}

	result := query.Offset(offset).Limit(pageSize).Preload("Replies").Find(&inquires)

	if result.Error != nil {
		return nil, result.Error
	}

	return inquires, nil
}

func (service *inquireService) SendInquire(inquire model.Inquire) (string, error) {
	result := service.db.Save(&inquire)

	if result.Error != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *inquireService) AnswerInquire(answer model.InquireReply) (string, error) {
	var user model.User
	var inquire model.Inquire

	result := service.db.First(&user, answer.Uid)

	if result.Error != nil {
		return "", errors.New("db error")
	}

	result2 := service.db.First(&inquire, answer.InquireId)

	if result2.Error != nil {
		return "", errors.New("db error2")
	}

	if answer.ReplyType { // true = 답변
		if !user.IsAdmin {
			return "", errors.New("unauthorized: user is not an admin")
		}
	} else { // 추가문의
		if user.Id != inquire.Uid {
			return "", errors.New("unauthorized: illegal user")
		}
	}

	answer.Id = 0
	result = service.db.Save(&answer)

	if result.Error != nil {
		return "", errors.New("db error")
	}

	go func() {
		reponse, err := service.emailClient.SendEmail(context.Background(), &pb.EmailRequest{
			Email:        inquire.Email,                                       // 받는 사람의 이메일
			Created:      inquire.Created.Format("2006년 01월 02일 15시 04분 05초"), // 문의 생성 날짜
			Title:        inquire.Title,                                       // 이메일 제목
			Content:      inquire.Content,                                     // 문의 내용
			ReplyContent: answer.Content,                                      // 답변 내용
			ReplyCreated: answer.Created.Format("2006년 01월 02일 15시 04분 05초"),  // 답변 생성 날짜
		})
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}
		log.Printf(" send email: %v", reponse)
	}()

	return "200", nil
}
