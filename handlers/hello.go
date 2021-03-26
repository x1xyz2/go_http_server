package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

// NewHello ...
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// ServeHttp ...
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ooops", http.StatusBadRequest)
		return
	}

	h.l.Printf("Hello::ServeHTTP called with [%s]", d)
	fmt.Fprintf(rw, "Hello %s", d)
}
