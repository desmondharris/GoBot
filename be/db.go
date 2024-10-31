package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	ID        uint
	ZIP       uint
	UTCOffset int
	CreatedAt time.Time
	UpdatedAt time.Time
	Events    []Event
	Reminders []Reminder
	ToDos     []ToDo
}

type Event struct {
	gorm.Model
	ID        uint
	UserID    uint
	Name      string
	Date      string // format DD-MM-YYYY
	Time      string // format HH:MM
	Reminders []Reminder
}

type RecurringEvent struct {
	gorm.Model
	ID        uint
	UserID    uint
	Name      string
	Frequency string
	Day       string
	Time      string
	Reminders []Reminder
}

type ToDo struct {
	gorm.Model
	ID        uint
	UserID    uint
	Name      string
	Completed bool
	UpdatedAt time.Time
}

type Reminder struct {
	gorm.Model
	ID       uint
	EventID  uint
	UserID   uint
	Unit     string
	Duration uint
}

func DBConn() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	user := os.Getenv("MYSQLT_USER")
	pass := os.Getenv("MYSQLT_PASSWORD")
	if user == "" || pass == "" {
		return nil, fmt.Errorf("MYSQLT_USER or MYSQLT_PASSWORD not set")
	}

	var mysqlConfig = user + ":" + pass + "@tcp(127.0.0.1:3306)/bottest?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(mysqlConfig), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database")
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Event{}, &Reminder{}, &ToDo{}, &RecurringEvent{})

	return db, nil
}

var db, _ = DBConn()

func ParseUInt(idStr string) (uint, error) {
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	id := uint(id64)
	return id, nil
}

func ParseZip(zipStr string) (uint, error) {
	zip64, err := strconv.ParseUint(zipStr, 10, 64)
	if err != nil {
		return 0, err
	}
	zip := uint(zip64)
	return zip, nil
}

func ExtractInt(str string) (int, error) {
	num64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 9999, err
	}
	num := int(num64)
	return num, nil
}

func HandleError(c *gin.Context, err error, message string) bool {
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, message+"\n"+err.Error())
		log.Printf("%v: %v", message, err)
		return true
	}
	return false
}

func HandleDBError(c *gin.Context, result *gorm.DB) bool {
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Database Error")
		log.Fatalf("Database Error: %v", result.Error)
		return true
	}
	return false
}
