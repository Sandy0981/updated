package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"job-portal-api/internal/models"
)

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate function will ONLY create tables, missing columns and missing indexes, and WON'T change existing column's type or delete unused columns
	err = db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Company{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Jobs{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return nil, err
	}
	return db, nil
}
