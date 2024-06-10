package db

import (
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
		log.Printf("Failed to create admin user: %v", err)
		return err
	}

	log.Println("Admin user created successfully")
	return nil
}

// SetupDB sets up the database
func (db *DBManager) SetupDB(bucketName, scopeName, collectionName, documentID string) error {
	err := db.CreateBucket(bucketName)
	if err != nil {
		return err
	}

	err = db.CreateScope(bucketName, scopeName)
	if err != nil {
		return err
	}

	err = db.CreateCollection(bucketName, scopeName, collectionName)
	if err != nil {
		return err
	}

	err = db.WriteDocument(bucketName, scopeName, collectionName, documentID, DocumentHistory{History: []Document{}})
	if err != nil {
		log.Printf("Failed to write order document: %v", err)
		return err
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
		log.Printf("Failed to create bucket: %v", err)
		return err
	}
}

// CreateScope creates a scope in a bucket
func (db *DBManager) CreateScope(bucketName, scopeName string) error {
	err := db.Cluster.Bucket(bucketName).Collections().CreateScope(scopeName, nil)
	if err != nil {
		log.Printf("Failed to create scope: %v", err)
		return err
	}

	log.Println("Scope created successfully")
	return nil
}

// CreateCollection creates a collection in a scope
func (db *DBManager) CreateCollection(bucketName, scopeName, collectionName string) error {

	collectionNameSettings := gocb.CollectionSpec{
		Name:      collectionName,
		ScopeName: scopeName,
	}

	err := db.Cluster.Bucket(bucketName).Collections().CreateCollection(collectionNameSettings, nil)
	if err != nil {
		log.Printf("Failed to create collection: %v", err)
		return err
	}

	log.Println("Collection created successfully")
	return nil
}

// WriteDocument writes a document to the database collection
func (db *DBManager) WriteDocument(bucketName, scopeName, collectionName, documentID string, content interface{}) error {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	_, err := collection.Upsert(documentID, content, &gocb.UpsertOptions{Timeout: 10 * time.Second})
	if err != nil {
		log.Printf("Failed to write document: %v", err)
		return err
	}

	log.Println("Document written successfully")
	return nil
}

// UpdateDocument updates a document in the database collection with new content
func (db *DBManager) UpdateDocument(bucketName, scopeName, collectionName, documentID string, content map[string]interface{}) error {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	mops := make([]gocb.MutateInSpec, 0, len(content))
	for k, v := range content {
		mops = append(mops, gocb.UpsertSpec(k, v, &gocb.UpsertSpecOptions{}))
	}

	_, err := collection.MutateIn(documentID, mops, &gocb.MutateInOptions{Timeout: 10 * time.Second})
	if err != nil {
		log.Printf("Failed to update document: %v", err)
		return err
	}

	log.Println("Document updated successfully")
	return nil
}
