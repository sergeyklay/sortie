package domain

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// Compile-time interface satisfaction check.
var _ TrackerAdapter = (*mockTrackerAdapter)(nil)

type mockTrackerAdapter struct{}

func (m *mockTrackerAdapter) FetchCandidateIssues(_ context.Context) ([]Issue, error) {
	return nil, nil
}

func (m *mockTrackerAdapter) FetchIssueByID(_ context.Context, _ string) (Issue, error) {
	return Issue{}, nil
}

func (m *mockTrackerAdapter) FetchIssuesByStates(_ context.Context, _ []string) ([]Issue, error) {
	return nil, nil
}

func (m *mockTrackerAdapter) FetchIssueStatesByIDs(_ context.Context, _ []string) (map[string]string, error) {
	return nil, nil
}

func (m *mockTrackerAdapter) FetchIssueComments(_ context.Context, _ string) ([]Comment, error) {
	return nil, nil
}

func TestTrackerError_Error(t *testing.T) {
	err := &TrackerError{
		Kind:    ErrTrackerTransport,
		Message: "connection refused",
	}
	got := err.Error()
	if got != "tracker: tracker_transport_error: connection refused" {
		t.Errorf("Error() = %q, want %q", got, "tracker: tracker_transport_error: connection refused")
	}
}

func TestTrackerError_ErrorWithWrapped(t *testing.T) {
	inner := fmt.Errorf("dial tcp: connect refused")
	err := &TrackerError{
		Kind:    ErrTrackerTransport,
		Message: "connection failed",
		Err:     inner,
	}
	got := err.Error()
	want := "tracker: tracker_transport_error: connection failed: dial tcp: connect refused"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestTrackerError_Unwrap(t *testing.T) {
	inner := fmt.Errorf("underlying error")
	trackerErr := &TrackerError{
		Kind:    ErrTrackerAuth,
		Message: "invalid token",
		Err:     inner,
	}

	if trackerErr.Unwrap() != inner {
		t.Error("Unwrap() did not return the inner error")
	}

	// Verify errors.As works through a wrapping chain.
	wrapped := fmt.Errorf("outer: %w", trackerErr)
	var extracted *TrackerError
	if !errors.As(wrapped, &extracted) {
		t.Fatal("errors.As failed to extract *TrackerError from wrapped chain")
	}
	if extracted.Kind != ErrTrackerAuth {
		t.Errorf("extracted.Kind = %q, want %q", extracted.Kind, ErrTrackerAuth)
	}
}

func TestTrackerError_UnwrapNil(t *testing.T) {
	err := &TrackerError{
		Kind:    ErrTrackerPayload,
		Message: "unexpected field",
	}
	if err.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", err.Unwrap())
	}
}
