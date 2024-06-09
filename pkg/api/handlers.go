// pkg/api/handlers.go

package api

import (
	"encoding/json"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/mxnyawi/gymSharkTask/internal/db"
	"github.com/mxnyawi/gymSharkTask/internal/model"
)

type OrderRequest struct {
	OrderAmount  int   `json:"orderAmount"`
	PackageSizes []int `json:"packageSizes"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DocumentRequest struct {
	Document db.Document `json:"document"`
}

type DBSetuo struct {
	BucketName     string `json:"bucketName"`
	ScopeName      string `json:"scopeName"`
	CollectionName string `json:"collectionName"`
	DocumentID     string `json:"documentID"`
}

var packages = &model.Packages{Sizes: []int{1, 5, 10}} // Default package sizes

func RegisterHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	user.Password = hash

	bucketName, scopeName, collectionName, _, err := dbManager.GetDBCreds()
	if err != nil {
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Store the user in the database
	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, user.Username, user)
	if err != nil {
		http.Error(w, "Could not store user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucketName, scopeName, collectionName, _, err := dbManager.GetDBCreds()
	if err != nil {
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Retrieve the user from the database
	storedUser, err := dbManager.GetUser(bucketName, scopeName, collectionName, user.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Compare the provided password with the stored hashed password
	match, err := argon2id.ComparePasswordAndHash(user.Password, storedUser.Password)
	if err != nil {
		http.Error(w, "Error while comparing password and hash", http.StatusInternalServerError)
		return
	}

	if !match {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// User is authenticated
	w.WriteHeader(http.StatusOK)
}

func PostOrderHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	var req OrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Fetch the existing document
	existingDocument, err := dbManager.GetDocument(bucketName, scopeName, collectionName, documentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the new document to the history
	existingDocument.History = append(existingDocument.History, document)

	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, documentID, existingDocument)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func CreateAdminUserHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dbManager.CreateAdminUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func SetDocumentHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	var req DocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, documentID, req.Document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetDocumentHandler retrieves a document from the database
func GetDocumentHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	// Call the GetDocument method
	content, err := dbManager.GetDocument(bucketName, scopeName, collectionName, documentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the content to JSON and write it to the response
	jsonContent, err := json.Marshal(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonContent)
}

// UpdateDocumentHandler updates a document in the database
func UpdateDocumentHandler(w http.ResponseWriter, r *http.Request, dbManager *db.DBManager) {
	var req DocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mapContent := map[string]interface{}{
		"order":    req.Document.Order,
		"packages": req.Document.Packages,
	}

	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		http.Error(w, "Could not get database credentials", http.StatusInternalServerError)
		return
	}

	err = dbManager.UpdateDocument(bucketName, scopeName, collectionName, documentID, mapContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
