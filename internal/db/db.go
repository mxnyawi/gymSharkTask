package db

import (
	"log"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/couchbase/gocb/v2"
	"github.com/joho/godotenv"
	"github.com/mxnyawi/gymSharkTask/internal/model"
)

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
		return nil, err
	}

	return &DBManager{Cluster: cluster}, nil
}

// GetDBCreds gets the database credentials
func (db *DBManager) GetDBCreds() (string, string, string, string, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return "", "", "", "", err
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
		log.Fatalf("Error loading .env file: %v", err)
		return "", "", err
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
		log.Fatalf("Error loading .env file: %v", err)
		return "", "", err
	}

	clusterUsername := os.Getenv("USERNAME")
	clusterPassword := os.Getenv("PASSWORD")

	return clusterUsername, clusterPassword, nil
}

// ConnectToCluster connects to the Couchbase cluster
func ConnectToCluster() (*gocb.Cluster, error) {
	adminUsername, adminPassword, err := GetDBAminCreds()
	if err != nil {
		log.Fatalf("Failed to get admin credentials: %v", err)
		return nil, err
	}

	cluster, err := gocb.Connect("couchbase://db", gocb.ClusterOptions{
		Username: adminUsername,
		Password: adminPassword,
	})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	return cluster, nil
}

// InitDB initializes the database
func InitDB() (*DBManager, error) {
	db, err := NewDBManager()
	if err != nil {
		return nil, err
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
		log.Fatalf("Could not get database credentials: %v", err)
		return err
	}

	err = dbManager.SetupDB(bucketName, scopeName, collectionName, documentID)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
		return err
	}

	log.Println("Bucket setup successfully")

	user, pass, err := dbManager.GetClusterCredentials()
	if err != nil {
		log.Fatalf("Could not get cluster credentials: %v", err)
		return err
	}

	hash, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		log.Fatalf("Could not hash password: %v", err)
		return err
	}

	pass = hash

	// Create a new user with full access to the bucket
	err = dbManager.WriteDocument(bucketName, scopeName, collectionName, user, User{Username: user, Password: pass})
	if err != nil {
		log.Fatalf("Failed to write user document: %v", err)
		return err
	}

	log.Println("User document written successfully")
	return nil
}
