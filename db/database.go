package db

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

// var err error

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(`Error loading .env file: `, err.Error())
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
	)
	db, err = sqlx.Connect("pgx", connStr)
	if err != nil {
		panic(err.Error())
	}
	db.DB.SetMaxOpenConns(100)
	db.DB.SetConnMaxIdleTime(time.Second * 5)
	db.DB.SetConnMaxLifetime(time.Hour)
	db.DB.SetMaxIdleConns(10)
}
func CreateConn() *sqlx.DB {
	if db != nil {
		fmt.Println("Connected to database")
		return db
	} else {
		fmt.Println("DB is not initialized.")
		return nil
	}
}
