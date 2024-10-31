package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserToDo(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Error fetching todo") {
		return
	}
	var todo []ToDo
	result := db.Where("user_id = ?", id).Find(&todo)
	if HandleDBError(c, result) {
		return
	}
	c.IndentedJSON(http.StatusOK, todo)
}

func CreateToDo(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Invalid user ID") {
		return
	}
	name := c.Query("name")
	todo := ToDo{UserID: id, Name: name}
	res := db.Create(&todo)
	if HandleDBError(c, res) {
		return
	}
	c.IndentedJSON(http.StatusOK, todo)
}

func ToggleToDo(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Invalid todo ID") {
		return
	}
	todo := ToDo{}
	result := db.First(&todo, id)
	if HandleDBError(c, result) {
		return
	}
	todo.Completed = !todo.Completed
	db.Save(&todo)
	c.IndentedJSON(http.StatusOK, todo)
}
