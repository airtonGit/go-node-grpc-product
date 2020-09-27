package main

import "testing"

func TestDiscount(t *testing.T) {
	pct, value, err := Discount("1", "1")
	if err != nil {
		t.Fatal("Err call discount err:", err)
	}

	t.Fatalf("Desc %v valor %v", pct, value)
}
