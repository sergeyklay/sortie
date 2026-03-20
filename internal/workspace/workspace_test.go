package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
		wantOp  string
	}{
		{"simple alphanumeric with hyphen", "ABC-123", "ABC-123", false, ""},
		{"dots and underscores preserved", "my.task_1", "my.task_1", false, ""},
		{"slashes replaced", "PROJ/sub-task", "PROJ_sub-task", false, ""},
		{"spaces replaced", "My Task 42", "My_Task_42", false, ""},
		{"unicode replaced", "日本語-タスク", "___-___", false, ""},
		{"special chars replaced", "a@b#c$d%e", "a_b_c_d_e", false, ""},
		{"all replaced chars", "///", "___", false, ""},
		{"single valid char", "a", "a", false, ""},
		{"consecutive replacements not collapsed", "A//B", "A__B", false, ""},
		{"backslash replaced", `A\B`, "A_B", false, ""},
		{"null byte replaced", "A\x00B", "A_B", false, ""},
		{"empty input", "", "", true, "sanitize"},
		{"result is dot", ".", "", true, "sanitize"},
		{"result is dotdot", "..", "", true, "sanitize"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := SanitizeKey(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("SanitizeKey(%q) = %q, want error", tt.input, got)
				}
				var pe *PathError
				if !errors.As(err, &pe) {
					t.Fatalf("SanitizeKey(%q) error type = %T, want *PathError", tt.input, err)
				}
				if pe.Op != tt.wantOp {
					t.Errorf("PathError.Op = %q, want %q", pe.Op, tt.wantOp)
				}
				return
			}

			if err != nil {
				t.Fatalf("SanitizeKey(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("SanitizeKey(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestComputePath(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		res, err := ComputePath(root, "ABC-123")
		if err != nil {
			t.Fatalf("ComputePath(%q, %q) error: %v", root, "ABC-123", err)
		}
		if res.Key != "ABC-123" {
			t.Errorf("Key = %q, want %q", res.Key, "ABC-123")
		}
		wantPath := filepath.Join(root, "ABC-123")
		if res.Path != wantPath {
			t.Errorf("Path = %q, want %q", res.Path, wantPath)
		}
		if !filepath.IsAbs(res.Path) {
			t.Errorf("Path %q is not absolute", res.Path)
		}
	})

	t.Run("root with trailing slash", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		res, err := ComputePath(root+"/", "X-1")
		if err != nil {
			t.Fatalf("ComputePath(%q, %q) error: %v", root+"/", "X-1", err)
		}
		wantPath := filepath.Join(root, "X-1")
		if res.Path != wantPath {
			t.Errorf("Path = %q, want %q", res.Path, wantPath)
		}
	})

	t.Run("identifier needs sanitization", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		res, err := ComputePath(root, "A/B#C")
		if err != nil {
			t.Fatalf("ComputePath(%q, %q) error: %v", root, "A/B#C", err)
		}
		if res.Key != "A_B_C" {
			t.Errorf("Key = %q, want %q", res.Key, "A_B_C")
		}
		wantPath := filepath.Join(root, "A_B_C")
		if res.Path != wantPath {
			t.Errorf("Path = %q, want %q", res.Path, wantPath)
		}
	})

	t.Run("empty root", func(t *testing.T) {
		t.Parallel()

		_, err := ComputePath("", "ABC-123")
		if err == nil {
			t.Fatal("ComputePath with empty root should error")
		}
		var pe *PathError
		if !errors.As(err, &pe) {
			t.Fatalf("error type = %T, want *PathError", err)
		}
		if pe.Op != "resolve" {
			t.Errorf("PathError.Op = %q, want %q", pe.Op, "resolve")
		}
	})

	t.Run("empty identifier", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		_, err := ComputePath(root, "")
		if err == nil {
			t.Fatal("ComputePath with empty identifier should error")
		}
		var pe *PathError
		if !errors.As(err, &pe) {
			t.Fatalf("error type = %T, want *PathError", err)
		}
		if pe.Op != "sanitize" {
			t.Errorf("PathError.Op = %q, want %q", pe.Op, "sanitize")
		}
	})

	t.Run("root does not exist yet", func(t *testing.T) {
		t.Parallel()
		base := t.TempDir()
		nonexistent := filepath.Join(base, "nonexistent", "sub")

		res, err := ComputePath(nonexistent, "X-1")
		if err != nil {
			t.Fatalf("ComputePath(%q, %q) error: %v", nonexistent, "X-1", err)
		}
		if res.Key != "X-1" {
			t.Errorf("Key = %q, want %q", res.Key, "X-1")
		}
		wantPath := filepath.Join(nonexistent, "X-1")
		if res.Path != wantPath {
			t.Errorf("Path = %q, want %q", res.Path, wantPath)
		}
	})

	// Section 9.5, Invariant 2: root="/" must not cause false rejection
	t.Run("root is filesystem root", func(t *testing.T) {
		t.Parallel()

		res, err := ComputePath("/", "X-1")
		if err != nil {
			t.Fatalf("ComputePath(%q, %q) error: %v", "/", "X-1", err)
		}
		if res.Path != "/X-1" {
			t.Errorf("Path = %q, want %q", res.Path, "/X-1")
		}
		if res.Key != "X-1" {
			t.Errorf("Key = %q, want %q", res.Key, "X-1")
		}
	})

	t.Run("path is always under root", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		identifiers := []string{"ABC-123", "A/B#C", "日本語", "my.task_1"}
		for _, id := range identifiers {
			res, err := ComputePath(root, id)
			if err != nil {
				t.Errorf("ComputePath(%q, %q) error: %v", root, id, err)
				continue
			}

			dir := filepath.Dir(res.Path)
			if dir != root {
				t.Errorf("ComputePath(%q, %q): Dir(Path) = %q, want %q", root, id, dir, root)
			}
			if filepath.Base(res.Path) != res.Key {
				t.Errorf("ComputePath(%q, %q): Base(Path) = %q, want Key %q", root, id, filepath.Base(res.Path), res.Key)
			}
		}
	})

	t.Run("dot identifier rejected", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		_, err := ComputePath(root, ".")
		if err == nil {
			t.Fatal("ComputePath with dot identifier should error")
		}
		var pe *PathError
		if !errors.As(err, &pe) {
			t.Fatalf("error type = %T, want *PathError", err)
		}
		if pe.Op != "sanitize" {
			t.Errorf("PathError.Op = %q, want %q", pe.Op, "sanitize")
		}
	})

	t.Run("dotdot identifier rejected", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()

		_, err := ComputePath(root, "..")
		if err == nil {
			t.Fatal("ComputePath with dotdot identifier should error")
		}
		var pe *PathError
		if !errors.As(err, &pe) {
			t.Fatalf("error type = %T, want *PathError", err)
		}
		if pe.Op != "sanitize" {
			t.Errorf("PathError.Op = %q, want %q", pe.Op, "sanitize")
		}
	})
}

