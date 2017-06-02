package proxy

import (
	"encoding/json"
	"net/http"
)

func (proxy *Proxy) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := w.Write(proxy.manager.Describe()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (proxy *Proxy) postBackend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Address string `json:"address"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	proxy.manager.Add(body.Address)
}

func (proxy *Proxy) serveHTTP() {
	http.HandleFunc("/", proxy.index)
	http.HandleFunc("/backend", proxy.postBackend)
	panic(http.ListenAndServe(proxy.config.HTTPAddress, nil))
}
