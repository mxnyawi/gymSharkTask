// pkg/api/handlers.go

package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/mxnyawi/gymSharkTask/internal/db"
	"github.com/mxnyawi/gymSharkTask/internal/model"
)

// OrderRequest is a struct that contains the order amount and package sizes
type OrderRequest struct {
	OrderAmount  int   `json:"orderAmount"`
	PackageSizes []int `json:"packageSizes"`
}

// UserRequest is a struct that contains the user credentials
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DocumentRequest is a struct that contains the document
type DocumentRequest struct {
	Document db.Document `json:"document"`
}

var packages = &model.Packages{Sizes: []int{1, 5, 10}} // Default package sizes

// RegisterHandler registers a new user
func RegisterHandler(w http.ResponseWriter, r *http.Request, dbManager db.DBManagerInterface) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	user.Password = hash

	bucketName, scopeName, collectionName, _, err := dbManager.GetDBCreds()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Store the user in the database
	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, user.Username, user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not store user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

// LoginHandler logs in a user and checks the password
func LoginHandler(w http.ResponseWriter, r *http.Request, dbManager db.DBManagerInterface) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	bucketName, scopeName, collectionName, _, err := dbManager.GetDBCreds()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Retrieve the user from the database
	storedUser, err := dbManager.GetUser(bucketName, scopeName, collectionName, user.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Compare the provided password with the stored hashed password
	match, err := argon2id.ComparePasswordAndHash(user.Password, storedUser.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error while comparing password and hash", http.StatusInternalServerError)
		return
	}

	if !match {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// User is authenticated
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User authenticated"})
}

// PostOrderHandler creates a new order and finds the packages for it
func PostOrderHandler(w http.ResponseWriter, r *http.Request, dbManager db.DBManagerInterface) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req OrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order := &model.Order{Amount: req.OrderAmount}
	packages.Sizes = req.PackageSizes

	packages.FindPackages(order, packages)

	// Create a document with the order and packages
	document := db.Document{
		Order:    *order,
		Packages: *packages,
	}

	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Fetch the existing document
	existingDocument, err := dbManager.GetDocument(bucketName, scopeName, collectionName, documentID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not find order history", http.StatusInternalServerError)
		return
	}

	// Add the new document to the history
	existingDocument.History = append(existingDocument.History, document)

	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, documentID, existingDocument)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not write order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// CreateAdminUserHandler creates an admin user in the database
func CreateAdminUserHandler(w http.ResponseWriter, r *http.Request, dbManager db.DBManagerInterface) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	err = dbManager.CreateAdminUser(req.Username, req.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin user created"})
}

// SetDocumentHandler sets a document in the database
func SetDocumentHandler(w http.ResponseWriter, r *http.Request, dbManager db.DBManagerInterface) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req DocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, documentID, req.Document)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not write document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Document created"})
}

// GetDocumentHandler retrieves a document from the database
func GetDocumentHandler(w http.ResponseWriter, r *http.Request, dbManager db.DBManagerInterface) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Call the GetDocument method
	content, err := dbManager.GetDocument(bucketName, scopeName, collectionName, documentID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not get document", http.StatusInternalServerError)
		return
	}

	// Convert the content to JSON and write it to the response
	jsonContent, err := json.Marshal(content)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not marshall response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonContent)
}
