package docs

import (
	"context"
	"time"
)

// DocRecord represents a stored documentation artifact.
type DocRecord struct {
	Repo      string    `firestore:"repo"`
	Content   string    `firestore:"content"`
	GeneratedAt time.Time `firestore:"generated_at"`
}

// Storer defines the interface for persisting documentation.
type Storer interface {
	Save(ctx context.Context, record DocRecord) error
	Get(ctx context.Context, repo string) (*DocRecord, error)
}
