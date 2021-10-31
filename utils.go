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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
