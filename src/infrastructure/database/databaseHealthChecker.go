package database

import "gorm.io/gorm"

type DatabaseHealthChecker struct {
	db *gorm.DB
}

func (checker *DatabaseHealthChecker) Check() error {
	if pinger, ok := checker.db.ConnPool.(interface{ Ping() error }); ok {
		return pinger.Ping()
	}
	return nil
}

func NewDatabaseHealthChecker(db *gorm.DB) *DatabaseHealthChecker {
	checker := DatabaseHealthChecker{
		db: db,
	}
	return &checker
}
