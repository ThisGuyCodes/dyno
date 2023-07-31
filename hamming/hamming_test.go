package hamming_test

import (
	"testing"

	"github.com/thisguycodes/dyno/hamming"
)

func TestHamming(t *testing.T) {
	// examples coppied from https://en.wikipedia.org/wiki/Hamming_distance
	inputs := []string{
		"karolin", "kathrin",
		"karolin", "kerstin",
		"kathrin", "kerstin",
		"0000", "1111",
		"2173896", "2233796",
	}
	expected := []int{
		3,
		3,
		4,
		4,
		3,
	}
	for i, expect := range expected {
		left, right := inputs[i*2], inputs[i*2+1]
		if result := hamming.Distance([]byte(left), []byte(right)); expect != result {
			t.Errorf("test %d got result %d, expected %d", i, result, expect)
		}
	}
}
