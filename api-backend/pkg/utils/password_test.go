package utils

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("admin123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty string")
	}
	if hash == "admin123" {
		t.Fatal("HashPassword returned plaintext")
	}
}

func TestComparePassword(t *testing.T) {
	hash, _ := HashPassword("test123")

	if !ComparePassword(hash, "test123") {
		t.Fatal("ComparePassword should return true for correct password")
	}
	if ComparePassword(hash, "wrong") {
		t.Fatal("ComparePassword should return false for wrong password")
	}
	if ComparePassword(hash, "") {
		t.Fatal("ComparePassword should return false for empty password")
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	h1, _ := HashPassword("same")
	h2, _ := HashPassword("same")
	if h1 == h2 {
		t.Fatal("Two hashes of same password should differ (bcrypt uses random salt)")
	}
}
