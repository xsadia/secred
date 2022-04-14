package internal

import "testing"

func TestValidateAuthHeader(t *testing.T) {
	t.Run("Should return the token if header is set correctly", func(t *testing.T) {
		expected := "token"
		got, _ := ValidateAuthHeader("bearer token")

		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("Should return error if header is missing", func(t *testing.T) {

		_, err := ValidateAuthHeader("")

		if err == nil {
			t.Error("expected an error got none")
		}

		if err.Error() != "authorization header missing" {
			t.Errorf("Expected error to be %q, got %q", "authorization header missing", err.Error())
		}
	})

	t.Run("Should return error if header is not set correctly", func(t *testing.T) {
		_, err := ValidateAuthHeader("token")

		if err == nil {
			t.Error("expected an error got none")
		}

		if err.Error() != "invalid authorization header" {
			t.Errorf("Expected error to be %q, got %q", "invalid authorization header", err.Error())
		}
	})
}
