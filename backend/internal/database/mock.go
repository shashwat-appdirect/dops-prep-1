package database

import (
	"context"

	"cloud.google.com/go/firestore"
)

// MockFirestoreClient is a mock implementation of DatabaseInterface for testing
type MockFirestoreClient struct {
	ContextFunc    func() context.Context
	CollectionFunc func(name string) *firestore.CollectionRef
	CloseFunc      func() error
}

func (m *MockFirestoreClient) Context() context.Context {
	if m.ContextFunc != nil {
		return m.ContextFunc()
	}
	return context.Background()
}

func (m *MockFirestoreClient) Collection(name string) *firestore.CollectionRef {
	if m.CollectionFunc != nil {
		return m.CollectionFunc(name)
	}
	return nil
}

func (m *MockFirestoreClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

