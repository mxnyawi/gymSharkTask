package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mxnyawi/gymSharkTask/internal/db"
)

func StartServer(dbManager *db.DBManager) {

	r := mux.NewRouter()

	// Apply the middleware to the routes that need it
	r.Use(AuthMiddleware)

	Routes(dbManager)
	http.ListenAndServe(":8080", nil)
}

// AuthMiddleware is a middleware function that checks for a valid authentication token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authentication token from config
		authToken := os.Getenv("AUTH_TOKEN")

		// Check for a valid authentication token
		token := r.Header.Get("Authorization")
		if token != authToken {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
