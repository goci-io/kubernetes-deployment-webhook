package main

import (
	"testing"
)

func TestContainsReturnTrueWhenArrayContainsItem(t *testing.T) {
	arr := []string{"a", "b", "c", "d"}
	result := contains(arr, "c")

	if !result {
		t.Error("expected array to contain item c")
	}
}

func TestRandStringBytesReturnsRandomStrings(t *testing.T) {
	r1 := randStringBytes(10)
	r2 := randStringBytes(10)

	if r1 == r2 {
		t.Error("random values must not be equal")
	}

	if len(r1) != 10 || len(r2) != 10 {
		t.Error("random values must be exactly 10 characters long")
	}
}

func TestGetEnvReturnsFallbackValue(t *testing.T) {
	val := getEnv("DOES_NOT_EXIST", "fallback")

	if val != "fallback" {
		t.Error("expected fallback value from getEnv")
	}
}
