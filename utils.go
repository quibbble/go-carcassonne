package go_carcassonne

func contains(items []string, item string) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}

func indexOf(items []string, item string) int {
	for index, it := range items {
		if it == item {
			return index
		}
	}
	return -1
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}
