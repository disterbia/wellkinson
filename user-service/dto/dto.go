// /user-service/dto/dto.go
package dto

type UserRequest struct {
	Id                   int    `json:"-"`
	Email                string `json:"-"`
	Birthday             string `json:"birthday" example:"YYYY-MM-DD"`
	DeviceID             string `json:"device_id"`
	Gender               bool   `json:"gender"` // true:남 false: 여
	FCMToken             string `json:"fcm_token"`
	IsFirst              bool   `json:"is_first"`
	Name                 string `json:"name"`
	PhoneNum             string `json:"phone_num" example:"01000000000"`
	UseAutoLogin         bool   `json:"use_auto_login"`
	UsePrivacyProtection bool   `json:"user_privacy_protection"`
	UseSleepTracking     bool   `json:"use_sleep_tracking"`
	UserType             string `json:"user_type"`
}

type UserResponse struct {
	Birthday             string `json:"birthday" example:"YYYY-MM-DD"`
	DeviceID             string `json:"device_id"`
	Gender               bool   `json:"gender"` // true:남 false: 여
	FCMToken             string `json:"fcm_token"`
	IsFirst              bool   `json:"is_first"`
	Name                 string `json:"name"`
	PhoneNum             string `json:"phone_num" example:"01000000000"`
	UseAutoLogin         bool   `json:"use_auto_login"`
	UsePrivacyProtection bool   `json:"user_privacy_protection"`
	UseSleepTracking     bool   `json:"use_sleep_tracking"`
	UserType             string `json:"user_type"`
	Email                string `json:"email"`
	Created              string `json:"created"`
	Updated              string `json:"updated"`
}

type LoginRequest struct {
	IdToken     string      `json:"id_token"`
	UserRequest UserRequest `json:"user"`
}

type LoginResponse struct {
	Jwt string `json:"jwt,omitempty"`
	Err string `json:"err,omitempty"`
}

type AutoLoginRequest struct {
	Jwt string `json:"jwt"`
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
