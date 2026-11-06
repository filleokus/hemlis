package hemlis

import "testing"

func TestGenerateSecretRejectsInvalidParameters(t *testing.T) {
	cases := []struct {
		name              string
		shares, threshold uint
	}{
		{"zero shares", 0, 0},
		{"one share", 1, 1},
		{"threshold one", 3, 1},
		{"threshold above shares", 3, 4},
		{"too many shares", 256, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			secret, err := GenerateSecret(c.shares, c.threshold)
			if err == nil {
				t.Fatalf("expected error for shares=%d threshold=%d", c.shares, c.threshold)
			}
			if secret != nil {
				t.Fatalf("expected nil secret on error")
			}
		})
	}
}

func TestGenerateSecretProducesNonEmptyShares(t *testing.T) {
	for _, p := range []struct{ shares, threshold uint }{{2, 2}, {5, 3}, {255, 2}} {
		secret, err := GenerateSecret(p.shares, p.threshold)
		if err != nil {
			t.Fatalf("shares=%d threshold=%d: unexpected error: %v", p.shares, p.threshold, err)
		}
		if got := uint(len(secret.Shares())); got != p.shares {
			t.Fatalf("expected %d shares, got %d", p.shares, got)
		}
		for i, share := range secret.Shares() {
			if len(share.Words()) != 33 {
				t.Fatalf("share %d has %d words, expected 33", i, len(share.Words()))
			}
		}
	}
}
