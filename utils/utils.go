package utils

// Check panics if the error is not nil.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
