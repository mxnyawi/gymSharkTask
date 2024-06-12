package db

import (
	"fmt"
	"log"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/couchbase/gocb/v2"
	"github.com/joho/godotenv"
	"github.com/mxnyawi/gymSharkTask/internal/model"
)

type DBManagerInterface interface {
	GetDocument(bucketName, scopeName, collectionName, documentID string) (*DocumentHistory, error)
	GetUser(bucketName, scopeName, collectionName, documentID string) (*User, error)
	WriteDocument(bucket, scope, collection, id string, data interface{}) error
	GetDBCreds() (string, string, string, string, error)
	CreateAdminUser(username, password string) error
	CreateBucket(bucketName string) error
	CreateScope(bucketName, scopeName string) error
	CreateCollection(bucketName, scopeName, collectionName string) error
	GetClusterCredentials() (string, string, error)
}

// DBManager is a struct that contains the Couchbase cluster
type DBManager struct {
	Cluster *gocb.Cluster
}

// DocumentHistory is a struct that contains the history of documents
type DocumentHistory struct {
	History []Document `json:"history"`
}

// Document is a struct that contains the packages and order
type Document struct {
	Packages model.Packages `json:"packages"`
	Order    model.Order    `json:"order"`
}

// User is a struct that contains the user credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewDBManager creates a new DBManager
func NewDBManager() (*DBManager, error) {
	cluster, err := ConnectToCluster()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}

	return &DBManager{Cluster: cluster}, nil
}

// GetDBCreds gets the database credentials
func (db *DBManager) GetDBCreds() (string, string, string, string, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to load .env file: %w", err)
	}

	bucketName := os.Getenv("BUCKET_NAME")
	scopeName := os.Getenv("SCOPE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	documentID := os.Getenv("DOCUMENT_ID")

	return bucketName, scopeName, collectionName, documentID, nil
}

// GetDBAminCreds gets the database admin credentials
func GetDBAminCreds() (string, string, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		return "", "", fmt.Errorf("failed to load .env file: %w", err)
	}

	// Get the admin credentials
	adminUsername := os.Getenv("USERNAME")
	adminPassword := os.Getenv("PASSWORD")

	return adminUsername, adminPassword, nil
}

// GetClusterCredentials gets the cluster credentials
func (db *DBManager) GetClusterCredentials() (string, string, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		return "", "", fmt.Errorf("failed to load .env file: %w", err)
	}

	clusterUsername := os.Getenv("USERNAME")
	clusterPassword := os.Getenv("PASSWORD")

	return clusterUsername, clusterPassword, nil
}

// ConnectToCluster connects to the Couchbase cluster
func ConnectToCluster() (*gocb.Cluster, error) {
	adminUsername, adminPassword, err := GetDBAminCreds()
	if err != nil {
		return nil, fmt.Errorf("failed to get admin credentials: %w", err)
	}

	cluster, err := gocb.Connect("couchbase://db", gocb.ClusterOptions{
		Username: adminUsername,
		Password: adminPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return cluster, nil
}

// InitDB initializes the database
func InitDB() (*DBManager, error) {
	db, err := NewDBManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create new DBManager: %w", err)
	}

	log.Println("Connected to Couchbase successfully")

	// Create the bucket
	SetupBucket(db)

	log.Println("Database setup successfully")
	return db, nil
}

// CreateBucketHandler creates a new bucket in the database
func SetupBucket(dbManager *DBManager) error {
	bucketName, scopeName, collectionName, documentID, err := dbManager.GetDBCreds()
	if err != nil {
		return fmt.Errorf("failed to get database credentials: %w", err)
	}

	err = dbManager.SetupDB(bucketName, scopeName, collectionName, documentID)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	log.Println("Bucket setup successfully")

	user, pass, err := dbManager.GetClusterCredentials()
	if err != nil {
		return fmt.Errorf("failed to get cluster credentials: %w", err)
	}

	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	pass = hash

	// Create a new user with full access to the bucket
	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, user, User{Username: user, Password: pass})
	if err != nil {
		return fmt.Errorf("failed to write user document: %w", err)
	}

	log.Println("User document written successfully")
	return nil
}
