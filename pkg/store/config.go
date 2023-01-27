package store

type Config struct {
	DatabaseURL string `json:"databaseurl"`
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL: "host=localhost dbname=sumodb sslmode=disable",
	}
}
