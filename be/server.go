package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	corsWrapper "github.com/rs/cors/wrapper/gin"
	"io"
	"log"
)

func logRequestData(c *gin.Context) {
	// Log request method and URL
	log.Printf("Request: %s %s", c.Request.Method, c.Request.URL)

	// Log request headers
	for name, values := range c.Request.Header {
		for _, value := range values {
			log.Printf("Header: %s=%s", name, value)
		}
	}

	// Log request body
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			log.Printf("Body: %s", string(bodyBytes))
			// Restore the io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	// Continue to the next middleware/handler
	c.Next()
}

func main() {
	router := gin.Default()
	router.Use(corsWrapper.New(corsWrapper.Options{
		AllowedOrigins:   []string{"*"}, // You can set specific origins if needed for production.
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.Use(logRequestData)

	// user table calls
	router.GET("/user", GetUser)
	router.POST("/user/create", CreateUser)
	router.PUT("/user/update", UpdateUser)

	// events table calls
	router.GET("/events", GetUserEvents)
	router.POST("/events/create", CreateEvent)

	// recurring events table calls
	router.GET("/recurringevents", GetUserRecurringEvents)

	// to do table calls
	router.GET("/todo", GetUserToDo)
	router.PUT("/todo/toggle", ToggleToDo)
	router.GET("/todo/create", CreateToDo)

	// serve static files
	router.Static("/assets", "./assets")
	router.StaticFile("/portal", "./assets/index.html")
	router.StaticFile("/portal/events", "./assets/events.html")

	router.Run("localhost:9090")
}