// Section 9.5, Invariant 2: symlinked root resolves to real directory
func TestComputePath_SymlinkRoot(t *testing.T) {
	t.Parallel()

	realRoot := t.TempDir()
	symlinkDir := t.TempDir()
	symlinkPath := filepath.Join(symlinkDir, "symlink-root")

	if err := os.Symlink(realRoot, symlinkPath); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	res, err := ComputePath(symlinkPath, "X-1")
	if err != nil {
		t.Fatalf("ComputePath(%q, %q) error: %v", symlinkPath, "X-1", err)
	}

	// Path must be under the real root, not the symlink path
	if strings.HasPrefix(res.Path, symlinkPath) {
		t.Errorf("Path %q is under symlink path %q, should be under real root %q", res.Path, symlinkPath, realRoot)
	}

	wantPath := filepath.Join(realRoot, "X-1")
	if res.Path != wantPath {
		t.Errorf("Path = %q, want %q", res.Path, wantPath)
	}
}

func TestPathError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  PathError
		want string
	}{
		{
			"op with root and identifier",
			PathError{Op: "containment", Root: "/tmp/root", Identifier: "ABC-123", Err: errors.New("escaped")},
			`workspace containment: root="/tmp/root" identifier="ABC-123": escaped`,
		},
		{
			"op with identifier only",
			PathError{Op: "sanitize", Identifier: "bad-id", Err: errors.New("empty")},
			`workspace sanitize: identifier="bad-id": empty`,
		},
		{
			"op with root only",
			PathError{Op: "resolve", Root: "/tmp/root", Err: errors.New("bad")},
			`workspace resolve: root="/tmp/root": bad`,
		},
		{
			"op only",
			PathError{Op: "resolve", Err: errors.New("missing")},
			`workspace resolve: missing`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPathError_Unwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner")
	pe := &PathError{Op: "test", Err: inner}

	if !errors.Is(pe, inner) {
		t.Error("errors.Is should find the wrapped error")
	}
}
