package router

import (
	"localhost/medium-mongo-go-driver/handlers"
	"localhost/medium-mongo-go-driver/middlewares"

	"github.com/gin-gonic/gin"
)

// GetMainEngine get main engine with routers
func GetMainEngine(dbSession middlewares.DatabaseSession, usersHandler handlers.User, userHandlers2020 handlers.User2020) *gin.Engine {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(dbSession.SetDbContext) // This becomes useless if userHandlers2020 are used

	v1 := r.Group("/api")
	{

		v1.GET("/users/:username", usersHandler.GetUser)
		v1.POST("/users", usersHandler.CreateUser)
		v1.DELETE("/users/:username", usersHandler.DeleteUser)

	}
	v2 := r.Group("/apiv2")
	{

		v2.GET("/users/:username", userHandlers2020.GetUser)
		v2.POST("/users", userHandlers2020.CreateUser)
		v2.DELETE("/users/:username", userHandlers2020.DeleteUser)

	}
	return r

}
