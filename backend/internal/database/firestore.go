package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

func NewFirestoreClient(cfg *config.Config) (*FirestoreClient, error) {
	ctx := context.Background()

	var app *firebase.App
	var err error

	// Check if we're running on Cloud Run (use Application Default Credentials)
	// Cloud Run automatically provides credentials via ADC
	if os.Getenv("K_SERVICE") != "" || os.Getenv("GOOGLE_CLOUD_PROJECT") != "" {
		// Use Application Default Credentials (ADC) for Cloud Run
		app, err = firebase.NewApp(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase app with ADC: %v", err)
		}
	} else {
		// Use service account for local development
		if cfg.FirebaseServiceAccount == nil || len(cfg.FirebaseServiceAccount) == 0 {
			return nil, fmt.Errorf("FIREBASE_SERVICE_ACCOUNT is required for local development")
		}

		// Convert service account map to JSON
		serviceAccountJSON, err := json.Marshal(cfg.FirebaseServiceAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal service account: %v", err)
		}

		// Initialize Firebase app with service account
		opt := option.WithCredentialsJSON(serviceAccountJSON)
		app, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase app: %v", err)
		}
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

