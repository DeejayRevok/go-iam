package app

import (
	"fmt"
	"go-iam/src/infrastructure/logging"
	"os"

	postgres "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ConnectDatabase(logger *zap.Logger, tracedLogger *logging.ZapGormTracedLogger) *gorm.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: tracedLogger,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error connecting to database: %s", err))
	}
	return db
}
