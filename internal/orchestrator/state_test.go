package orchestrator

import (
	"testing"

	"github.com/sortie-ai/sortie/internal/domain"
)

func TestNewState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		pollIntervalMS       int
		maxConcurrentAgents  int
		maxConcurrentByState map[string]int
		totals               AgentTotals
		wantMaxByStateLen    int
		checkAlias           bool
	}{
		{
			name:                 "nil state limits map becomes empty non-nil map",
			pollIntervalMS:       5000,
			maxConcurrentAgents:  10,
			maxConcurrentByState: nil,
			totals: AgentTotals{
				InputTokens:    1,
				OutputTokens:   2,
				TotalTokens:    3,
				SecondsRunning: 4.5,
			},
			wantMaxByStateLen: 0,
			checkAlias:        false,
		},
		{
			name:                "non-nil state limits map is stored as-is",
			pollIntervalMS:      1000,
			maxConcurrentAgents: 6,
			maxConcurrentByState: map[string]int{
				"to do": 2,
			},
			totals: AgentTotals{
				InputTokens:    10,
				OutputTokens:   20,
				TotalTokens:    30,
				SecondsRunning: 40.25,
			},
			wantMaxByStateLen: 1,
			checkAlias:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewState(tt.pollIntervalMS, tt.maxConcurrentAgents, tt.maxConcurrentByState, tt.totals)

			if s == nil {
				t.Fatal("NewState() = nil, want non-nil")
			}
			if s.PollIntervalMS != tt.pollIntervalMS {
				t.Errorf("PollIntervalMS = %d, want %d", s.PollIntervalMS, tt.pollIntervalMS)
			}
			if s.MaxConcurrentAgents != tt.maxConcurrentAgents {
				t.Errorf("MaxConcurrentAgents = %d, want %d", s.MaxConcurrentAgents, tt.maxConcurrentAgents)
			}
			if s.AgentTotals != tt.totals {
				t.Errorf("AgentTotals = %+v, want %+v", s.AgentTotals, tt.totals)
			}
			if s.AgentRateLimits != nil {
				t.Errorf("AgentRateLimits = %v, want nil", s.AgentRateLimits)
			}

			if s.MaxConcurrentByState == nil {
				t.Fatal("MaxConcurrentByState = nil, want non-nil")
			}
			if len(s.MaxConcurrentByState) != tt.wantMaxByStateLen {
				t.Errorf("len(MaxConcurrentByState) = %d, want %d", len(s.MaxConcurrentByState), tt.wantMaxByStateLen)
			}

			if s.Running == nil {
				t.Fatal("Running = nil, want non-nil")
			}
			if s.Claimed == nil {
				t.Fatal("Claimed = nil, want non-nil")
			}
			if s.RetryAttempts == nil {
				t.Fatal("RetryAttempts = nil, want non-nil")
			}
			if s.Completed == nil {
				t.Fatal("Completed = nil, want non-nil")
			}

			if tt.checkAlias {
				tt.maxConcurrentByState["in progress"] = 3
				if got := s.MaxConcurrentByState["in progress"]; got != 3 {
					t.Errorf("MaxConcurrentByState aliasing check = %d, want 3", got)
				}
			}
		})
	}
}

func TestRunningCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		running map[string]*RunningEntry
		want    int
	}{
		{
			name:    "empty running map",
			running: map[string]*RunningEntry{},
			want:    0,
		},
		{
			name: "three running entries",
			running: map[string]*RunningEntry{
				"1": {Issue: domain.Issue{State: "To Do"}},
				"2": {Issue: domain.Issue{State: "In Progress"}},
				"3": {Issue: domain.Issue{State: "Done"}},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &State{Running: tt.running}
			got := s.RunningCount()
			if got != tt.want {
				t.Errorf("RunningCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRunningCountByState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		running map[string]*RunningEntry
		state   string
		want    int
	}{
		{
			name:    "empty running map",
			running: map[string]*RunningEntry{},
			state:   "in progress",
			want:    0,
		},
		{
			name: "case-insensitive match with mixed states",
			running: map[string]*RunningEntry{
				"1": {Issue: domain.Issue{State: "To Do"}},
				"2": {Issue: domain.Issue{State: "In Progress"}},
				"3": {Issue: domain.Issue{State: "in progress"}},
			},
			state: "IN PROGRESS",
			want:  2,
		},
		{
			name: "absent state",
			running: map[string]*RunningEntry{
				"1": {Issue: domain.Issue{State: "To Do"}},
				"2": {Issue: domain.Issue{State: "In Progress"}},
			},
			state: "blocked",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := RunningCountByState(tt.running, tt.state)
			if got != tt.want {
				t.Errorf("RunningCountByState(..., %q) = %d, want %d", tt.state, got, tt.want)
			}
		})
	}
}
