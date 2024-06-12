package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mxnyawi/gymSharkTask/internal/db"
	"github.com/rs/cors"
)

func Routes(dbManager db.DBManagerInterface) {
	r := mux.NewRouter()

	// Middleware to authenticate users
	r.Use(AuthMiddleware)

	// User management routes
	r.HandleFunc("/registerUser", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r, dbManager)
	}).Methods("POST")

	r.HandleFunc("/loginUser", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, dbManager)
	}).Methods("POST")

	r.HandleFunc("/createAdminUser", func(w http.ResponseWriter, r *http.Request) {
		CreateAdminUserHandler(w, r, dbManager)
	}).Methods("POST")

	// Order management route
	r.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		PostOrderHandler(w, r, dbManager)
	}).Methods("POST")

	// Document management routes
	r.HandleFunc("/setDocument", func(w http.ResponseWriter, r *http.Request) {
		SetDocumentHandler(w, r, dbManager)
	}).Methods("POST")

	r.HandleFunc("/getDocument", func(w http.ResponseWriter, r *http.Request) {
		GetDocumentHandler(w, r, dbManager)
	}).Methods("GET")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	http.Handle("/", handler)
}
