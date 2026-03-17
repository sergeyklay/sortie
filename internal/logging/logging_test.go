package logging_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/sergeyklay/sortie/internal/logging"
)

func TestSetup(t *testing.T) {
	var buf bytes.Buffer
	logging.Setup(&buf, slog.LevelInfo)

	slog.Default().Info("startup complete")
	output := buf.String()

	if !strings.Contains(output, "startup complete") {
		t.Errorf("expected log output to contain message, got: %s", output)
	}

	buf.Reset()
	slog.Default().Debug("should be filtered")
	if buf.Len() != 0 {
		t.Errorf("expected DEBUG message to be filtered at INFO level, got: %s", buf.String())
	}
}

func TestWithIssue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	issueLogger := logging.WithIssue(logger, "10042", "PROJ-123")
	issueLogger.Info("processing issue")

	output := buf.String()

	if !strings.Contains(output, "issue_id=10042") {
		t.Errorf("expected issue_id=10042 in output, got: %s", output)
	}
	if !strings.Contains(output, "issue_identifier=PROJ-123") {
		t.Errorf("expected issue_identifier=PROJ-123 in output, got: %s", output)
	}
}

func TestWithSession(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	sessionLogger := logging.WithSession(logger, "sess-abc-def")
	sessionLogger.Info("session started")

	output := buf.String()

	if !strings.Contains(output, "session_id=sess-abc-def") {
		t.Errorf("expected session_id=sess-abc-def in output, got: %s", output)
	}
}

func TestWithIssueAndSession(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	combined := logging.WithSession(logging.WithIssue(logger, "10042", "PROJ-123"), "sess-abc-def")
	combined.Info("dispatching agent")

	output := buf.String()

	for _, key := range []string{"issue_id=10042", "issue_identifier=PROJ-123", "session_id=sess-abc-def"} {
		if !strings.Contains(output, key) {
			t.Errorf("expected %s in output, got: %s", key, output)
		}
	}
}
