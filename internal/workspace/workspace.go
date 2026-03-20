// Package workspace manages per-issue workspace directories, path
// safety, and lifecycle hooks. Start with [ComputePath] for safe
// workspace path derivation from issue identifiers.
package workspace

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

// unsafeChars matches any character not in the allowed set for
// workspace directory names.
var unsafeChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// PathResult holds the computed workspace path and its sanitized key.
type PathResult struct {
	// Key is the sanitized directory name derived from the issue
	// identifier. Contains only [A-Za-z0-9._-] characters.
	Key string

	// Path is the absolute workspace path: <resolved_root>/<key>.
	Path string
}

// SanitizeKey derives a safe directory name from an issue identifier
// by replacing every character not in [A-Za-z0-9._-] with underscore.
// Returns a [*PathError] if the input is empty or the sanitized result
// is "." or ".." (filesystem special names).
func SanitizeKey(identifier string) (string, error) {
	if identifier == "" {
		return "", &PathError{
			Op:         "sanitize",
			Identifier: identifier,
			Err:        errors.New("identifier must not be empty"),
		}
	}

	key := unsafeChars.ReplaceAllString(identifier, "_")

	if key == "." || key == ".." {
		return "", &PathError{
			Op:         "sanitize",
			Identifier: identifier,
			Err:        errors.New("sanitized key is a filesystem special name"),
		}
	}

	return key, nil
}

// ComputePath computes and validates the absolute workspace path for
// the given issue identifier under the specified workspace root. It
// sanitizes the identifier into a workspace key, resolves the root to
// an absolute path, joins root and key, and validates that the result
// is contained within the root directory. Returns a [PathResult] with
// the sanitized key and validated absolute path.
//
// Returns a [*PathError] on sanitization failure, root resolution
// failure, or containment violation.
func ComputePath(root, identifier string) (PathResult, error) {
	if root == "" {
		return PathResult{}, &PathError{
			Op:  "resolve",
			Err: errors.New("workspace root must not be empty"),
		}
	}

	key, err := SanitizeKey(identifier)
	if err != nil {
		return PathResult{}, err
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return PathResult{}, &PathError{
			Op:   "resolve",
			Root: root,
			Err:  err,
		}
	}

	// Resolve symlinks on the root to get the real filesystem path.
	// This prevents a symlink at root from pointing outside the
	// intended directory tree.
	resolvedRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		// Root does not exist yet — fall back to the cleaned absolute
		// path. The root directory may be created later (task 5.2).
		resolvedRoot = filepath.Clean(absRoot)
	}

	workspacePath := filepath.Join(resolvedRoot, key)

	// Containment check: workspace_path must be a direct child of
	// resolved_root. filepath.Rel handles edge cases like root="/"
	// where a naive string prefix check would fail.
	rel, err := filepath.Rel(resolvedRoot, workspacePath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == "." || strings.Contains(rel, string(filepath.Separator)) {
		return PathResult{}, &PathError{
			Op:         "containment",
			Root:       root,
			Identifier: identifier,
			Err:        errors.New("workspace path is not under root"),
		}
	}

	return PathResult{Key: key, Path: workspacePath}, nil
}
