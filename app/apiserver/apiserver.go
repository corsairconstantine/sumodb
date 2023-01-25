package apiserver

import (
	"log"
	"net/http"

	"github.com/corsairconstantine/http-rest-api/app/store"
)

type rikishi struct {
	ID      int    `json:"id"`
	Shikona string `json:"shikona"`
	Rank    string `json:"rank"`
	Height  uint8  `json:"height"`
	Weight  uint8  `json:"weight"`
}

type bout struct {
	ID         int     `json:"id"`
	Winner     rikishi `json:"winner"`
	Loser      rikishi `json:"loser"`
	Tournament string  `json:"tournament"`
	Division   string  `json:"division"`
	Day        uint8   `json:"day"`
}

type APIserver struct {
	config *Config
	logger *log.Logger
	mux    *http.ServeMux
	store  *store.Store
}

// create a new api server
func new(serverConfig *Config, storeConfig *store.Config) *APIserver {
	return &APIserver{
		config: serverConfig,
		logger: log.Default(),
		mux:    http.NewServeMux(),
		store:  store.Open(storeConfig),
	}
}

// Start a server
func Start() error {
	serverConfig := NewConfig()
	storeConfig := store.NewConfig()

	server := new(serverConfig, storeConfig)
	server.logger.Println("starting a new server")

	err := http.ListenAndServe(serverConfig.Port, server)
	if err != nil {
		server.logger.Fatal(err)
	}
	return err
}

func (s *APIserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/bouts":
		s.boutsHandler(w, r)
	case "/rikishis":
		s.rikishiHandler(w, r)
	}
}
