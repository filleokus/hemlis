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
	decodedBytes := DecodeWordsToBytes(words)
	if !bytes.Equal(randomBytes, decodedBytes) {
		t.Errorf("Decoded bytes are not equal to original bytes")
		t.Errorf("Original bytes: %v", randomBytes)
		t.Errorf("Decoded bytes: %v", decodedBytes)
	}
}
