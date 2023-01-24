package apiserver

type Config struct {
	Port string `json:"port"`
	//loggerconfig
}

func NewConfig() *Config {
	return &Config{
		Port: ":8080",
	}
}
