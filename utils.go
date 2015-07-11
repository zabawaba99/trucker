package main

// stringOr returns the first non empty string given
func stringOr(s ...string) string {
	for _, canidate := range s {
		if canidate != "" {
			return canidate
		}
	}

	return ""
}
