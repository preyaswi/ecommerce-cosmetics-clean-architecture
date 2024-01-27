package db

import (
	"clean/pkg/config"
	"clean/pkg/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDatabase(cfg config.Config) (*gorm.DB,error) {
	psqlInfo:=fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db,dbErr:=gorm.Open(postgres.Open(psqlInfo),&gorm.Config{SkipDefaultTransaction: true})

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Address{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Products{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Category{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Cart{}); err != nil {
		return db, err
	}



	return db,dbErr
}