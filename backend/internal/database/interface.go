package database

import (
	"context"

	"cloud.google.com/go/firestore"
)

// DatabaseInterface defines the interface for database operations
// This allows for easy mocking in tests
type DatabaseInterface interface {
	Context() context.Context
	Collection(name string) *firestore.CollectionRef
	Close() error
}

