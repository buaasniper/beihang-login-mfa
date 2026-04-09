package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func InitDB() *gorm.DB {
	// 尝试加载 .env 文件（如果存在）
	_ = godotenv.Load()

	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "bot_defender")
	password := os.Getenv("DB_PASSWORD")
	dbName := getEnv("DB_NAME", "anti_bot_db")

	if password == "" {
		log.Fatal("DB_PASSWORD environment variable must be set")
	}
	if !dbNamePattern.MatchString(dbName) {
		log.Fatal("DB_NAME contains invalid characters; only letters, numbers, and underscore are allowed")
	}

	createDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		user, password, host, port)
	createDBIfNotExists(createDSN, dbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func createDBIfNotExists(dsn, dbName string) {
	rawDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect MySQL for database bootstrap: %v", err)
	}
	defer rawDB.Close()

	if err := rawDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL for database bootstrap: %v", err)
	}

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if _, err := rawDB.Exec(query); err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}
}
