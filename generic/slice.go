package generic

func MapDeref[T any](ptrs []*T) []T {
	var result []T
	for _, ptr := range ptrs {
		if ptr != nil {
			result = append(result, *ptr)
		}
	}
	return result
}

func Prepend[T any](slice []T, element T) []T {
	return append([]T{element}, slice...)
}
