package slice

import "testing"

func TestMap(t *testing.T) {
	original := makeInterfaceSlice([]string{"str1", "str2", "str3"})
	result := Map(original, func(str interface{}) interface{} {
		return str.(string) + "1"
	})

	expected := makeInterfaceSlice([]string{"str11", "str21", "str31"})

	if !isEqual(expected, result) {
		t.Errorf("Failed to assert the result contains the same elements as expected.\nExpected: %v\nActual: %v", expected, result)
	}
}

func TestContains(t *testing.T) {
	haystack := makeInterfaceSlice([]string{"str1", "str2", "str3"})

	var simpleEqual = func(item interface{}, needle interface{}) bool { return item.(string) == needle.(string) }

	if Contains(haystack, "str4", simpleEqual) {
		t.Errorf("Fail to detect non-existing element \"str4\" in %v", haystack)
	}

	if !Contains(haystack, "str1", simpleEqual) {
		t.Errorf("Fail to include existing element \"str1\" in %v", haystack)
	}
}

func TestFirst(t *testing.T) {
	haystack := makeInterfaceSlice([]string{"str1", "str2", "str3"})

	var simpleEqual = func(item interface{}, needle interface{}) bool { return item.(string) == needle.(string) }

	firstEle, hasMatch := First(haystack, "str1", simpleEqual)

	if !hasMatch {
		t.Errorf("Fail to find an element that matches the condition.")
	}

	if firstEle.(string) != "str1" {
		t.Errorf("Fail to get the first element in %v\nGot %v as result", haystack, firstEle)
	}
}

func makeInterfaceSlice(values []string) []interface{} {
	result := make([]interface{}, len(values))
	for i, s := range values {
		result[i] = s
	}
	return result
}

func isEqual(a, b []interface{}) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
