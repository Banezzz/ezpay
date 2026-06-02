package mdb

import "testing"

func TestNormalizeNetworkMapsLegacyBscName(t *testing.T) {
	if got := NormalizeNetwork(" BINANCE "); got != NetworkBsc {
		t.Fatalf("NormalizeNetwork() = %q, want %q", got, NetworkBsc)
	}
}

func TestNetworkAliasesIncludesLegacyBscName(t *testing.T) {
	got := NetworkAliases(NetworkBsc)
	want := []string{NetworkBsc, NetworkBscLegacy}
	if len(got) != len(want) {
		t.Fatalf("NetworkAliases() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("NetworkAliases() = %#v, want %#v", got, want)
		}
	}
}

func TestSameNetworkTreatsBscAndLegacyNameAsEquivalent(t *testing.T) {
	if !SameNetwork(NetworkBsc, NetworkBscLegacy) {
		t.Fatalf("SameNetwork(%q, %q) = false, want true", NetworkBsc, NetworkBscLegacy)
	}
}
