package hemlis

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"filippo.io/age"
	"github.com/filleokus/hemlis/internal/bech32"
)

type GeneratedSecret struct {
	privateKeyBytes, publicKeyBytes []byte
	recipient                       *age.X25519Identity
	numberOfShares, threshold       uint
	shares                          []Share
}

type Share struct {
	bytes []byte
}

func GenerateSecret(numberOfShares, threshold uint) (*GeneratedSecret, error) {
	if numberOfShares < 2 {
		return nil, errors.New("number of shares must be at least 2")
	}
	if threshold < 1 {
		return nil, errors.New("threshold must be at least 1")
	}
	if threshold > numberOfShares {
		return nil, errors.New("threshold must be less than or equal to the number of shares")
	}

	recipient, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipient key: %v", err)
	}
	_, publicKeyBytes, err := bech32.Decode(recipient.Recipient().String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}
	_, privateKeyBytes, err := bech32.Decode(recipient.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	shares := make([]Share, numberOfShares)
	sharesBytes, _ := SplitSecret(privateKeyBytes, threshold, numberOfShares)
	for shareIndex, shareBytes := range sharesBytes {
		shares[shareIndex] = Share{bytes: shareBytes}
	}

	secret := &GeneratedSecret{
		recipient:       recipient,
		publicKeyBytes:  publicKeyBytes,
		privateKeyBytes: privateKeyBytes,
		numberOfShares:  numberOfShares,
		threshold:       threshold,
		shares:          shares,
	}

	return secret, nil
}

func (s *GeneratedSecret) PublicKeyString() string {
	return s.recipient.Recipient().String()
}

func (s *GeneratedSecret) PrivateKeyString() string {
	return s.recipient.String()
}

func (s *GeneratedSecret) Threshold() uint {
	return s.threshold
}

func (s *GeneratedSecret) NumberOfShares() uint {
	return s.numberOfShares
}

func ShareIdentifier(share []byte) string {
	hash := sha256.Sum256(share)
	hashHex := hex.EncodeToString(hash[:])
	lastFiveChars := hashHex[len(hashHex)-5:]
	return lastFiveChars
}

func (s *Share) Identifier() string {
	return ShareIdentifier(s.bytes)
}

func (s *Share) Words() []string {
	return EncodeBytesToWords(s.bytes)
}

func (s *GeneratedSecret) Shares() []Share {
	return s.shares
}
