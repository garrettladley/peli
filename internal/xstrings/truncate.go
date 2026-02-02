package xstrings

import "github.com/charmbracelet/x/ansi"

// Truncate truncates a string to maxLen characters.
// If the string is longer than maxLen, it returns the first maxLen-1 characters
// followed by an ellipsis.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

// TruncateWithEllipsis truncates a string to maxWidth display width.
// It uses ANSI-aware width calculation to handle unicode correctly.
func TruncateWithEllipsis(s string, maxWidth int) string {
	if ansi.StringWidth(s) <= maxWidth {
		return s
	}
	return ansi.Truncate(s, maxWidth, "…")
}
