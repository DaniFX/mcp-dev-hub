package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// Client wraps the Firestore client with helper methods.
type Client struct {
	fs *firestore.Client
}

// New creates a new Firestore client for the given GCP project.
func New(ctx context.Context, projectID string) (*Client, error) {
	fs, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	log.Printf("Firestore client initialized for project: %s", projectID)
	return &Client{fs: fs}, nil
}

// Close closes the underlying Firestore connection.
func (c *Client) Close() error {
	return c.fs.Close()
}
