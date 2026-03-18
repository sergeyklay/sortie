// Package persistence provides SQLite-backed durable storage for retry queues,
// run history, session metadata, and aggregate metrics.
package persistence

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Register the "sqlite" database/sql driver.
)

// Store provides single-writer access to the SQLite database. All write
// operations must be serialized through a single Store instance. Concurrent
// reads are safe in WAL mode.
type Store struct {
	db *sql.DB
}

// Open creates or opens a SQLite database at the given path and enables WAL
// journal mode. Use ":memory:" for an in-memory database (useful for testing).
// The caller must call [Store.Close] when the store is no longer needed.
func Open(ctx context.Context, path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %q: %w", path, err)
	}

	// Verify the connection is usable.
	if err := db.PingContext(ctx); err != nil {
		db.Close() //nolint:errcheck // best-effort cleanup on open failure
		return nil, fmt.Errorf("ping sqlite %q: %w", path, err)
	}

	// Enable WAL mode for concurrent reads with single-writer semantics.
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		db.Close() //nolint:errcheck // best-effort cleanup on open failure
		return nil, fmt.Errorf("enable WAL mode: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
