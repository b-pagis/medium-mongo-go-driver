package config_test

import (
	"localhost/medium-mongo-go-driver/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	conf := config.GetConfig()

	assert.NotEmpty(t, conf)
}
