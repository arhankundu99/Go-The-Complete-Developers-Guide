package main

import "testing"

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 48 {
		t.Errorf("Expected 48, but got %v", len(d))
	}
}
