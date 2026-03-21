package orchestrator

import "strings"

// GlobalAvailableSlots returns the number of dispatch slots available under
// the global concurrency limit. The result is never negative.
func GlobalAvailableSlots(maxConcurrentAgents, runningCount int) int {
	return max(maxConcurrentAgents-runningCount, 0)
}

// StateAvailableSlots returns the number of dispatch slots available for the
// given issue state. If a per-state limit exists in maxConcurrentByState, it
// returns max(limit - stateRunningCount, 0). If no per-state limit exists,
// it returns globalAvailable as a fallback. The state parameter is normalized
// to lowercase before lookup.
func StateAvailableSlots(state string, maxConcurrentByState map[string]int, stateRunningCount, globalAvailable int) int {
	normalized := strings.ToLower(state)
	if limit, ok := maxConcurrentByState[normalized]; ok {
		return max(limit-stateRunningCount, 0)
	}
	return globalAvailable
}

// HasAvailableSlots reports whether dispatch is permitted for an issue in
// the given state, considering both global and per-state limits.
func HasAvailableSlots(s *State, issueState string) bool {
	globalAvail := GlobalAvailableSlots(s.MaxConcurrentAgents, len(s.Running))
	if globalAvail == 0 {
		return false
	}
	stateRunning := RunningCountByState(s.Running, issueState)
	stateAvail := StateAvailableSlots(issueState, s.MaxConcurrentByState, stateRunning, globalAvail)
	return stateAvail > 0
}
