package auth

import "testing"

func TestCSRFServiceCreatesSessionBoundTokens(t *testing.T) {
	service, err := NewCSRFService("01234567890123456789012345678901")
	if err != nil {
		t.Fatalf("create service: %v", err)
	}
	token, err := service.Create("session-one")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}
	if err := service.Validate("session-one", token); err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if err := service.Validate("session-two", token); err == nil {
		t.Fatal("expected token to be rejected for another session")
	}
	if err := service.Validate("session-one", token+"tampered"); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}

func TestNewCSRFServiceRequiresStrongSecret(t *testing.T) {
	if _, err := NewCSRFService("short"); err == nil {
		t.Fatal("expected short secret to be rejected")
	}
}
