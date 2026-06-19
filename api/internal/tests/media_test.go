package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"embox/internal/models"

	"gorm.io/gorm"
)

func TestUploadMedia_Image(t *testing.T) {
	server, db, cfg, teardown := SetupTestApp(t)
	defer teardown()

	_, cookie := CreateTestUser(t, db, server)
	imgData := createTestPNG()

	meta := `[{"fileName":"test.png","type":"image/png","date":"2024-06-01T12:00:00Z","caption":"integration test"}]`

	resp := doMultipart(t, server, "/media/", func(w *multipart.Writer) {
		part, _ := w.CreateFormFile("files", "test.png")
		part.Write(imgData)
		w.WriteField("meta", meta)
	}, cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("upload returned %d: %s", resp.StatusCode, body)
	}

	// Decode response to get media ID
	var envelope struct {
		Data []struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(envelope.Data) == 0 {
		t.Fatal("expected at least one media item in response")
	}

	mediaID := uint(envelope.Data[0].ID)

	// DB entry must exist
	var media models.Media
	if err := db.First(&media, mediaID).Error; err != nil {
		t.Fatalf("media not found in DB: %v", err)
	}

	// Original file must be present in LocalStorageAdapter dir
	expectedPath := filepath.Join(cfg.Storage.LocalDir, media.RemotePath())
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("original file not found at %s: %v", expectedPath, err)
	}
}

func TestUploadMedia_WrongMIME(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	_, cookie := CreateTestUser(t, db, server)
	imgData := createTestPNG()

	// Declare file as video/mp4 but upload a PNG
	meta := `[{"fileName":"test.mp4","type":"video/mp4","date":"2024-06-01T12:00:00Z","caption":""}]`

	resp := doMultipart(t, server, "/media/", func(w *multipart.Writer) {
		part, _ := w.CreateFormFile("files", "test.mp4")
		part.Write(imgData)
		w.WriteField("meta", meta)
	}, cookie)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected non-200 for MIME mismatch, got 200")
	}
}

func TestGetMediaList(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email, cookie := CreateTestUser(t, db, server)
	user := getUserFromDB(t, db, email)
	createTestMedia(t, db, &user.ID)
	createTestMedia(t, db, &user.ID)

	resp := doJSON(t, server, "GET", "/media/", "", cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var envelope struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(envelope.Data) < 2 {
		t.Errorf("expected at least 2 media items, got %d", len(envelope.Data))
	}
}

func TestDeleteMedia_OwnedByUser(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email, cookie := CreateTestUser(t, db, server)
	user := getUserFromDB(t, db, email)
	media := createTestMedia(t, db, &user.ID)

	body := fmt.Sprintf(`{"ids":[%d]}`, media.ID)
	resp := doJSON(t, server, "DELETE", "/media/", body, cookie)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Media must be gone from DB
	err := db.First(&models.Media{}, media.ID).Error
	if err == nil {
		t.Error("media still present in DB after deletion")
	}
}

func TestDeleteMedia_OtherUser(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email1, _ := CreateTestUser(t, db, server)
	user1 := getUserFromDB(t, db, email1)
	media := createTestMedia(t, db, &user1.ID)

	// Create a second user and try to delete user1's media
	_, cookie2 := CreateTestUser(t, db, server)

	// Insert media owned by user1, but user2 has date in the past to avoid any overlap
	media.Date = time.Now().Add(-time.Hour) // just to have a valid non-zero date
	_ = db.Save(media)

	body := fmt.Sprintf(`{"ids":[%d]}`, media.ID)
	resp := doJSON(t, server, "DELETE", "/media/", body, cookie2)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when deleting another user's media, got %d", resp.StatusCode)
	}

	// Media must still exist
	if err := db.First(&models.Media{}, media.ID).Error; err != nil {
		t.Error("media was deleted by wrong user")
	}
}

// getMediaCount returns the number of media records in the DB.
func getMediaCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	db.Model(&models.Media{}).Count(&count)
	return count
}
