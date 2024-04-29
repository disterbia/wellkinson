// /medicine-service/dto/dto.go
package dto

type GetParams struct {
	StartDate string `form:"start_date" example:"YYYY-MM-DD"`
	EndDate   string `form:"end_date" example:"YYYY-MM-DD"`
}

type MedicineRequest struct {
	Id            uint     `json:"id"`
	Uid           uint     `json:"-"`
	Timestamp     []string `json:"timestamp" example:"HH:mm,HH:mm"`
	Weekdays      []uint   `json:"weekdays"`
	Dose          float32  `json:"dose"`
	IntervalType  uint8    `json:"interval_type"`
	IsActive      *bool    `json:"is_active"`
	LeastStore    *float32 `json:"least_store"`
	UseLeastStore *bool    `json:"use_least_store"`
	MedicineType  string   `json:"medicine_type"`
	Name          string   `json:"name"`
	Store         float32  `json:"store"`
	StartAt       string   `json:"start_at" example:"YYYY-MM-dd"`
	EndAt         string   `json:"end_at"  example:"YYYY-MM:dd"`
	UsePrivacy    *bool    `json:"use_privacy"`
}

type MedicineOriginResponse struct {
	Id            uint     `json:"id"`
	Timestamp     []string `json:"timestamp" example:"HH:mm,HH:mm"`
	Weekdays      []uint   `json:"weekdays"`
	Dose          float32  `json:"dose"`
	IntervalType  uint8    `json:"interval_type"`
	IsActive      bool     `json:"is_active"`
	LeastStore    float32  `json:"least_store"`
	UseLeastStore bool     `json:"use_least_store"`
	MedicineType  string   `json:"medicine_type"`
	Name          string   `json:"name"`
	Store         float32  `json:"store"`
	StartAt       string   `json:"start_at" example:"YYYY-MM-dd"`
	EndAt         string   `json:"end_at"  example:"YYYY-MM:dd"`
	UsePrivacy    bool     `json:"use_privacy"`
	Created       string   `json:"created"  example:"YYYY-mm-ddTHH:mm:ss "`
	Updated       string   `json:"updated"  example:"YYYY-mm-ddTHH:mm:ss "`
}

type MedicineDateInfo struct {
	Date      string             `json:"date" example:"YYYY-MM-DD"`
	Medicines []MedicineResponse `json:"medicines"`
}

type MedicineTakenInfo2 struct {
	Medicine MedicineResponse `json:"medicine"`
	Taken    bool             `json:"taken"`
}
type MedicineResponse struct {
	Id            uint              `json:"id"`
	Timestamp     map[string]string `json:"timestamp"`
	Weekdays      []uint            `json:"weekdays"`
	Dose          float32           `json:"dose"`
	IntervalType  uint8             `json:"interval_type"`
	IsActive      bool              `json:"is_active"`
	LeastStore    float32           `json:"least_store"`
	UseLeastStore bool              `json:"use_least_store"`
	MedicineType  string            `json:"medicine_type"`
	Name          string            `json:"name"`
	Store         float32           `json:"store"`
	StartAt       string            `json:"start_at" example:"YYYY-MM-dd"`
	EndAt         string            `json:"end_at"  example:"YYYY-MM:dd"`
	UsePrivacy    bool              `json:"use_privacy"`
	Created       string            `json:"created"  example:"YYYY-mm-ddTHH:mm:ss "`
	Updated       string            `json:"updated"  example:"YYYY-mm-ddTHH:mm:ss "`
}

type TakeMedicine struct {
	Uid        uint    `json:"-"`
	MedicineId uint    `json:"medicine_id"`
	DateTaken  string  `json:"date_taken"  example:"YYYY-MM-DD"`
	TimeTaken  string  `json:"time_taken"  example:"HH:mm"`
	RealTaken  string  `json:"real_taken"  example:"HH:mm"`
	Dose       float32 `json:"dose"`
}

type UnTakeMedicine struct {
	Uid        uint   `json:"-"`
	MedicineId uint   `json:"medicine_id"`
	DateTaken  string `json:"date_taken"  example:"YYYY-MM-DD"`
	TimeTaken  string `json:"time_taken"  example:"YYYY-MM-DD"`
}

type MedicinBridge struct {
	Id            uint    `json:"id"`
	Weekdays      []uint  `json:"weekdays"`
	Dose          float32 `json:"dose"`
	IntervalType  uint8   `json:"interval_type"`
	IsActive      bool    `json:"is_active"`
	LeastStore    float32 `json:"least_store"`
	UseLeastStore bool    `json:"use_least_store"`
	MedicineType  string  `json:"medicine_type"`
	Name          string  `json:"name"`
	Store         float32 `json:"store"`
	StartAt       string  `json:"start_at" example:"YYYY-MM-dd"`
	EndAt         string  `json:"end_at"  example:"YYYY-MM:dd"`
	UsePrivacy    bool    `json:"use_privacy"`
	Created       string  `json:"created"  example:"YYYY-mm-ddTHH:mm:ss "`
	Updated       string  `json:"updated"  example:"YYYY-mm-ddTHH:mm:ss "`
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}
type ErrorResponse struct {
	Err string `json:"err"`
}

type BasicResponse struct {
	Code string `json:"code"`
}
