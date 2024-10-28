package utils

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func RemoveItem[T comparable](slice []T, item T) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}
