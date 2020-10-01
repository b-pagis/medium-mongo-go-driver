package databases_test

import (
	"context"
	"errors"
	"testing"

	"localhost/medium-mongo-go-driver/config"
	"localhost/medium-mongo-go-driver/databases"
	"localhost/medium-mongo-go-driver/databases/mocks"

	"localhost/medium-mongo-go-driver/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

// In order to exclude mocks there are few ways we can do it.
//
// * first - way it to put mocks into a separate package that is outside our project
// * second - way is to put mocked data into the same package where it is used
// (downside that it involves manual labour after it was generated. So if your
// code contains a lot of interfaces to mock, then you will have more manual
// labour
// * third - is to use build options and exclude mocks while using coverage (need to CONFIRM)
// * fourth - exclude mocks while running coverage
// * firth - write tests for mocks

// As mocks can be auto generated so it would be wise to either do the manual
// labour and include it into the main test functions or either to use external
// package for mocks that would not be part of project.
// or simply cheat and ignore it in the coverage ;]

func TestNewUserDatabase(t *testing.T) {
	conf := config.GetConfig()

	dbClient, err := databases.NewClient(conf)
	assert.NoError(t, err)

	db := databases.NewDatabase(conf, dbClient)

	userDB := databases.NewUserDatabase(db)

	assert.NotEmpty(t, userDB)
}

func TestFindOne(t *testing.T) {

	// Define variables for interfaces
	var dbHelper databases.DatabaseHelper
	var collectionHelper databases.CollectionHelper
	var srHelperErr databases.SingleResultHelper
	var srHelperCorrect databases.SingleResultHelper

	// Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}
	srHelperErr = &mocks.SingleResultHelper{}
	srHelperCorrect = &mocks.SingleResultHelper{}

	// Because interfaces does not implement mock.Mock functions we need to use
	// type assertion to mock implemented methods
	srHelperErr.(*mocks.SingleResultHelper).
		On("Decode", mock.AnythingOfType("*models.User")).
		Return(errors.New("mocked-error"))

	srHelperCorrect.(*mocks.SingleResultHelper).
		On("Decode", mock.AnythingOfType("*models.User")).
		Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.User)
		arg.Username = "mocked-user"
	})

	collectionHelper.(*mocks.CollectionHelper).
		On("FindOne", context.Background(), bson.M{"error": true}).
		Return(srHelperErr)

	collectionHelper.(*mocks.CollectionHelper).
		On("FindOne", context.Background(), bson.M{"error": false}).
		Return(srHelperCorrect)

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "users").Return(collectionHelper)

	// Create new database with mocked Database interface
	userDba := databases.NewUserDatabase(dbHelper)

	// Call method with defined filter, that in our mocked function returns
	// mocked-error
	user, err := userDba.FindOne(context.Background(), bson.M{"error": true})

	assert.Empty(t, user)
	assert.EqualError(t, err, "mocked-error")

	// Now call the same function with different different filter for correct
	// result
	user, err = userDba.FindOne(context.Background(), bson.M{"error": false})

	assert.Equal(t, &models.User{Username: "mocked-user"}, user)
	assert.NoError(t, err)

	// If we try following approach the cover tool will not detect it, because
	// it will think that we are testing the mock.
	// userDba = &mocks.UserDatabase{}
	// userDba.(*mocks.UserDatabase).On("FindOne", context.Background(), bson.M{"error": true}).Return(nil, errors.New("mocked-user-db-error"))
	// userDba.(*mocks.UserDatabase).On("FindOne", context.Background(), bson.M{"error": false}).Return(&models.User{Username: "mocked-db-user"}, nil)

	// user, err := userDba.FindOne(context.Background(), bson.M{"error": true})

	// assert.Empty(t, user)
	// assert.EqualError(t, err, "mocked-user-db-error")

	// user, err = userDba.FindOne(context.Background(), bson.M{"error": false})

	// assert.Equal(t, &models.User{Username: "mocked-db-user"}, user)
	// assert.NoError(t, err)

}

func TestCreate(t *testing.T) {

	// Define variables for interfaces
	var dbHelper databases.DatabaseHelper
	var collectionHelper databases.CollectionHelper

	// Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}

	collectionHelper.(*mocks.CollectionHelper).
		On("InsertOne", context.Background(), mock.AnythingOfType("*models.User")).
		Return(nil, errors.New("mocked-error"))

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "users").Return(collectionHelper)

	// Create new database with mocked Database interface
	userDba := databases.NewUserDatabase(dbHelper)

	err := userDba.Create(context.Background(), &models.User{Username: "mocked-user"})

	assert.EqualError(t, err, "mocked-error")

}

func TestDeleteByUsername(t *testing.T) {

	// Define variables for interfaces
	var dbHelper databases.DatabaseHelper
	var collectionHelper databases.CollectionHelper

	// Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}

	var errResult int64
	errResult = 0

	collectionHelper.(*mocks.CollectionHelper).
		On("DeleteOne", context.Background(), mock.AnythingOfType("*models.User")).
		Return(errResult, errors.New("mocked-error"))

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "users").Return(collectionHelper)

	// Create new database with mocked Database interface
	userDba := databases.NewUserDatabase(dbHelper)

	err := userDba.DeleteByUsername(context.Background(), "mocked-user")

	assert.EqualError(t, err, "mocked-error")

}
