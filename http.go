package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/damoye/moproxy/backend"
)

func serveHTTP(address string, m *backend.Manager) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Write(m.Describe())
	})

	http.HandleFunc("/backend", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		var body struct {
			Address string `json:"address"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.Add(body.Address)
	})

	log.Fatal(http.ListenAndServe(address, nil))
}
