package utils

func Find[T any](arr []T, b func(T) bool) T {
	for _, elem := range arr {
		if b(elem) {
			return elem
		}
	}
	return GetDefault[T]()
}

func Map[T, G any](arr []T, b func(T) G) []G {
	result := []G{}
	for _, elem := range arr {
		result = append(result, b(elem))
	}
	return result
}

func GetDefault[T any]() T {
	var result T
	return result
}
