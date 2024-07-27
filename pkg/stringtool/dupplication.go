package stringtool

import "strings"

func RemoveDuplicate(s string, r rune) string {
	// Remove all dupplicate rune if string
	for modified := strings.ReplaceAll(s, string(r)+string(r), string(r)); modified != s; modified = strings.ReplaceAll(s, string(r)+string(r), string(r)) {
		s = modified
	}

	return s
}
