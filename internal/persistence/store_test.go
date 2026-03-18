package persistence

import (
	"context"
	"path/filepath"
	"testing"
)

func closeStore(t *testing.T, s *Store) {
	t.Helper()
	if err := s.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

func TestOpen(t *testing.T) {
	ctx := context.Background()

	s, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:): %v", err)
	}
	defer closeStore(t, s)
}

func TestOpen_WALEnabled(t *testing.T) {
	ctx := context.Background()

	// WAL requires a file-backed database; in-memory always reports "memory".
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	s, err := Open(ctx, dbPath)
	if err != nil {
		t.Fatalf("Open(%q): %v", dbPath, err)
	}
	defer closeStore(t, s)

	var mode string
	if err := s.db.QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&mode); err != nil {
		t.Fatalf("query journal_mode: %v", err)
	}
	if mode != "wal" {
		t.Fatalf("journal_mode = %q, want %q", mode, "wal")
	}
}

// TestOpen_BasicCRUD is a smoke test for the modernc.org/sqlite driver.
// It verifies that basic SQL operations work through the configured connection.
// This test will be superseded by CRUD method tests in tasks 2.3–2.5.
func TestOpen_BasicCRUD(t *testing.T) {
	ctx := context.Background()

	s, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:): %v", err)
	}
	defer closeStore(t, s)

	// Create.
	if _, err := s.db.ExecContext(ctx, `CREATE TABLE test_items (
		id   INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	)`); err != nil {
		t.Fatalf("CREATE TABLE: %v", err)
	}

	// Insert.
	if _, err := s.db.ExecContext(ctx, `INSERT INTO test_items (id, name) VALUES (1, 'alpha')`); err != nil {
		t.Fatalf("INSERT: %v", err)
	}

	// Select.
	var id int
	var name string
	if err := s.db.QueryRowContext(ctx, `SELECT id, name FROM test_items WHERE id = 1`).Scan(&id, &name); err != nil {
		t.Fatalf("SELECT: %v", err)
	}
	if id != 1 || name != "alpha" {
		t.Fatalf("got (%d, %q), want (1, %q)", id, name, "alpha")
	}

	// Delete.
	res, err := s.db.ExecContext(ctx, `DELETE FROM test_items WHERE id = 1`)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("RowsAffected: %v", err)
	}
	if n != 1 {
		t.Fatalf("deleted %d rows, want 1", n)
	}
}
