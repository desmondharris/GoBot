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
	router.GET("/events/:id", GetUserEvents)
	router.GET("/recurringevents/:id", GetUserRecurringEvents)
	router.GET("/todo/:id", GetUserToDo)

	router.PUT("/user/:id/setzip/:zip", SetUserZip)
	router.PUT("/user/:id/setutcoffset/:offset", SetUserUTCOffset)
	router.Run("localhost:9090")
}
