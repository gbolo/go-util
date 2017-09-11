package common

import "strings"

/* returns true if substr is in string s */
func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func CaseInsensitiveEquals(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}