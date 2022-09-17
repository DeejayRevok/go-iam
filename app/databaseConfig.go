package app

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

func ConnectDatabase(logger *zap.Logger) *gorm.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error connecting to database: %s", err))
	}
	if err = db.Use(gormopentracing.New()); err != nil {
		logger.Fatal(fmt.Sprintf("Error setting tracing for gorm: %s", err))
	}
	return db
}
