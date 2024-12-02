package utils

func IsIn[T comparable](value T, okValues []T) bool {
	for _, okValue := range okValues {
		if value == okValue {
			return true
		}
	}
	return false
}
func IsOneOf[T comparable](value T, okValues ...T) bool {
	return IsIn(value, okValues)
}
