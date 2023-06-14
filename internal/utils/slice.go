package utils

func FindIndexInString(entries []string, predicate func(string, int) bool) int {
	for i, entry := range entries {
		if predicate(entry, i) {
			return i
		}
	}
	return -1
}
