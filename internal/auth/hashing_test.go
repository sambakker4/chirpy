package auth

import (
	"testing"
)

func TestHashing(t *testing.T) {
	cases := []string{
		"password123",
		"somethingelse",
		"asdf#@#FJADSKFLAJE3234",
		"asdfghjkl;",
		"",
	}

	for _, testcase := range cases {
		hash, err := HashPassword(testcase)
		if err != nil {
			t.Errorf("Error hashing password")
		}
		if err := CheckPasswordHash(hash, testcase); err != nil {
			t.Errorf("Hash does not match password")
		}
	}
}
