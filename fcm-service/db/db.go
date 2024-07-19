// /login-service/pkg/db/db.go
package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 데이터베이스 연결을 초기화
func NewDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
