package slice

func Map(haystack []interface{}, f func(interface{}) interface{}) []interface{} {
	target := make([]interface{}, len(haystack))
	for i, v := range haystack {
		target[i] = f(v)
	}
	return target
}

func Contains(haystack []interface{}, needle interface{}, f func(interface{}, interface{}) bool) bool {
	for _, item := range haystack {
		if f(item, needle) == true {
			return true
		}
	}
	return false
}

func First(haystack []interface{}, needle interface{}, f func(interface{}, interface{}) bool) (interface{}, bool) {
	for _, item := range haystack {
		if f(item, needle) == true {
			return item, true
		}
	}
	return nil, false
}

func ToMap(item interface{}) map[interface{}]interface{} {
	value, ok := item.(map[interface{}]interface{})
	if ok == false {
		panic("Unable to convert parameters.  Please check the settings in config file.")
	}

	return value
}
