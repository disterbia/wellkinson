// /medicine-service/service/service.go
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"medicine-service/common/model"
	"medicine-service/common/util"
	"medicine-service/dto"
	pb "medicine-service/proto"
	"reflect"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type MedicineService interface {
	SaveMedicine(medicineRequest dto.MedicineRequest) (string, error)
	RemoveMedicines(ids []uint, uid uint) (string, error)
	GetTakens(id uint, startDateStr, endDateStr string) ([]dto.MedicineDateInfo, error)
	GetMedicines(id uint) ([]dto.MedicineOriginResponse, error)
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

	newWeekdays, unique, err := validateWeek(medicine)
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
	if medicine.IntervalType == 1 {
		medicine.Timestamp = json.RawMessage("[]")
		medicine.Weekdays = json.RawMessage("[]")
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		medicine.Id = 0
		medicine.IsActive = true
		if err := service.db.Create(&medicine).Error; err != nil {
			return "", err
		}

		if medicine.IntervalType != 1 {

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
		}
	} else if result.Error != nil {
		return "", errors.New("db error")
	} else {
		if medicine.IntervalType == 1 {
			medicineRequest.Timestamp = make([]string, 0)
			medicineRequest.Weekdays = make([]uint, 0)
		}
		updateFields := make(map[string]interface{})

		userRequestValue := reflect.ValueOf(medicineRequest)
		userRequestType := userRequestValue.Type()
		for i := 0; i < userRequestValue.NumField(); i++ {
			field := userRequestValue.Field(i)
			fieldName := userRequestType.Field(i).Tag.Get("json")
			if fieldName == "-" {
				continue
			}
			if !field.IsZero() {
				if fieldName == "weekdays" || fieldName == "timestamp" {
					// weekdays 필드를 JSON 형식으로 변환
					updateFields[fieldName], _ = json.Marshal(field.Interface())
				} else {
					updateFields[fieldName] = field.Interface()
				}
			}
		}
		// 레코드가 존재하면 업데이트
		if err := service.db.Model(&medicine).Updates(updateFields).Error; err != nil {
			return "", err
		}

		if medicine.IntervalType != 1 {
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

	}

	return "200", nil
}

func (service *medicineService) RemoveMedicines(ids []uint, uid uint) (string, error) {
	result := service.db.Model(&model.Medicine{}).Where("id IN (?) AND uid= ?", ids, uid).Select("is_delete").Updates(map[string]interface{}{"is_delete": true})
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
	var medicineTemp []dto.MedicineOriginResponse
	var medicineBridge []dto.MedicinBridge
	var medicineResponses []dto.MedicineResponse
	err = service.db.Debug().Where("uid = ? AND (start_at ='' OR start_at <= ?) AND (end_at ='' OR end_at >= ?)",
		id, endDate.Format("2006-01-02"), startDate.Format("2006-01-02")).Find(&medicines).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if err := util.CopyStruct(medicines, &medicineTemp); err != nil {
		return nil, err
	}

	if err := util.CopyStruct(medicineTemp, &medicineBridge); err != nil {
		return nil, err
	}

	if err := util.CopyStruct(medicineBridge, &medicineResponses); err != nil {
		return nil, err
	}

	// 해당 약물 복용내역 조회
	var medicineIds []uint
	for _, medicine := range medicineResponses {
		medicineIds = append(medicineIds, medicine.Id)
	}
	var takenMedicines []model.MedicineTake
	err = service.db.Debug().Where("medicine_id IN (?) AND date_taken BETWEEN ? AND ?",
		medicineIds, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&takenMedicines).Error
	if err != nil {
		return nil, err
	}

	takenMap := make(map[uint]map[string]map[string]map[string]float32)
	for _, tm := range takenMedicines {
		if takenMap[tm.MedicineId] == nil {
			takenMap[tm.MedicineId] = make(map[string]map[string]map[string]float32)
		}
		if takenMap[tm.MedicineId][tm.DateTaken] == nil {
			takenMap[tm.MedicineId][tm.DateTaken] = make(map[string]map[string]float32)
		}
		if takenMap[tm.MedicineId][tm.DateTaken][tm.TimeTaken] == nil {
			takenMap[tm.MedicineId][tm.DateTaken][tm.TimeTaken] = make(map[string]float32)
		}
		takenMap[tm.MedicineId][tm.DateTaken][tm.TimeTaken][tm.RealTaken] = tm.Dose
	}

	// 전체날짜에서 약물 복용날짜 체크
	medicineDates := make([]dto.MedicineDateInfo, 0)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		var dayMedicineResponses []dto.MedicineResponse
		for i, m := range medicineTemp {
			startAt, errStart := time.Parse("2006-01-02", m.StartAt)
			endAt, errEnd := time.Parse("2006-01-02", m.EndAt)

			// startAt 또는 endAt이 null (또는 파싱 에러)일 경우, 무기한으로 간주
			if errStart != nil {
				startAt = time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC) // 무기한 시작일 경우의 처리
			}
			if errEnd != nil {
				endAt = time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC) // 무기한 종료일 경우의 처리
			}

			var a = make(map[string]map[string]float32)

			// m.Timestamp가 비어있을 경우, takenMap에서 해당 날짜의 모든 시간에 대한 데이터를 가져옴

			for takenTime := range takenMap[m.Id][d.Format("2006-01-02")] {
				if takenTime != "" {
					taken := takenMap[m.Id][d.Format("2006-01-02")][takenTime]
					a[takenTime] = taken
					if !d.Before(startAt) && d.Before(endAt.AddDate(0, 0, 1)) && !(isMedicineDay(m.Weekdays, d.Weekday())) {
						tempMedicineResponse := medicineResponses[i]
						tempMedicineResponse.Timestamp = a
						dayMedicineResponses = append(dayMedicineResponses, tempMedicineResponse)
					}
				}
			}

			if !d.Before(startAt) && d.Before(endAt.AddDate(0, 0, 1)) && (isMedicineDay(m.Weekdays, d.Weekday())) {
				// var a = make(map[string]string)

				for _, v := range m.Timestamp {

					taken := takenMap[m.Id][d.Format("2006-01-02")][v]
					a[v] = taken
					// log.Println(d.Format("2006-01-02"), v, taken)

				}
				tempMedicineResponse := medicineResponses[i]
				tempMedicineResponse.Timestamp = a
				dayMedicineResponses = append(dayMedicineResponses, tempMedicineResponse)
			}
		}
		if len(medicineResponses) > 0 {
			medicineDates = append(medicineDates, dto.MedicineDateInfo{Date: d.Format("2006-01-02"), Medicines: dayMedicineResponses})
		}
	}

	return medicineDates, nil
}

