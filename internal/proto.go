package proto

import "math/rand"

const (
	SHA256 = iota
)

const (
	Difficulty = 10
)

func RandSeq(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

// checks if a byte slice has the specified number of leading zeros
func HasLeadingZeros(data [32]byte, numZeros int) bool {
	for i := 0; i < numZeros/8; i++ {
		if data[i] != 0 {
			return false
		}
	}
	remainingZeros := numZeros % 8
	if remainingZeros == 0 {
		return true
	}
	mask := byte(1 << (8 - remainingZeros))
	return (data[numZeros/8] & mask) == 0
}

var WordOfWisdom = map[int]string{
	0: "Doubt is the origin of wisdom. Augustine of Hippo",
	1: "The truest wisdom is a resolute determination. Napoleon Bonaparte",
	2: "Science is organized knowledge. Wisdom is organized life. Immanuel Kant",
}
