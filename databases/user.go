package databases

import (
	"context"

	"localhost/medium-mongo-go-driver/models"
)

const collectionName = "users"

// UserDatabase user database representation to find, update, delete users
type UserDatabase interface {
	FindOne(context.Context, interface{}) (*models.User, error)
	Create(context.Context, *models.User) error
	DeleteByUsername(context.Context, string) error
}

type userDatabase struct {
	db DatabaseHelper
}

// NewUserDatabase creates new user database instance
func NewUserDatabase(db DatabaseHelper) UserDatabase {
	return &userDatabase{
		db: db,
	}
}

// FindOne finds single record by passed filter
func (u *userDatabase) FindOne(ctx context.Context, filter interface{}) (*models.User, error) {
	user := &models.User{}
	err := u.db.Collection(collectionName).FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Create creates new user in the database
func (u *userDatabase) Create(ctx context.Context, usr *models.User) error {
	_, err := u.db.Collection(collectionName).InsertOne(ctx, usr)
	return err
}

// DeleteByUsername deletes user by provided username
func (u *userDatabase) DeleteByUsername(ctx context.Context, username string) error {
	// In this case it is possible to use bson.M{"username":username} but I tend
	// to avoid another dependency in this layer and for demonstration purposes
	// used omitempty in the model
	user := &models.User{
		Username: username,
	}
	_, err := u.db.Collection(collectionName).DeleteOne(ctx, user)
	return err
}
