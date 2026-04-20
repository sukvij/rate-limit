package main

import (
	"rate-limiter/middleware"
	"rate-limiter/user"
	"sync"

	"github.com/gin-gonic/gin"
)

// var Mp map[string]*middleware.Node

func main() {
	app := gin.Default()
	Mp := map[string]*middleware.Node{} // map[ip]Node{}
	mtx := &sync.Mutex{}
	user.UserController(app, Mp, mtx)
	app.Run(":8080")
}
