package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	Date      string //format DD-MM-YYYY
	Reminders []Reminder
}

type RecurringEvent struct {
	gorm.Model
	ID        uint
	UserID    uint
	Name      string
	Frequency string
	Day       string
	Reminders []Reminder
}
type ToDo struct {
	gorm.Model
	ID        uint
	UserID    uint
	Name      string
	Completed bool
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
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Event{}, &Reminder{}, &ToDo{}, &RecurringEvent{})

	if err != nil {
		return nil, fmt.Errorf("Error connecting to database")
	}
	return db, nil
}

var db, _ = DBConn()

func parseId(idStr string) uint {
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting id to UINT in GetUser")
	}
	id := uint(id64)
	return id
}

/*
GET Methods
*/
func GetUser(c *gin.Context) {
	id := parseId(c.Param("id"))
	usr := User{}
	result := db.First(&usr, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusNoContent, "")
		return
	}
	c.IndentedJSON(http.StatusOK, usr)
}

func GetUserEvents(c *gin.Context) {
	id := parseId(c.Param("id"))
	var events []Event
	db.First(&events, "user_id = ?", id)
	c.IndentedJSON(http.StatusOK, events)
}

func GetUserRecurringEvents(c *gin.Context) {
	id := parseId(c.Param("id"))
	var events []RecurringEvent
	db.First(&events, "user_id = ?", id)
	c.IndentedJSON(http.StatusOK, events)
}

func GetUserToDo(c *gin.Context) {
	id := parseId(c.Param("id"))
	var todo []ToDo
	db.First(&todo, "id = ?", id)
	c.IndentedJSON(http.StatusOK, todo)
}

/*
Setter methods
*/
func isValidZip(zip uint) bool {
	return true
}
func parseZip(zipStr string) (uint, error) {
	zip64, err := strconv.ParseUint(zipStr, 10, 64)
	if err != nil {
		return 0, err
	}
	zip := uint(zip64)
	return zip, nil
}

func SetUserZip(c *gin.Context) {
	zip, err := parseZip(c.Param("zip"))
	id := parseId(c.Param("id"))
	println(zip)
	println(id)
	if err != nil {
		fmt.Println("Error parsing ZIP")
	}

	usr := User{}
	db.First(&usr, id)
	usr.ZIP = zip
	db.Save(&usr)
	c.IndentedJSON(http.StatusOK, true)
}

func extractint(str string) (int, error) {
	num64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 9999, err
	}
	num := int(num64)
	return num, nil
}

func SetUserUTCOffset(c *gin.Context) {
	offsetStr := c.Param("offset")
	offset, err := extractint(offsetStr)

	idStr := c.Param("id")
	id, err := extractint(idStr)
	if err != nil {
		println(err)
	}

	usr := User{}
	db.First(&usr, id)
	usr.UTCOffset = offset
	db.Save(&usr)

}
