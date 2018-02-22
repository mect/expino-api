package main

import (
	"log"
	"net/http"

	bolt "github.com/coreos/bbolt"
	"github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/cors"
)

var db *bolt.DB
var io *socketio.Server

func main() {
	var err error
	db, err = bolt.Open("backend.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	initDB()

	io, _ = socketio.NewServer(nil)
	io.On("connection", func(so socketio.Socket) {
		so.Join("update")
	})

	setTimers()
	go runTimers()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Static("/", "frontend/build")

	e.GET("/api/news", getAllNewsHandler)
	e.GET("/api/news/current", getCurrentNewsHandler)
	e.GET("/api/news/:id", getNewsHandler)
	e.POST("/api/news", addNewsHandler)
	e.PUT("/api/news", editNewsHandler)
	e.DELETE("/api/news/:id", deleteNewsHandler)

	e.PUT("/api/settings/featureslides", editFeatureSlides)
	e.GET("/api/settings/featureslides", getFeatureSlides)

	e.GET("/api/traffic", getTrafficHandler)

	e.PUT("/api/keukendienst", setKeukendienst)
	e.GET("/api/keukendienst", getKeukendienst)

	e.POST("/api/image", uploadImage)
	e.GET("/api/images/:name", getImage)

	e.GET("/api/ticker", getTickerItemsHandler)
	e.POST("/api/ticker", addTickerItemHandler)
	e.DELETE("/api/ticker/:id", deleteTickerItemsHandler)

	e.GET("/api/graphs", getGraphItemsHandler)
	e.POST("/api/graphs", addGraphItemHandler)
	e.DELETE("/api/graphs/:id", deleteGraphItemsHandler)

	http.Handle("/", e)
	http.Handle("/socket.io/", cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}).Handler(io))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
