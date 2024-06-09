package api

import (
	"net/http"

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
		// Check for a valid authentication token
		// This is just a placeholder - you'll need to replace this with your actual authentication logic
		token := r.Header.Get("Authorization")
		if token != "valid-token" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// ProtectedHandler is a handler that requires authentication
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the request passed the AuthMiddleware
	w.Write([]byte("Hello, authenticated user!"))
}
