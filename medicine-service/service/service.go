// /medicine-service/service/service.go
package service

import (
	"common/model"
	"common/util"
	"context"
	"errors"
	"fmt"
	"log"
	"medicine-service/dto"
	pb "medicine-service/proto"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type MedicineService interface {
	SaveMedicine(medicineRequest dto.MedicineRequest) (string, error)
	RemoveMedicines(ids []uint, uid uint) (string, error)
	GetTakens(id uint, startDateStr, endDateStr string) ([]dto.MedicineDateInfo, error)
	GetMedicines(id uint) ([]dto.MedicineResponse, error)
	TakeMedicine(takeMedicine dto.TakeMedicine) (string, error)
	UnTakeMedicine(takeMedicine dto.UnTakeMedicine) (string, error)
	SearchMedicines(keyword string) ([]string, error)
}

type medicineService struct {
	db          *gorm.DB
	alarmClient pb.AlarmServiceClient
}

func NewMedicineService(db *gorm.DB, conn *grpc.ClientConn) MedicineService {
	alarmClient := pb.NewAlarmServiceClient(conn)
	return &medicineService{db: db, alarmClient: alarmClient}

}

func (service *medicineService) SaveMedicine(medicineRequest dto.MedicineRequest) (string, error) {
	if err := validateMedicine(medicineRequest); err != nil {
		return "", err
	}
	var medicine model.Medicine

	result := service.db.Where("id=? AND uid=?", medicineRequest.Id, medicineRequest.Uid).First(&model.Medicine{})

	if err := util.CopyStruct(medicineRequest, &medicine); err != nil {
		return "", err
	}

	medicine.Uid = medicineRequest.Uid //  json: "-" 이라서

	newWeekdays, unique, err := validateWeek(medicine.Weekdays)
	if err != nil {
		return "", err
	}
	medicine.Weekdays = newWeekdays
	bodyMessage := "약 먹을 시간입니다. 드시고 나면 잊지 말고 표시해주세요."

	mar := &pb.MultiAlarmRequest{}
	var ars []*pb.AlarmRequest

	if medicine.UsePrivacy {
		bodyMessage = medicine.Name + " " + fmt.Sprintf("%v", medicine.Dose) + " " + medicine.MedicineType + " 먹을 시간입니다. 드시고 나면 잊지 말고 표시해주세요."
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		medicine.Id = 0
		medicine.IsActive = true
		if err := service.db.Create(&medicine).Error; err != nil {
			return "", err
		}

		for _, v := range medicineRequest.Timestamp {
			ar := &pb.AlarmRequest{
				ParentId:  int32(medicine.Id),
				Uid:       int32(medicine.Uid),
				Body:      bodyMessage,
				Type:      int32(util.MedicineType),
				StartAt:   medicine.StartAt,
				EndAt:     medicine.EndAt,
				Timestamp: v,
				Week:      unique,
			}
			ars = append(ars, ar)
		}

		mar.AlarmRequests = ars

		go sendAlarm(service, mar)

	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&medicine).Updates(medicine).Error; err != nil {
			return "", err
		}
		for _, v := range medicineRequest.Timestamp {
			ar := &pb.AlarmRequest{
				ParentId:  int32(medicine.Id),
				Uid:       int32(medicine.Uid),
				Body:      bodyMessage,
				Type:      int32(util.MedicineType),
				StartAt:   medicine.StartAt,
				EndAt:     medicine.EndAt,
				Timestamp: v,
				Week:      unique,
			}
			ars = append(ars, ar)
		}

		mar.AlarmRequests = ars
		if medicine.IsActive {
			go updateAlarm(service, mar)
		} else {
			b := make([]int32, 1)

			b[0] = int32(medicine.Id)

			arr := &pb.AlarmRemoveRequest{
				ParentIds: b,
				Uid:       int32(medicine.Uid),
				Type:      int32(util.MedicineType),
			}
			go removeAlarm(service, arr)
		}
	}

	return "200", nil
}

func (service *medicineService) RemoveMedicines(ids []uint, uid uint) (string, error) {
	result := service.db.Where("id IN (?) AND uid= ?", ids, uid).Delete(&model.Medicine{})

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
		Type:      int32(util.MedicineType),
	}

	go removeAlarm(service, arr)
	return "200", nil
}

