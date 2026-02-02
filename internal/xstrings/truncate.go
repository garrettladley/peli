package xstrings

import "github.com/charmbracelet/x/ansi"

// Truncate truncates a string to maxWidth display width.
// It uses ANSI-aware width calculation to handle unicode correctly.
func Truncate(s string, maxWidth int) string {
	if ansi.StringWidth(s) <= maxWidth {
		return s
	}
	return ansi.Truncate(s, maxWidth, "…")
}
