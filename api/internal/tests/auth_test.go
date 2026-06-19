package tests

import (
	"io"
	"net/http"
	"testing"
	"time"

	"embox/internal/models"
)

func TestRequestLoginToken_ValidEmail(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	user := &models.User{Email: "token-test@example.com", Name: "Token Test"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("db.Create: %v", err)
	}

	resp := doJSON(t, server, "POST", "/auth/token", `{"email":"token-test@example.com"}`, "")
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	// SMTP is not configured in tests, so the endpoint may return 500 due to email failure.
	// The primary assertion is that GenerateToken ran and set the token in the DB.
	var updated models.User
	if err := db.First(&updated, user.ID).Error; err != nil {
		t.Fatalf("fetch updated user: %v", err)
	}
	if updated.Token == "" {
		t.Error("expected token to be set in DB after /auth/token, got empty")
	}
}

func TestRequestLoginToken_UnknownEmail(t *testing.T) {
	server, _, _, teardown := SetupTestApp(t)
	defer teardown()

	resp := doJSON(t, server, "POST", "/auth/token", `{"email":"nobody@example.com"}`, "")
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown email, got %d", resp.StatusCode)
	}
}

func TestVerifyToken_Valid(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	_, cookie := CreateTestUser(t, db, server)

	if cookie == "" {
		t.Error("expected access_token cookie after successful login")
	}
}

func TestVerifyToken_Expired(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	user := &models.User{
		Email:          "expired@example.com",
		Token:          "777777",
		TokenCreatedAt: time.Now().Add(-20 * time.Minute), // exceeds 10-minute limit
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("db.Create: %v", err)
	}

	resp := doJSON(t, server, "POST", "/auth/login", `{"token":"777777"}`, "")
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for expired token, got %d", resp.StatusCode)
	}
}

func TestVerifyToken_AlreadyConsumed(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	user := &models.User{
		Email:          "consumed@example.com",
		Token:          "888888",
		TokenCreatedAt: time.Now(),
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("db.Create: %v", err)
	}

	resp1 := doJSON(t, server, "POST", "/auth/login", `{"token":"888888"}`, "")
	defer resp1.Body.Close()
	io.Copy(io.Discard, resp1.Body)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("first login: expected 200, got %d", resp1.StatusCode)
	}

	// Token is now consumed — second attempt must fail.
	resp2 := doJSON(t, server, "POST", "/auth/login", `{"token":"888888"}`, "")
	defer resp2.Body.Close()
	io.Copy(io.Discard, resp2.Body)
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Errorf("second login with consumed token: expected 401, got %d", resp2.StatusCode)
	}
}
