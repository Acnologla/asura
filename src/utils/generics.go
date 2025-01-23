package utils

func Find[T any](arr []T, b func(T) bool) T {
	for _, elem := range arr {
		if b(elem) {
			return elem
		}
	}
	return GetDefault[T]()
}

func Has[T comparable](arr []T, b T) bool {
	for _, elem := range arr {
		if elem == b {
			return true
		}
	}
	return false
}

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
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

func Filter[T any](arr []T, b func(T) bool) []T {
	result := []T{}

	for _, elem := range arr {
		if b(elem) {
			result = append(result, elem)
		}
	}

	return result
}

func GetIndex[T any](arr []T, b func(T) bool) int {
	for i, elem := range arr {
		if b(elem) {
			return i
		}
	}

	return -1
}

func Chunk[T any](slice []T, size int) (chunks [][]T) {
	if size < 1 {
		panic("chunk size cannot be less than 1")
	}
	for i := 0; ; i++ {
		next := i * size
		if len(slice[next:]) > size {
			end := next + size
			chunks = append(chunks, slice[next:end:end])
		} else {
			chunks = append(chunks, slice[i*size:])
			return
		}
	}
}
