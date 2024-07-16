package main

import (
	"fmt"
	"testing"

	"github.com/filleokus/hemlis/internal/hemlis"
)

func TestExpectedWorkingRoundTrip(t *testing.T) {
	secret, _ := hemlis.GenerateSecret(6, 4)
	shares := secret.Shares()
	t.Logf("Generated secret: %s\n", secret.PrivateKeyString())
	var randomWords [][]string = make([][]string, 4)
	for i := 0; i < 4; i++ {
		randomWords[i] = shares[i].Words
	}
	combinedSecret, err := combineShares(randomWords)
	if err != nil {
		fmt.Println(err)
	} else {
		t.Logf("Combined secret:  %s\n", combinedSecret)
	}
	if secret.PrivateKeyString() != combinedSecret {
		t.Errorf("Combined secret is not equal to orginal secret")
		t.Errorf("Original secret: %v", secret.PrivateKeyString())
		t.Errorf("Combined secret: %v", combinedSecret)
	}
}
