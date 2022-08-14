package infra_db

import (
	"fmt"

	_config "websocket-in-go-boilerplate/src/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection() (*gorm.DB, error) {
	fmt.Println("Connecting into database ...")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		_config.SystemParams.DB_HOST,
		_config.SystemParams.DB_USER,
		_config.SystemParams.DB_PASSWORD,
		_config.SystemParams.DB_DATABASE_NAME,
		_config.SystemParams.DB_PORT,
		_config.SystemParams.DB_TIME_ZONE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
