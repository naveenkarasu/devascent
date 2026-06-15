package tui

import "testing"

// TestDigestMatches_NormalizationContract proves the cross-tool contract with a
// THROWAWAY phrase (not the secret): Go's sha256(TrimSpace(...)) must be byte-
// identical to the PowerShell `UTF8.GetBytes($s)` one-liner used to mint the
// stored digest. "testphrase123" → d8fab…c053 from that exact one-liner. If this
// passes, the real phrase triggers the overlay by construction; if a future edit
// breaks normalization (e.g. drops TrimSpace, lowercases, or changes encoding),
// this fails loudly instead of the cheat silently never firing.
func TestDigestMatches_NormalizationContract(t *testing.T) {
	const dummy = "d8fab52293600430741c6d907854586dc06681d0b6803ce471cf76a9a035c053"
	if !digestMatches("  testphrase123\n", dummy) {
		t.Fatal("Go sha256(TrimSpace) must equal the PowerShell UTF-8 hash of the same phrase")
	}
	if digestMatches("testphrase124", dummy) {
		t.Fatal("a different phrase must not match")
	}
}

// Ordinary submitted code must never accidentally open the navigation overlay.
func TestMatchesContentChecksum_RealCodeDoesNotTrigger(t *testing.T) {
	for _, code := range []string{
		"def solve(x):\n    return x + 1\n",
		"",
		"   \n\n",
		"// just a comment\n",
	} {
		if matchesContentChecksum(code) {
			t.Fatalf("ordinary submission must not match the checksum: %q", code)
		}
	}
}
