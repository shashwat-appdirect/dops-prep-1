package database

import (
	"context"
	"sync"

	"appdirect-workshop-backend/internal/config"

	"cloud.google.com/go/firestore"
)

// MockDB is an in-memory mock implementation of DBInterface for testing
type MockDB struct {
	mu          sync.RWMutex
	collections map[string]map[string]map[string]interface{}
	ctx         context.Context
	cfg         *config.Config
	nextID      int
}

// NewMockDB creates a new mock database
func NewMockDB(cfg *config.Config) *MockDB {
	return &MockDB{
		collections: make(map[string]map[string]map[string]interface{}),
		ctx:         context.Background(),
		cfg:         cfg,
		nextID:      1,
	}
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) Context() context.Context {
	return m.ctx
}

func (m *MockDB) Collection(name string) *firestore.CollectionRef {
	// Return a real CollectionRef but we'll need to intercept operations
	// For now, return nil and handle in tests differently
	// Actually, we need to create a wrapper that intercepts calls
	return nil
}

// Helper methods for test setup
func (m *MockDB) AddDocument(collectionName, docID string, data map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.collections[collectionName] == nil {
		m.collections[collectionName] = make(map[string]map[string]interface{})
	}
	m.collections[collectionName][docID] = data
}

func (m *MockDB) GetDocument(collectionName, docID string) (map[string]interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	col, exists := m.collections[collectionName]
	if !exists {
		return nil, false
	}
	doc, exists := col[docID]
	return doc, exists
}

func (m *MockDB) GetAllDocuments(collectionName string) map[string]map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	col, exists := m.collections[collectionName]
	if !exists {
		return make(map[string]map[string]interface{})
	}
	result := make(map[string]map[string]interface{})
	for k, v := range col {
		result[k] = v
	}
	return result
}

func (m *MockDB) DeleteDocument(collectionName, docID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if col, exists := m.collections[collectionName]; exists {
		delete(col, docID)
	}
}

// Ensure MockDB implements DBInterface
var _ DBInterface = (*MockDB)(nil)

