package handlers

import (
	"context"
	"localhost/medium-mongo-go-driver/databases"
	"localhost/medium-mongo-go-driver/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User2020 struct {
	// local abstraction layer so it would be possible to define only those methods
	// that are needed in this scope. alternatively if coupling is not a problem then
	// it is possible to use local abstraction layer like
	// 	type repository interface {
	// 	  FindOne(context.Context, interface{}) (*models.User, error)
	// 	  Create(context.Context, *models.User) error
	// 	  DeleteByUsername(context.Context, string) error
	//  }
	// and models.User could be divided into one used for database and one used for API
	// but this would require more tweaking around the old code.
	DB databases.UserDatabase
}

func (u User2020) GetUser(c *gin.Context) {
	username := c.Param("username")

	filter := bson.M{
		"username": username,
	}
	user, err := u.DB.FindOne(context.TODO(), filter)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

	c.SecureJSON(200, user)

}

func (u User2020) CreateUser(c *gin.Context) {

	requestParams := &struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{}

	err := c.BindJSON(&requestParams)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

	err = u.DB.Create(context.TODO(), &models.User{
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

func (u User2020) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	err := u.DB.DeleteByUsername(context.TODO(), username)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
		return
	}

}
