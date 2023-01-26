package apiserver

import "net/http"

func (s *APIserver) rikishiQuery(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "HEHE", 404)
}
