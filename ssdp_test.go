package main

import (
	"testing"
)

func TestMakeMSearchString(t *testing.T) {
	expected := "M-SEARCH * HTTP/1.1\r\nHost:239.255.255.250:1900\r\nST:ssdp:all\r\nMan:\"ssdp:discover\"\r\nMX:3\r\n\r\n"
	x := makeMSearchString("ssdp:all")

	if x != expected {
		t.Errorf("actual = %v, expected = %v", x, expected)
	}
}
