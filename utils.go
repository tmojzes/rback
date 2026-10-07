package main

// iff returns valTrue if cond is true, otherwise valFalse.
func iff[T any](cond bool, valTrue, valFalse T) T {
	if cond {
		return valTrue
	}
	return valFalse
}

// getOrDefault safely returns a value of type T from map m for the given key,
// or fallback if the key is missing or not of type T.
func getOrDefault[T any](m map[string]any, key string, fallback T) T {
	if m == nil {
		return fallback
	}
	if v, ok := m[key]; ok && v != nil {
		if typed, ok := v.(T); ok {
			return typed
		}
	}
	return fallback
}

// filter returns a new slice containing only the elements of items that satisfy predicate.
func filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// mapSlice transforms each element of items using the transform function.
func mapSlice[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = transform(item)
	}
	return result
}
