package handlers

import (
	"localhost/medium-mongo-go-driver/databases"
	"localhost/medium-mongo-go-driver/middlewares"
	"localhost/medium-mongo-go-driver/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	DB databases.UserDatabase
}

func (u User) GetUser(c *gin.Context) {
	username := c.Param("username")

	sessionContext := middlewares.GetDbSessionContext(c)

	filter := bson.M{
		"username": username,
	}

	user, err := u.DB.FindOne(sessionContext, filter)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, user)

}

func (u User) CreateUser(c *gin.Context) {
	sessionContext := middlewares.GetDbSessionContext(c)

	requestParams := &struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{}

	err := c.BindJSON(&requestParams)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

	err = u.DB.Create(sessionContext, &models.User{
		Username: requestParams.Username,
		Email:    requestParams.Email,
		ID:       primitive.NewObjectID(),
	})

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

	c.AbortWithStatusJSON(200, gin.H{})

}

func (u User) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	sessionContext := middlewares.GetDbSessionContext(c)

	err := u.DB.DeleteByUsername(sessionContext, username)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

}
