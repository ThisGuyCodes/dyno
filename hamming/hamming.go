package hamming

// Distance returns the hamming distance of the two input byte slices.
// Differences in length are added to the distance.
func Distance(left []byte, right []byte) int {
	// ensure left is always the shorter one
	if len(left) > len(right) {
		return Distance(right, left)
	}

	accumulator := len(right) - len(left)
	for i, leftByte := range left {
		if leftByte != right[i] {
			accumulator += 1
		}
	}

	return accumulator
}
