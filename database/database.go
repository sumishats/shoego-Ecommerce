package database

import (
	"fmt"
	"shoego/config"
	"shoego/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)

	fmt.Println("connecting database", cfg.DBName)

	db, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if dberr != nil {
		return nil, fmt.Errorf("faild to connect to database:%w", dberr)
	}

	DB = db
	DB.AutoMigrate(&domain.User{})
	DB.AutoMigrate(&domain.OTPVerification{})

	fmt.Println("Database connected successfully")

	return DB, nil
}
