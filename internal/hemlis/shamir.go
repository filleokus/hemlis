package hemlis

import (
	"github.com/filleokus/hemlis/internal/shamir"
)

func SplitSecret(secret []byte, threshold, shares uint) ([][]byte, error) {

	return shamir.Split(secret, int(shares), int(threshold))
}

func CombineSecret(shares [][]byte) ([]byte, error) {
	return shamir.Combine(shares)
}
