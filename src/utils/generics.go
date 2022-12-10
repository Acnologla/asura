package utils

func Find[T any](arr []T, b func(T) bool) T {
	for _, elem := range arr {
		if b(elem) {
			return elem
		}
	}
	return GetDefault[T]()
}

func GetDefault[T any]() T {
	var result T
	return result
}
