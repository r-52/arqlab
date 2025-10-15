package test262

import "strings"

var asyncKeywords = []string{"async", "await", "async-functions"}

// IsAsyncRelated returns true when a test case should be excluded because it
// targets async/await semantics that are intentionally unsupported.
func IsAsyncRelated(tc TestCase) bool {
	lowerPath := strings.ToLower(tc.Path)
	for _, keyword := range asyncKeywords {
		if strings.Contains(lowerPath, keyword) {
			return true
		}
	}

	for _, flag := range tc.Flags {
		if strings.Contains(strings.ToLower(flag), "async") {
			return true
		}
	}

	return false
}

// FilterAsync removes async-related Test262 cases from the provided slice and
// returns the filtered list.
func FilterAsync(cases []TestCase) []TestCase {
	if len(cases) == 0 {
		return cases
	}

	filtered := make([]TestCase, 0, len(cases))
	for _, tc := range cases {
		if IsAsyncRelated(tc) {
			continue
		}
		filtered = append(filtered, tc)
	}
	return filtered
}
