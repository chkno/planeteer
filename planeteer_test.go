package main

import "testing"

func TestEncodeDecode(t *testing.T) {
	dims := []int{3, 2, 4, 17, 26, 15, 1, 2, 1}
	for i := 0; i < 318240; i++ {   // Product of dims
		addr := DecodeIndex(dims, i)
		for j := 0; j < len(dims); j++ {
			if addr[j] >= dims[j] {
				t.Fail()
			}
		}
		if EncodeIndex(dims, addr) != i {
			t.Fail()
		}
	}
}
