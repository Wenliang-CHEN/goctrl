package slice

func ContainsString(haystack []interface{}, needle string) bool {
	for _, item := range haystack {
		if item.(string) == needle {
			return true
		}
	}
	return false
}
