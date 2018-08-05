package slice

func ContainsName(kubeConfigObjects []interface{}, name string) bool {
	for _, item := range kubeConfigObjects {
		if toMap(item)["name"].(string) == name {
			return true
		}
	}
	return false
}

func toMap(item interface{}) map[interface{}]interface{} {
	value, ok := item.(map[interface{}]interface{})
	if ok == false {
		panic("Unable to convert parameters.  Please check the settings in config file.")
	}

	return value
}
