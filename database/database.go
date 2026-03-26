package database

import (
	"fmt"
	"log"
	"shoego/config"
	"shoego/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)

	log.Println("connecting database", cfg.DBName)

	db, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if dberr != nil {
		return nil, fmt.Errorf("faild to connect to database:%w", dberr)
	}

	DB = db
	DB.AutoMigrate(&domain.User{})
	DB.AutoMigrate(&domain.OTPVerification{})
	DB.AutoMigrate(&domain.Address{})
	DB.AutoMigrate(&domain.Product{})
	DB.AutoMigrate(&domain.ProductImage{})
	DB.AutoMigrate(&domain.Category{})
	DB.AutoMigrate(&domain.BlacklistToken{})
	DB.AutoMigrate(&domain.Cart{})
	DB.AutoMigrate(&domain.CartItem{})

	log.Println("database connected successfully")

	return DB, nil
}
