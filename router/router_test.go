package router_test

import (
	"localhost/medium-mongo-go-driver/handlers"
	"localhost/medium-mongo-go-driver/middlewares"
	"localhost/medium-mongo-go-driver/router"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEngine(t *testing.T) {
	// We can use same approach as with databases, so using mocked user database
	// and mocked middlewares we can get more information. But lets not repeat
	// ourselves in this case and simply just do it like this.

	dbSession := middlewares.DatabaseSession{DB: nil}

	userHandlers := handlers.User{DB: nil}
	userHandlers2020 := handlers.User2020{DB: nil}
	e := router.GetMainEngine(dbSession, userHandlers, userHandlers2020)
	assert.NotEmpty(t, e)
}
