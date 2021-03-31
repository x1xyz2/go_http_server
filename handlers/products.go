package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/x1xyz2/go_http_server/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		//p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		//p.addProducts(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		//url := r.URL.Path
		//p.l.Println(url)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// GetProducts blablabla
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle getProducts")
	lp := data.GetProducts()
	//d, err := json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle addProducts")

	/*
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to un-marshal json", http.StatusInternalServerError)
		}
	*/

	ctx := r.Context()
	prod := ctx.Value(KeyProduct{}).(data.Product)
	data.AddProducts(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to Atoi", http.StatusInternalServerError)
		return
	}

	p.l.Println("Handle updateProduct", id)

	/*
		prod := &data.Product{}
		err = prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to un-marshal json", http.StatusInternalServerError)
			return
		}
	*/

	ctx := r.Context()
	prod := ctx.Value(KeyProduct{}).(data.Product)

	data.UpdateProduct(id, &prod)
}

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to Atoi", http.StatusInternalServerError)
		return
	}

	p.l.Println("Handle delProduct", id)

	data.DeleteProduct(id)
}

type KeyProduct struct{}

func (p *Products) ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to un-marshal json", http.StatusInternalServerError)
			p.l.Println("Unable to un-marshal json,", http.StatusInternalServerError, "Error")
			return
		}

		err = prod.Validate()
		if err != nil {
			errstr := fmt.Sprintf("Json validate failed: %s", err)
			http.Error(rw, errstr, http.StatusBadRequest)
			p.l.Println(errstr, http.StatusBadRequest, "Error")
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
