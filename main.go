package main

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-match/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Map is a concurrent-safe hash map.
	// see https://pkg.go.dev/sync#Map
	//
	// memoize users[user ID, user]
	users := new(sync.Map)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// upgrader holds websocket connection
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1 << 12,
		WriteBufferSize: 1 << 12,
	}

	// register URIs
	e.GET("/connect", router.Connect(users, upgrader))
	e.POST("/disconnect", router.Disconnect(users))
	e.POST("/recommend", router.Recommend(users))
	e.POST("/match", router.Match(users))

	// start server
	e.Logger.Fatal(e.Start(":1324"))
}
