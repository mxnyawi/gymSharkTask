package db

import (
	"log"

	"github.com/couchbase/gocb/v2"
)

func (db *DBManager) GetDocument(bucketName, scopeName, collectionName, documentID string) (*DocumentHistory, error) {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	var document DocumentHistory
	docOut, err := collection.Get(documentID, &gocb.GetOptions{})
	if err != nil {
		log.Printf("Failed to get document: %v", err)
		return nil, err
	}

	err = docOut.Content(&document)
	if err != nil {
		log.Printf("Failed to get document content: %v", err)
		return nil, err
	}

	log.Println("Document retrieved successfully")
	return &document, nil
}

func (db *DBManager) GetUser(bucketName, scopeName, collectionName, documentID string) (*User, error) {
	collection := db.Cluster.Bucket(bucketName).Scope(scopeName).Collection(collectionName)

	var document User
	docOut, err := collection.Get(documentID, &gocb.GetOptions{})
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, err
	}

	err = docOut.Content(&document)
	if err != nil {
		log.Printf("Failed to get user content: %v", err)
		return nil, err
	}

	log.Println("User retrieved successfully")
	return &document, nil
}
