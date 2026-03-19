package domain

import "fmt"

// TrackerErrorKind enumerates the normalized error categories that
// tracker adapters map their native errors to. The orchestrator uses
// these categories to decide retry, skip, or fail behavior.
type TrackerErrorKind string

const (
	// ErrUnsupportedTrackerKind indicates the configured tracker kind
	// has no registered adapter.
	ErrUnsupportedTrackerKind TrackerErrorKind = "unsupported_tracker_kind"

	// ErrMissingTrackerAPIKey indicates the tracker API key is absent
	// after environment variable resolution.
	ErrMissingTrackerAPIKey TrackerErrorKind = "missing_tracker_api_key"

	// ErrMissingTrackerProject indicates the tracker project is absent
	// when required by the adapter.
	ErrMissingTrackerProject TrackerErrorKind = "missing_tracker_project"

	// ErrTrackerTransport indicates a network or transport failure.
	ErrTrackerTransport TrackerErrorKind = "tracker_transport_error"

	// ErrTrackerAuth indicates an authentication or authorization failure.
	ErrTrackerAuth TrackerErrorKind = "tracker_auth_error"

	// ErrTrackerAPI indicates a non-200 HTTP or API-level error.
	ErrTrackerAPI TrackerErrorKind = "tracker_api_error"

	// ErrTrackerPayload indicates a malformed or unexpected response
	// structure from the tracker.
	ErrTrackerPayload TrackerErrorKind = "tracker_payload_error"

	// ErrTrackerMissingCursor indicates a pagination integrity error
	// where the expected end cursor is absent.
	ErrTrackerMissingCursor TrackerErrorKind = "tracker_missing_end_cursor"
)

// TrackerError is a structured error returned by [TrackerAdapter]
// implementations. The Kind field enables the orchestrator to make
// category-based decisions (retry on transport, skip on auth, etc.)
// without inspecting error messages.
type TrackerError struct {
	// Kind is the normalized error category.
	Kind TrackerErrorKind

	// Message is an operator-friendly description of the failure.
	Message string

	// Err is the underlying error, if any.
	Err error
}

// Error returns a human-readable diagnostic including the error
// category and message.
func (e *TrackerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("tracker: %s: %s: %v", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("tracker: %s: %s", e.Kind, e.Message)
}

// Unwrap returns the underlying error for use with [errors.Is] and
// [errors.As].
func (e *TrackerError) Unwrap() error {
	return e.Err
}
