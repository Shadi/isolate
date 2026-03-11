package sandbox

// containsSubsequence checks if haystack contains all elements of needle in order.
func containsSubsequence(haystack, needle []string) bool {
	j := 0
	for i := 0; i < len(haystack) && j < len(needle); i++ {
		if haystack[i] == needle[j] {
			j++
		}
	}
	return j == len(needle)
}

// containsElement checks if slice contains a specific string.
func containsElement(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
