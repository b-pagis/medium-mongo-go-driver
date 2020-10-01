package config

// Config for database
type Config struct {
	Username     string
	Password     string
	DatabaseName string
	URL          string
}

// GetConfig returns hardcoded config
func GetConfig() *Config {
	return &Config{
		Username:     "user.name",
		Password:     "pass.word",
		DatabaseName: "database.db",
		URL:          "mongodb://127.0.0.1:27017",
	}
}
