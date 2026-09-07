package kidwords

import (
	"bytes"
	"slices"
	"testing"
)

func TestSecretFingerPrinting(t *testing.T) {
	secret, err := NewSecret([]byte("testSecret"), 12, 6)
	if err != nil {
		t.Fatal(err)
	}
	fp := secret.GetFingerprint()
	if len(fp) != len(secret) {
		t.Fatalf("expected fingerprint length of %d, got %d", len(secret), len(fp))
	}

	if !secret.MatchFingerprint(fp) {
		t.Fatal("secret does not match its own finger print")
	}
	secret2 := slices.Delete(secret, 0, 3)
	secret2 = slices.Delete(secret2, 7, 9)
	if !secret2.MatchFingerprint(fp) {
		t.Fatal("secret with dropped shards does not match its parent's finger print")
	}

	another, err := NewSecret([]byte("testSecre8"), 12, 6)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(another.GetFingerprint(), fp) {
		t.Fatal("finger prints match when they should not")
	}
}
