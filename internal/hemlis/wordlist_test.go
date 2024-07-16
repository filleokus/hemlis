package hemlis

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	words := EncodeBytesToWords(randomBytes)
	t.Logf("Encoded words: %v", words)
	decodedBytes, _ := DecodeWordsToBytes(words)
	if !bytes.Equal(randomBytes, decodedBytes) {
		t.Errorf("Decoded bytes are not equal to original bytes")
		t.Errorf("Original bytes: %v", randomBytes)
		t.Errorf("Decoded bytes: %v", decodedBytes)
	}
}

func TestFailDecode(t *testing.T) {
	words := []string{"totally", "random", "words", "not", "in", "list"}
	decodedBytes, err := DecodeWordsToBytes(words)
	if decodedBytes != nil {
		t.Errorf("Decoded bytes should be nil")
	}
	if err == nil {
		t.Errorf("Error should not be nil")
	}
}
