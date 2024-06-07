package main

import (
	"net/http"

	"github.com/mxnyawi/gymSharkTask/pkg/api"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/order", api.OrderHandler).Methods("POST")
	http.ListenAndServe(":8080", r)

}
