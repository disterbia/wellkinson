// /user-service/pkg/service/login-service.go

package service

import (
	"common/model"
)

type LoginService interface {
	Login(token string, user model.User) (string, error)
}
