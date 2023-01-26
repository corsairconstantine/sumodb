package apiserver

import (
	"log"
	"net/http"

	"github.com/corsairconstantine/http-rest-api/app/store"
)

type APIserver struct {
	config *Config
	logger *log.Logger
	mux    *http.ServeMux
	store  *store.Store
}

// create a new api server
func new(serverConfig *Config, storeConfig *store.Config) *APIserver {
	s := &APIserver{
		config: serverConfig,
		logger: log.Default(),
		mux:    http.NewServeMux(),
		store:  store.Open(storeConfig),
	}
	rHandler := &rikishiHandler{s.logger, s.store.Db}
	s.mux.Handle("/rikishis/", rHandler)
	rsHandler := &rikishisHandler{s.logger, s.store.Db}
	s.mux.Handle("/rikishis", rsHandler)

	return s
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
	s.mux.ServeHTTP(w, r)
}
