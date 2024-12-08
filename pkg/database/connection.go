package database

import (
	"fmt"
	"log"

	"github.com/stwrtrio/coffee-shop/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect connects to MySQL database and returns the DB object.
func Connect(config *utils.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	}
	return db, nil
}
