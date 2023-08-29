package test

import "testing"

func TestHello(t *testing.T) {
	if "Hello" == "hello" {
		t.Errorf("It should be case sensitive")
	}
}
