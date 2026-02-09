package app

import (
	"database/sql"
	"errors"
	"os"

	// Обязательный импорт для регистрации драйвера "postgres"
	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("No database url were provided")
	}

	// 1. Открываем sql.DB для миграций
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// 2. Выполняем миграции
	if err := goose.Up(sqlDB, "./migrations"); err != nil {
		sqlDB.Close() // Закрываем при ошибке
		return nil, err
	}

	// Закрываем временное соединение для миграций, так как GORM откроет свое
	sqlDB.Close()

	// 3. Инициализируем GORM
	gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
