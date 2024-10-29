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
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Event{}, &Reminder{}, &ToDo{}, &RecurringEvent{})

	if err != nil {
		return nil, fmt.Errorf("Error connecting to database")
	}
	return db, nil
}

var db, _ = DBConn()

func parseUInt(idStr string) (uint, error) {
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	id := uint(id64)
	return id, nil
}

// TODO: implement this
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

/*
GET Methods
*/
func GetUser(c *gin.Context) {
	id, err := parseUInt(c.Param("id"))
	if err != nil {
		c.IndentedJSON(500, "Error fetching user")
	}
	usr := User{}
	result := db.First(&usr, id)
	if result.Error != nil {
		c.IndentedJSON(500, "Database Error")
		return
	}
	c.IndentedJSON(http.StatusOK, usr)
}

func GetUserEvents(c *gin.Context) {
	id, err := parseUInt(c.Param("id"))
	if err != nil {
		c.IndentedJSON(500, "Error fetching user\n"+err.Error())
		return
	}
	var events []Event
	result := db.First(&events, "user_id = ?", id)
	if result.Error != nil {
		c.IndentedJSON(500, "Database Error")
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func GetUserRecurringEvents(c *gin.Context) {
	id, err := parseUInt(c.Param("id"))
	if err != nil {
		c.IndentedJSON(500, "Error fetching recurring event:\n"+err.Error())
		return
	}
	var events []RecurringEvent
	result := db.First(&events, "user_id = ?", id)
	if result.Error != nil {
		c.IndentedJSON(500, "Database Error")
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func GetUserToDo(c *gin.Context) {
	id, err := parseUInt(c.Param("id"))
	if err != nil {
		c.IndentedJSON(500, "Error fetching todo:\n"+err.Error())
		return
	}
	var todo []ToDo
	result := db.First(&todo, "id = ?", id)
	if result.Error != nil {
		c.IndentedJSON(500, "Database Error")
		return
	}
	c.IndentedJSON(http.StatusOK, todo)
}

/*
PUT methods
*/
func UpdateUser(c *gin.Context) bool {

}

func SetUserZip(c *gin.Context) {
	zip, err := parseZip(c.Param("zip"))
	if err != nil {
		c.IndentedJSON(500, "Invalid ZIP:\n"+err.Error())
	}
	id, err := parseUInt(c.Param("id"))
	if err != nil {
		c.IndentedJSON(500, "Error fetching user\n"+err.Error())
	}
	b
	usr := User{}
	db.First(&usr, id)
	usr.ZIP = zip
	db.Save(&usr)
	c.IndentedJSON(http.StatusOK, true)
}

func extractInt(str string) (int, error) {
	num64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 9999, err
	}
	num := int(num64)
	return num, nil
}

func SetUserUTCOffset(c *gin.Context) {
	offsetStr := c.Param("offset")
	offset, err := extractInt(offsetStr)

	idStr := c.Param("id")
	id, err := extractInt(idStr)
	if err != nil {
		println(err)
	}

	usr := User{}
	db.First(&usr, id)
	usr.UTCOffset = offset
	db.Save(&usr)

	c.IndentedJSON(http.StatusOK, true)
}

func CreateUser(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	usr := User{ID: id}
	db.Create(&usr)
	c.IndentedJSON(http.StatusOK, usr)
}

func CreateToDo(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	name := c.Param("name")
	todo := ToDo{UserID: id, Name: name}
	db.Create(&todo)
	c.IndentedJSON(http.StatusOK, todo)
}

func ToggleToDo(c *gin.Context) {
	id := parseUInt(c.Param("id"))
	todo := ToDo{}
	db.First(&todo, id)
	todo.Completed = !todo.Completed
	db.Save(&todo)
	c.IndentedJSON(http.StatusOK, todo)
}
