package utils

// StringSliceIndex returns index of string element in string slice
func StringSliceIndex(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
