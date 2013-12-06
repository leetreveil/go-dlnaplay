package fileserver

import (
	"net/http"
	"os"
	"time"
)

type handler struct {
	path string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(h.path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()
	http.ServeContent(w, r, h.path, time.Time{}, file)
}

func Serve(host string, path string) error {
	return http.ListenAndServe(host, &handler{path})
}