func (service *medicineService) GetTakens(id uint, startDateStr, endDateStr string) ([]dto.MedicineDateInfo, error) {

	// 문자열을 time.Time 타입으로 변환
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)

	if err != nil {
		return nil, err
	}

	var medicines []model.Medicine
	var medicineResponse []dto.MedicineResponse
	err = service.db.Where("uid = ? AND start_at <= ? AND end_at >= ?", id, endDate, startDate).Find(&medicines).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(medicines, &medicineResponse); err != nil {
		return nil, err
	}

	// 해당 약물 복용내역 조회
	var medicineIds []uint
	for _, medicine := range medicineResponse {
		medicineIds = append(medicineIds, medicine.Id)
	}
	var takenMedicines []model.MedicineTake
	err = service.db.Where("medicine_id IN (?) AND date_taken BETWEEN ? AND ?", medicineIds, startDate, endDate).Find(&takenMedicines).Error
	if err != nil {
		return nil, err
	}

	// 복용내역 응답형식으로 가공 및 복용일 반영
	takenMap := make(map[uint]map[string]bool)
	for _, tm := range takenMedicines {
		if takenMap[tm.MedicineId] == nil {
			takenMap[tm.MedicineId] = make(map[string]bool)
		}
		takenMap[tm.MedicineId][tm.DateTaken] = true
	}

	// 전체날짜에서 약물 복용날짜 체크
	medicineDates := make([]dto.MedicineDateInfo, 0)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		var dailyMediciens []dto.MedicineTakenInfo
		for _, m := range medicineResponse {
			startAt, _ := time.Parse("2006-01-02", m.StartAt)
			endAt, _ := time.Parse("2006-01-02", m.EndAt)

			if d.After(startAt) && d.Before(endAt) && isMedicineDay(m.Weekdays, d.Weekday()) {
				taken := takenMap[m.Id][d.Format("2006-01-02")]
				dailyMediciens = append(dailyMediciens, dto.MedicineTakenInfo{Medicine: m, Taken: taken})
			}
		}
		if len(dailyMediciens) > 0 {
			medicineDates = append(medicineDates, dto.MedicineDateInfo{Date: d.Format("2006-01-02"), Medicines: dailyMediciens})
		}
	}

	return medicineDates, nil
}

func (service *medicineService) GetMedicines(id uint) ([]dto.MedicineResponse, error) {
	var medicines []model.Medicine
	var medicineResponses []dto.MedicineResponse

	err := service.db.Where("uid = ?", id).Find(&medicines).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(medicines, &medicineResponses); err != nil {
		return nil, err
	}

	return medicineResponses, nil
}

func (service *medicineService) TakeMedicine(takeMedicine dto.TakeMedicine) (string, error) {
	var medicineTake model.MedicineTake
	if err := util.ValidateTime(takeMedicine.TimeTaken); err != nil {
		return "", err
	}
	if err := util.ValidateDate(takeMedicine.DateTaken); err != nil {
		return "", err
	}

	result := service.db.Where("medicine_id = ? AND uid = ? AND date_taken = ? AND time_taken = ? ", takeMedicine.MedicineId, takeMedicine.Uid, takeMedicine.DateTaken, takeMedicine.TimeTaken).First(&medicineTake)

	if result.RowsAffected == 0 {

		medicineTake.Uid = takeMedicine.Uid
		if err := util.CopyStruct(takeMedicine, &medicineTake); err != nil {
			return "", err
		}
		if err := service.db.Create(&medicineTake).Error; err != nil {
			return "", errors.New("db error")
		}
	} else {
		return "", errors.New("duplicated")

	}
	return "200", nil
}

func (service *medicineService) UnTakeMedicine(unTakeMedicine dto.UnTakeMedicine) (string, error) {

	if err := util.ValidateTime(unTakeMedicine.TimeTaken); err != nil {
		return "", err
	}

	if err := util.ValidateDate(unTakeMedicine.DateTaken); err != nil {
		return "", err
	}

	result := service.db.Where("medicine_id = ? AND uid=? AND date_taken = ? AND time_taken = ?", unTakeMedicine.MedicineId, unTakeMedicine.Uid, unTakeMedicine.DateTaken, unTakeMedicine.TimeTaken).Delete(&model.MedicineTake{})
	if result.Error != nil {
		return "", errors.New("db error2")
	}

	return "200", nil
}

func (service *medicineService) SearchMedicines(keyword string) ([]string, error) {
	var names []string
	err := service.db.Model(&model.MedicineSearch{}).Where("name LIKE ?", "%"+keyword+"%").Pluck("name", &names).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	return names, nil
}

func isMedicineDay(weekdays []uint, day time.Weekday) bool {
	for _, d := range weekdays {
		if uint(day) == d {
			return true
		}
	}
	return false
}

func sendAlarm(service *medicineService, mar *pb.MultiAlarmRequest) {
	reponse, err := service.alarmClient.MultiSetAlarm(context.Background(), mar)
	if err != nil {
		log.Printf("Failed to set Alarm: %v", err)
	}
	log.Printf("set Alarm: %v", reponse)
}

func removeAlarm(service *medicineService, arr *pb.AlarmRemoveRequest) {
	reponse, err := service.alarmClient.RemoveAlarm(context.Background(), arr)
	if err != nil {
		log.Printf("Failed to remove Alarm: %v", err)
	}
	log.Printf("remove Alarm: %v", reponse)
}

func updateAlarm(service *medicineService, mar *pb.MultiAlarmRequest) {
	reponse, err := service.alarmClient.MultiUpdateAlarm(context.Background(), mar)
	if err != nil {
		log.Printf("Failed to update Alarm: %v", err)
	}
	log.Printf("update Alarm: %v", reponse)
}
