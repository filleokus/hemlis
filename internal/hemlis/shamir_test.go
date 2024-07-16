package hemlis

import (
	"bytes"
	"testing"
)

func TestSplitAndCombineSecret(t *testing.T) {
	secret_bytes := []byte("This is a secret message")
	num_threshold := uint(2)
	num_shares := uint(3)
	shares, err := SplitSecret(secret_bytes, num_threshold, num_shares)
	if err != nil {
		t.Errorf("Error splitting secret: %v", err)
	}

	combinedSecret, err := CombineSecret(shares[:num_threshold])
	if err != nil {
		t.Errorf("Error combining secret: %v", err)
	}

	if bytes.Equal(combinedSecret, secret_bytes) == false {
		t.Errorf("Combined secret does not match original secret")
	}
}
