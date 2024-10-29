package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func main() {
	router := gin.Default()

	router.GET("/user/:id", GetUser)
	router.PUT("/user/:id/create", CreateUser)
	router.PUT("/user/:id/setutcoffset/:offset", SetUserUTCOffset)
	router.PUT("/user/:id/setzip/:zip", SetUserZip)

	router.GET("/events/:id", GetUserEvents)

	router.GET("/recurringevents/:id", GetUserRecurringEvents)

	router.GET("/todo/:id", GetUserToDo)
	router.PUT("/todo/:id/toggle", ToggleToDo)
	router.GET("/todo/:id/create/:name", CreateToDo)

	router.LoadHTMLGlob("assets/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.Run("localhost:9090")
}
