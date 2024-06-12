package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
)

// CreateAdminUser creates an admin user
func (db *DBManager) CreateAdminUser(username, password string) error {
	userSettings := gocb.User{
		Username: username,
		Password: password,
		Roles:    []gocb.Role{{Name: "admin", Bucket: ""}},
	}

	err := db.Cluster.Users().UpsertUser(userSettings, nil)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println("Admin user created successfully")
	return nil
}

// SetupDB sets up the database
func (db *DBManager) SetupDB(bucketName, scopeName, collectionName, documentID string) error {
	err := db.CreateBucket(bucketName)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	err = db.CreateScope(bucketName, scopeName)
	if err != nil {
		return fmt.Errorf("failed to create scope: %w", err)
	}

	err = db.CreateCollection(bucketName, scopeName, collectionName)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	_, err = db.GetDocument(bucketName, scopeName, collectionName, documentID)
	if err == nil {
		log.Println("Document already exists")
		return nil
	}

	err = db.WriteDocument(bucketName, scopeName, collectionName, documentID, DocumentHistory{History: []Document{}})
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}

	return nil
}

// CreateBucket creates a bucket
func (db *DBManager) CreateBucket(bucketName string) error {
	createBucket := gocb.CreateBucketSettings{
		BucketSettings: gocb.BucketSettings{
			Name:       bucketName,
			BucketType: gocb.CouchbaseBucketType,
			RAMQuotaMB: 100,
		},
	}

	err := db.Cluster.Buckets().CreateBucket(createBucket, nil)
	switch {
	case err == nil:
		log.Println("Bucket created successfully")
		return nil
	case strings.Contains(err.Error(), "Bucket with given name already exists"):
		log.Println("Bucket already exists")
		return nil
	default:
		return fmt.Errorf("failed to create bucket: %w", err)
	}
}

// CreateScope creates a scope in a bucket
func (db *DBManager) CreateScope(bucketName, scopeName string) error {
	err := db.Cluster.Bucket(bucketName).Collections().CreateScope(scopeName, nil)
	switch {
	case err == nil:
		log.Println("Scope created successfully")
		return nil
	case strings.Contains(err.Error(), "already exists"):
		log.Println("Scope already exists")
		return nil
	default:
		return fmt.Errorf("failed to create scope: %w", err)
	}
}

// CreateCollection creates a collection in a scope
func (db *DBManager) CreateCollection(bucketName, scopeName, collectionName string) error {

	collectionNameSettings := gocb.CollectionSpec{
		Name:      collectionName,
		ScopeName: scopeName,
	}

	err := db.Cluster.Bucket(bucketName).Collections().CreateCollection(collectionNameSettings, nil)
	switch {
	case err == nil:
		log.Println("Collection created successfully")
		return nil
	case strings.Contains(err.Error(), "already exists"):
		log.Println("Collection already exists")
		return nil
	default:
		return fmt.Errorf("failed to create collection: %w", err)
	}
}

// WriteDocument writes a document to the database collection
func (db *DBManager) WriteDocument(bucketName, scopeName, collectionName, documentID string, content interface{}) error {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	_, err := collection.Upsert(documentID, content, &gocb.UpsertOptions{Timeout: 10 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}

	log.Println("Document written successfully")
	return nil
}
