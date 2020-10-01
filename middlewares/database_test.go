package middlewares_test

import (
	"net/http/httptest"
	"testing"

	"localhost/medium-mongo-go-driver/config"
	"localhost/medium-mongo-go-driver/databases"
	"localhost/medium-mongo-go-driver/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetDbSessionContext(t *testing.T) {
	// Create gin test context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Create database context. I use real database, but it is possible to mock
	// database and configuration through interfaces.
	conf := config.GetConfig()
	client, err := databases.NewClient(conf)

	assert.NoError(t, err) // simple check for client error
	db := databases.NewDatabase(conf, client)
	dbSession := middlewares.DatabaseSession{
		DB: db,
	}
	dbSession.SetDbContext(c)

	// At first GetDbSessionContext will panic as it will not get "db" key,
	// because SetDbContext will get and error "client not connected". So we
	// set assert the function that will panic
	var noDbCtx assert.PanicTestFunc
	noDbCtx = func() { middlewares.GetDbSessionContext(c) }

	assert.Panics(t, noDbCtx)

	// Now we connect client like it is done in real life scenario and call
	// SetDbContext again.
	client.Connect()
	dbSession.SetDbContext(c)

	// Then lets retrieve the context and check if it is not empty
	mongoCtx := middlewares.GetDbSessionContext(c)
	assert.NotEmpty(t, mongoCtx)
}
