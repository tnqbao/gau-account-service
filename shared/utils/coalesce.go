package utils

// Coalesce returns the first non-nil value in the list
func Coalesce[T any](vals ...*T) *T {
	for _, v := range vals {
		if v != nil {
			return v
		}
	}
	return nil
}
