package tests

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"embox/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var userCounter int64

// noRedirectClient does not follow redirects, allowing tests to assert exact status codes.
var noRedirectClient = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// generateCSRFToken creates a valid CSRF token signed with the test secret.
func generateCSRFToken() string {
	nonce := make([]byte, 32)
	_, _ = rand.Read(nonce)
	encodedNonce := base64.RawURLEncoding.EncodeToString(nonce)
	mac := hmac.New(sha256.New, []byte(testCSRFSecret))
	mac.Write([]byte(encodedNonce))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedNonce + "." + sig
}

// CreateTestUser inserts a user directly into the DB, calls /auth/login to obtain a session,
// and returns (email, access_token cookie value).
func CreateTestUser(t *testing.T, db *gorm.DB, server *httptest.Server) (string, string) {
	t.Helper()
	n := atomic.AddInt64(&userCounter, 1)
	email := fmt.Sprintf("user%d@test.example.com", n)
	token := fmt.Sprintf("%06d", n%1000000)

	user := &models.User{
		Email:          email,
		Name:           "Test User",
		Token:          token,
		TokenCreatedAt: time.Now(),
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("CreateTestUser: db.Create: %v", err)
	}

	resp := doJSON(t, server, "POST", "/auth/login", fmt.Sprintf(`{"token":%q}`, token), "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("CreateTestUser: login returned %d: %s", resp.StatusCode, body)
	}
	io.Copy(io.Discard, resp.Body)

	for _, c := range resp.Cookies() {
		if c.Name == "access_token" {
			return email, c.Value
		}
	}
	t.Fatal("CreateTestUser: access_token cookie not set in response")
	return "", ""
}

// doJSON sends a JSON request with CSRF token and optional auth cookie.
func doJSON(t *testing.T, server *httptest.Server, method, path, body, cookie string) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, server.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("doJSON: NewRequest: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if method != http.MethodGet && method != http.MethodHead {
		req.Header.Set("X-XSRF-TOKEN", generateCSRFToken())
	}
	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: "access_token", Value: cookie})
	}
	resp, err := noRedirectClient.Do(req)
	if err != nil {
		t.Fatalf("doJSON: Do: %v", err)
	}
	return resp
}

// doMultipart sends a multipart/form-data POST with CSRF token and optional auth cookie.
func doMultipart(t *testing.T, server *httptest.Server, path string, buildForm func(*multipart.Writer), cookie string) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	buildForm(writer)
	if err := writer.Close(); err != nil {
		t.Fatalf("doMultipart: writer.Close: %v", err)
	}
	req, err := http.NewRequest("POST", server.URL+path, &body)
	if err != nil {
		t.Fatalf("doMultipart: NewRequest: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-XSRF-TOKEN", generateCSRFToken())
	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: "access_token", Value: cookie})
	}
	resp, err := noRedirectClient.Do(req)
	if err != nil {
		t.Fatalf("doMultipart: Do: %v", err)
	}
	return resp
}

// decodeJSON decodes the response body into a map.
func decodeJSON(t *testing.T, r io.Reader) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		t.Fatalf("decodeJSON: %v", err)
	}
	return result
}

// createTestMedia inserts a minimal Media record directly into the DB.
func createTestMedia(t *testing.T, db *gorm.DB, userID *uuid.UUID) *models.Media {
	t.Helper()
	media := &models.Media{
		Date:    time.Now(),
		UserID:  userID,
		FileExt: "jpg",
		Type:    "image",
		Caption: "test media",
	}
	if err := db.Create(media).Error; err != nil {
		t.Fatalf("createTestMedia: db.Create: %v", err)
	}
	return media
}

// getUserFromDB returns the User for the given email.
func getUserFromDB(t *testing.T, db *gorm.DB, email string) *models.User {
	t.Helper()
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		t.Fatalf("getUserFromDB(%q): %v", email, err)
	}
	return &user
}

// createTestPNG returns a minimal valid 10×10 PNG image.
func createTestPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
