package main

import (
	"testing"
)

func TestKroetPkgEncodeDecode(t *testing.T) {

	kroetpkg := KroetPkg{"1", "10", "UP", "36", "10"}
	encoded := kroetpkg.Encode()
	decoded := DecodeKroetPkg(encoded)
	if kroetpkg != decoded {
		t.Errorf("Encoded and Decoded structs are not the same")
	}
}
