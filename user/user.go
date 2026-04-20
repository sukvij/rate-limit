package user

import (
	"rate-limiter/middleware"
	"sync"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id   int
	Name string
}

func UserController(app *gin.Engine, Mp map[string]*middleware.Node, mtx *sync.Mutex) {
	mw := &middleware.Middleware{Mp: Mp, Lock: mtx}
	app.POST("/", mw.RateLimitCheck, createUser)
	app.GET("/stats", mw.UserStatus)
}

func createUser(ctx *gin.Context) {
	ctx.JSON(200, "Rate limit Nahi h..")
}
