package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"embox/internal/models"
)

func TestCreateAlbum(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	_, cookie := CreateTestUser(t, db, server)

	resp := doJSON(t, server, "POST", "/album/", `{"name":"My Album","description":"test"}`, cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
	}

	var envelope struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	albumID := uint(envelope.Data.ID)
	var album models.Album
	if err := db.First(&album, albumID).Error; err != nil {
		t.Fatalf("album not found in DB: %v", err)
	}
	if album.Name != "My Album" {
		t.Errorf("expected name %q, got %q", "My Album", album.Name)
	}
}

func TestUpdateAlbum_OwnedByUser(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email, cookie := CreateTestUser(t, db, server)
	user := getUserFromDB(t, db, email)

	album := &models.Album{Name: "Original", UserID: &user.ID}
	if err := db.Create(album).Error; err != nil {
		t.Fatalf("db.Create album: %v", err)
	}

	resp := doJSON(t, server, "PUT", fmt.Sprintf("/album/%d", album.ID),
		`{"name":"Updated"}`, cookie)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updated models.Album
	if err := db.First(&updated, album.ID).Error; err != nil {
		t.Fatalf("fetch updated album: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected name %q, got %q", "Updated", updated.Name)
	}
	if updated.UpdatedByID == nil {
		t.Error("expected updated_by_id to be set after update")
	}
}

func TestUpdateAlbum_OtherUser(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email1, _ := CreateTestUser(t, db, server)
	user1 := getUserFromDB(t, db, email1)

	album := &models.Album{Name: "User1 Album", UserID: &user1.ID}
	if err := db.Create(album).Error; err != nil {
		t.Fatalf("db.Create album: %v", err)
	}

	_, cookie2 := CreateTestUser(t, db, server)

	resp := doJSON(t, server, "PUT", fmt.Sprintf("/album/%d", album.ID),
		`{"name":"Hijacked"}`, cookie2)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 when updating another user's album, got %d", resp.StatusCode)
	}
}

func TestDeleteAlbum_CascadesMedia(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email, cookie := CreateTestUser(t, db, server)
	user := getUserFromDB(t, db, email)

	album := &models.Album{Name: "To Delete", UserID: &user.ID}
	if err := db.Create(album).Error; err != nil {
		t.Fatalf("db.Create album: %v", err)
	}

	media := createTestMedia(t, db, &user.ID)
	albumMedia := &models.AlbumMedia{AlbumID: album.ID, MediaID: media.ID}
	if err := db.Create(albumMedia).Error; err != nil {
		t.Fatalf("db.Create album_media: %v", err)
	}

	resp := doJSON(t, server, "DELETE", fmt.Sprintf("/album/%d", album.ID), "", cookie)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Album must be gone
	if err := db.First(&models.Album{}, album.ID).Error; err == nil {
		t.Error("album still present in DB after deletion")
	}

	// album_media entries must be cascade-deleted
	var count int64
	db.Model(&models.AlbumMedia{}).Where("album_id = ?", album.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 album_media entries after cascade delete, got %d", count)
	}
}

func TestAddMediaToAlbum(t *testing.T) {
	server, db, _, teardown := SetupTestApp(t)
	defer teardown()

	email, cookie := CreateTestUser(t, db, server)
	user := getUserFromDB(t, db, email)

	album := &models.Album{Name: "With Media", UserID: &user.ID}
	if err := db.Create(album).Error; err != nil {
		t.Fatalf("db.Create album: %v", err)
	}

	media := createTestMedia(t, db, &user.ID)

	body := fmt.Sprintf(`{"mediaIds":[%d],"isCover":false}`, media.ID)
	resp := doJSON(t, server, "POST", fmt.Sprintf("/album/%d/media", album.ID), body, cookie)
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var count int64
	db.Model(&models.AlbumMedia{}).Where("album_id = ? AND media_id = ?", album.ID, media.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 album_media entry, got %d", count)
	}
}
