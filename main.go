package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/x1xyz2/go_http_server/handlers"
)

func main() {
	l := log.New(os.Stdout, "#product-api# ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	r := mux.NewRouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.ValidateMiddleware)

	addRouter := r.Methods(http.MethodPost).Subrouter()
	addRouter.HandleFunc("/", ph.AddProduct)
	addRouter.Use(ph.ValidateMiddleware)

	delRouter := r.Methods(http.MethodDelete).Subrouter()
	delRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	s := &http.Server{
		Addr:         ":9999",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig.String())

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
