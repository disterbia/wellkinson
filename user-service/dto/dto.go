// /user-service/dto/dto.go
package dto

type UserRequest struct {
	Id                    uint   `json:"-"`
	Email                 string `json:"-"`
	Birthday              string `json:"birthday" example:"YYYY-MM-DD"`
	DeviceID              string `json:"device_id"`
	Gender                *bool  `json:"gender"` // true:남 false: 여
	FCMToken              string `json:"fcm_token"`
	IsFirst               *bool  `json:"is_first"`
	SnsType               uint   `json:"sns_type"`
	Name                  string `json:"name"`
	PhoneNum              string `json:"phone_num" example:"01000000000"`
	UseAutoLogin          *bool  `json:"use_auto_login"`
	IndemnificationClause *bool  `json:"indemnification_clause"`
	UsePrivacyProtection  *bool  `json:"user_privacy_protection"`
	UseSleepTracking      *bool  `json:"use_sleep_tracking"`
	UserType              *uint  `json:"user_type"`
	UserServices          []int  `json:"user_services"`
	ProfileImage          string `json:"profile_image" example:"base64 encoding string"`
}

type UserResponse struct {
	Birthday              string                `json:"birthday" example:"YYYY-MM-DD"`
	DeviceID              string                `json:"device_id"`
	Gender                bool                  `json:"gender"` // true:남 false: 여
	FCMToken              string                `json:"fcm_token"`
	IsFirst               bool                  `json:"is_first"`
	Name                  string                `json:"name"`
	PhoneNum              string                `json:"phone_num" example:"01000000000"`
	UseAutoLogin          bool                  `json:"use_auto_login"`
	IndemnificationClause bool                  `json:"indemnification_clause"`
	UsePrivacyProtection  bool                  `json:"user_privacy_protection"`
	UseSleepTracking      bool                  `json:"use_sleep_tracking"`
	UserType              uint                  `json:"user_type"`
	SnsType               uint                  `json:"sns_type"`
	Email                 string                `json:"email"`
	Created               string                `json:"created" example:"YYYY-mm-ddTHH:mm:ss "`
	Updated               string                `json:"updated" example:"YYYY-mm-ddTHH:mm:ss "`
	UserServices          []MainServiceResponse `json:"user_services"`
	ProfileImage          ImageResponse         `json:"profile_image"`
	LinkedEmails          []LinkedResponse      `json:"linked_emails"`
}

// func (r *UserResponse) MarshalJSON() ([]byte, error) {
// 	type Alias UserResponse
// 	loc, _ := time.LoadLocation("Asia/Seoul")
// 	parsedTime, _ := time.Parse(time.RFC3339, r.Created)
// 	seoulTime := parsedTime.In(loc).Format("2006-01-02 15:04:05")
// 	return json.Marshal(&struct {
// 		Created string `json:"created"`
// 		*Alias
// 	}{
// 		Created: seoulTime,
// 		Alias:   (*Alias)(r),
// 	})
// }

type ImageResponse struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type LinkedResponse struct {
	SnsType uint   `json:"sns_type"`
	Email   string `json:"email"`
}

type MainServiceResponse struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

type TempUser struct {
	Id                    uint   `json:"-"`
	Email                 string `json:"-"`
	Birthday              string `json:"birthday" example:"YYYY-MM-DD"`
	DeviceID              string `json:"device_id"`
	Gender                bool   `json:"gender"` // true:남 false: 여
	FCMToken              string `json:"fcm_token"`
	IndemnificationClause bool   `json:"indemnification_clause"`
	IsFirst               bool   `json:"is_first"`
	SnsType               uint   `json:"sns_type"`
	Name                  string `json:"name"`
	PhoneNum              string `json:"phone_num" example:"01000000000"`
	UseAutoLogin          bool   `json:"use_auto_login"`
	UsePrivacyProtection  bool   `json:"user_privacy_protection"`
	UseSleepTracking      bool   `json:"use_sleep_tracking"`
	UserType              uint   `json:"user_type"`
	UserServices          []int  `json:"user_services"`
}

type LoginRequest struct {
	IdToken     string      `json:"id_token"`
	UserRequest UserRequest `json:"user"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phone_number" example:"01000000000"`
	Code        string `json:"code" example:"인증번호 6자리"`
}

type LinkRequest struct {
	Id      uint   `json:"-"`
	IdToken string `json:"id_token"`
}

type AppVersionResponse struct {
	LatestVersion string `json:"latest_version"`
	AndroidLink   string `json:"android_link"`
	IosLink       string `json:"ios_link"`
}

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}

type AutoLoginRequest struct {
	Email    string `json:"-"`
	FcmToken string `json:"fcm_token"`
	DeviceId string `json:"device_id"`
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
