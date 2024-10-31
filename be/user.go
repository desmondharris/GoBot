package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUserZip(c *gin.Context) {
	zip, err := ParseZip(c.Query("zip"))
	if HandleError(c, err, "Invalid ZIP") {
		return
	}
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Error fetching user") {
		return
	}
	usr := User{}
	db.First(&usr, id)
	usr.ZIP = zip
	db.Save(&usr)
	c.IndentedJSON(http.StatusOK, true)
}

func GetUser(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Error fetching user") {
		return
	}
	usr := User{}
	result := db.First(&usr, id)
	if HandleDBError(c, result) {
		return
	}
	c.IndentedJSON(http.StatusOK, usr)
}

func UpdateUser(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Error parsing userid") {
		return
	}
	usr := User{}
	result := db.First(&usr, id)
	if HandleDBError(c, result) {
		return
	}

	zip := c.DefaultQuery("zip", "")
	if zip != "" {
		zipInt, err := ParseZip(zip)
		if HandleError(c, err, "Invalid ZIP") {
			return
		}
		usr.ZIP = zipInt
		result = db.Save(&usr)
		if HandleDBError(c, result) {
			return
		}

	}

	utcOffset := c.DefaultQuery("utcoffset", "")
	if utcOffset != "" {
		offsetInt, err := ExtractInt(utcOffset)
		if HandleError(c, err, "Invalid offset") {
			return
		}
		usr.UTCOffset = offsetInt
		result = db.Save(&usr)
		if HandleDBError(c, result) {
			return
		}
	}

}
func SetUserUTCOffset(c *gin.Context) {
	offsetStr := c.Query("offset")
	offset, err := ExtractInt(offsetStr)
	if HandleError(c, err, "Invalid offset") {
		return
	}
	idStr := c.Query("id")
	id, err := ExtractInt(idStr)
	if HandleError(c, err, "Error fetching user") {
		return
	}
	usr := User{}
	db.First(&usr, id)
	usr.UTCOffset = offset
	db.Save(&usr)
	c.IndentedJSON(http.StatusOK, true)
}

func CreateUser(c *gin.Context) {
	id, err := ParseUInt(c.Query("id"))
	if HandleError(c, err, "Invalid user ID") {
		return
	}
	usr := User{ID: id}
	db.Create(&usr)
	c.IndentedJSON(http.StatusOK, usr)
}
