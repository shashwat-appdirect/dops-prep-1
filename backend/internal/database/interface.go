package database

import (
	"context"

	"cloud.google.com/go/firestore"
)

// DBInterface defines the interface for database operations
// This allows us to mock Firestore in tests
type DBInterface interface {
	Close() error
	Context() context.Context
	Collection(name string) *firestore.CollectionRef
}