func (service *medicineService) GetMedicines(id uint) ([]dto.MedicineOriginResponse, error) {
	var medicines []model.Medicine
	var medicineTemp []dto.MedicineOriginResponse

	err := service.db.Where("uid = ? AND is_delete = false ", id).Find(&medicines).Error
	if err != nil {
		return nil, errors.New("db error")
	}
	if err := util.CopyStruct(medicines, &medicineTemp); err != nil {
		return nil, err
	}

	return medicineTemp, nil
}

func (service *medicineService) TakeMedicine(takeMedicine dto.TakeMedicine) (string, error) {
	var medicineTake model.MedicineTake
	var medicine model.Medicine
	if err := util.ValidateTime(takeMedicine.TimeTaken); err != nil {
		return "", err
	}
	if err := util.ValidateTime(takeMedicine.RealTaken); err != nil {
		return "", err
	}
	if err := util.ValidateDate(takeMedicine.DateTaken); err != nil {
		return "", err
	}

	tx := service.db.Begin()
	result := service.db.Where("medicine_id = ? AND uid = ? AND date_taken = ? AND time_taken = ?",
		takeMedicine.MedicineId, takeMedicine.Uid, takeMedicine.DateTaken, takeMedicine.TimeTaken).First(&medicineTake)

	if result.RowsAffected == 0 {

		medicineTake.Uid = takeMedicine.Uid
		if err := util.CopyStruct(takeMedicine, &medicineTake); err != nil {
			return "", err
		}

		err := tx.Where("id = ? AND uid = ?", takeMedicine.MedicineId, takeMedicine.Uid).First(&medicine).Error
		if err != nil {
			return "", errors.New("db error2")
		} else {
			store := medicine.Store
			dose := takeMedicine.Dose
			if medicine.UseLeastStore {
				if store < dose {
					return "", errors.New("over dose")
				} else {
					err := tx.Model(&medicine).Select("store").Updates(map[string]interface{}{"store": store - dose}).Error
					if err != nil {
						return "", errors.New("db error3")
					}
				}
			}
		}

		if err := tx.Create(&medicineTake).Error; err != nil {
			tx.Rollback()
			return "", errors.New("db error")
		}

	} else {
		err := tx.Where("id = ? AND uid = ?", takeMedicine.MedicineId, takeMedicine.Uid).First(&medicine).Error
		if err != nil {
			return "", errors.New("db error2")
		}
		if medicine.IntervalType == 1 {
			takeMedicine.RealTaken = takeMedicine.TimeTaken
		}
		err2 := tx.Model(&medicineTake).UpdateColumn("real_taken", takeMedicine.RealTaken).Error
		if err2 != nil {
			tx.Rollback()
			return "", errors.New("db error3")
		}

	}
	tx.Commit()
	return "200", nil
}

func (service *medicineService) UnTakeMedicine(unTakeMedicine dto.UnTakeMedicine) (string, error) {

	if err := util.ValidateTime(unTakeMedicine.TimeTaken); err != nil {
		return "", err
	}

	if err := util.ValidateDate(unTakeMedicine.DateTaken); err != nil {
		return "", err
	}

	var medicineTake model.MedicineTake
	var medicine model.Medicine

	result := service.db.Where("medicine_id = ? AND uid=? AND date_taken = ? AND time_taken = ?",
		unTakeMedicine.MedicineId, unTakeMedicine.Uid, unTakeMedicine.DateTaken, unTakeMedicine.TimeTaken).First(&medicineTake)
	if result.Error != nil {
		return "", errors.New("db error2")
	}
	tx := service.db.Begin()

	result2 := tx.Where("medicine_id = ? AND uid=? AND date_taken = ? AND time_taken = ?",
		unTakeMedicine.MedicineId, unTakeMedicine.Uid, unTakeMedicine.DateTaken, unTakeMedicine.TimeTaken).Delete(&model.MedicineTake{})
	if result2.Error != nil {
		tx.Rollback()
		return "", errors.New("db error2")
	}

	err := tx.Model(&medicine).Where("id = ? AND uid = ?", medicineTake.MedicineId, medicineTake.Uid).UpdateColumn("store", gorm.Expr("store + ?", medicineTake.Dose)).Error
	if err != nil {
		tx.Rollback()
		return "", errors.New("db error3")
	}
	tx.Commit()

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
