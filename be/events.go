package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetUserEvents(c *gin.Context) {
	id, err := ParseUInt(c.Query("userId"))
	if HandleError(c, err, "Error fetching user") {
		return
	}
	var events []Event
	result := db.Where("user_id = ?", id).Find(&events)
	if HandleDBError(c, result) {
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func GetUserRecurringEvents(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Error fetching recurring event") {
		return
	}
	var events []RecurringEvent
	result := db.Where("user_id = ?", id).Find(&events)
	if HandleDBError(c, result) {
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func CreateEvent(c *gin.Context) {
	id, err := ParseUInt(c.Query("userId"))
	if HandleError(c, err, "Error fetching user") {
		return
	}
	event := Event{}
	event.UserID = id

	// date will be like &date=2024-12-31
	date := c.Query("date")
	event.Date = date

	// time will be like &time=13:00
	time := c.Query("time")
	event.Time = time

	name := c.Query("name")
	event.Name = name

	// reminders will be like:
	// &reminders=5-minutes,15-minutes etc
	remindersStr := c.Query("reminders")
	remindersSlc := strings.Split(remindersStr, ",")
	var reminders []Reminder
	for _, reminderStr := range remindersSlc {
		reminder := Reminder{}
		reminder.UserID = id
		reminder.EventID = event.ID

		splitReminder := strings.Split(reminderStr, "-")
		reminder.Unit = splitReminder[0]
		duration, err := strconv.ParseInt(splitReminder[0], 10, 64)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Error parsing reminder duration")
			return
		}
		reminder.Duration = uint(duration)
		result := db.Create(&reminder)
		if result.Error != nil {
			// TODO: do something smarter here
			log.Printf("Error creating reminder: %v", result.Error)
		}
		reminders = append(reminders, reminder)
	}
	event.Reminders = reminders

	result := db.Create(&event)
	if HandleDBError(c, result) {
		return
	}
	c.IndentedJSON(http.StatusOK, event)
}
