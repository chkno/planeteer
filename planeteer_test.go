package main

import "testing"

func TestEncodeDecode(t *testing.T) {
	dims := []int{3, 2, 4, 17, 26, 15, 2, 1, 2, 1}
	var i PhysicalIndex
	for i = 0; i < 636480; i++ { // Product of dims
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

func TestCommas(t *testing.T) {
	cases := map[int32]string{
		1: "1",
		10: "10",
		100: "100",
		1000: "1,000",
		10000: "10,000",
		100000: "100,000",
		1000000: "1,000,000",
		1234567: "1,234,567",
		1000567: "1,000,567",
		1234000: "1,234,000",
		525000: "525,000",
	}
	for n, s := range cases {
		if Commas(n) != s {
			t.Error(n, "not", s)
		}
	}
}
