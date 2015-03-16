package utils

func Extend(slice []byte, sliceTwo []byte) []byte {
	for i := range sliceTwo {
		slice = append(slice, sliceTwo[i])
	}

	return slice
}
