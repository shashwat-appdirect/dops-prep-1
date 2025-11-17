package database

import (
	"context"
	"encoding/json"
	"fmt"

	"appdirect-workshop-backend/internal/config"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FirestoreClient struct {
	client *firestore.Client
	ctx    context.Context
	cfg    *config.Config
}

// Ensure FirestoreClient implements DBInterface
var _ DBInterface = (*FirestoreClient)(nil)

func NewFirestoreClient(cfg *config.Config) (*FirestoreClient, error) {
	ctx := context.Background()

	// Convert service account map to JSON
	serviceAccountJSON, err := json.Marshal(cfg.FirebaseServiceAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service account: %v", err)
	}

	// Initialize Firebase app
	opt := option.WithCredentialsJSON(serviceAccountJSON)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firestore client: %v", err)
	}

	return &FirestoreClient{
		client: client,
		ctx:    ctx,
		cfg:    cfg,
	}, nil
}

func (f *FirestoreClient) Close() error {
	return f.client.Close()
}

func (f *FirestoreClient) Context() context.Context {
	return f.ctx
}

func (f *FirestoreClient) Collection(name string) *firestore.CollectionRef {
	// Use subcollection ID as a document reference, then access collections as subcollections
	docRef := f.client.Collection("workshops").Doc(f.cfg.SubcollectionID)
	return docRef.Collection(name)
}

