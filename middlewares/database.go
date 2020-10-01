package middlewares

import (
	"localhost/medium-mongo-go-driver/databases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// ALL MIDDLEWARES PACKAGE CAN BE REMOVED IF user2020 HANDLERS ARE USED

const dbContextName = "db"

type DatabaseSession struct {
	DB databases.DatabaseHelper
}

func GetDbSessionContext(c *gin.Context) mongo.SessionContext {
	return c.MustGet(dbContextName).(mongo.SessionContext)
}

func (d DatabaseSession) SetDbContext(c *gin.Context) {
	session, err := d.DB.Client().StartSession()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	mongo.WithSession(c, session, func(sessionContext mongo.SessionContext) error {
		c.Set(dbContextName, sessionContext)
		c.Next()
		return nil
	})

}
