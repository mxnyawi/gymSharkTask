package db

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb/v2"
)

// GetDocument gets the document from the database
func (db *DBManager) GetDocument(bucketName, scopeName, collectionName, documentID string) (*DocumentHistory, error) {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	var document DocumentHistory
	docOut, err := collection.Get(documentID, &gocb.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	err = docOut.Content(&document)
	if err != nil {
		return nil, fmt.Errorf("failed to get document content: %w", err)
	}

	log.Println("Document retrieved successfully")
	return &document, nil
}

// GetUser gets the user from the database
func (db *DBManager) GetUser(bucketName, scopeName, collectionName, documentID string) (*User, error) {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	var document User
	docOut, err := collection.Get(documentID, &gocb.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = docOut.Content(&document)
	if err != nil {
		return nil, fmt.Errorf("failed to get user content: %w", err)
	}

	log.Println("User retrieved successfully")
	return &document, nil
}
