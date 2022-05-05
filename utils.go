package main

// Basic contains function to check if element is contained in slice
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

