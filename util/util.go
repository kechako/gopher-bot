// Package util provides utility functions for plugins.
package util

import "strings"

// MentionToBot returns true when mentions include the bot userID.
func MentionToBot(userID string, mentions []string) bool {
	for _, m := range mentions {
		if m == userID {
			return true
		}
	}

	return false
}

// HasKeywords find keywords in the text.
func HasKeywords(text string, partial bool, keywords ...string) bool {
	if partial {
		for _, k := range keywords {
			if strings.Contains(text, k) {
				return true
			}
		}
	} else {
		for _, k := range keywords {
			if text == k {
				return true
			}
		}
	}

	return false
}
