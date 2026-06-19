package tests

import (
	"net/http/httptest"
	"os/exec"
	"testing"

	"embox/internal/api/routes"
	"embox/internal/config"
	"embox/internal/models"
	"embox/internal/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testCSRFSecret = "test-csrf-secret"

// SetupTestApp starts a full in-memory test server backed by SQLite.
// Tests that exercise media uploads require ffmpeg in PATH; the test is skipped if not found.
func SetupTestApp(t *testing.T) (*httptest.Server, *gorm.DB, *config.ApiConfig, func()) {
	t.Helper()

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found in PATH — skipping integration tests")
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("SetupTestApp: open sqlite: %v", err)
	}
	db.Exec("PRAGMA foreign_keys = ON")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Media{},
		&models.Album{},
		&models.AlbumMedia{},
		&models.Favourite{},
	); err != nil {
		t.Fatalf("SetupTestApp: auto-migrate: %v", err)
	}

	storageDir := t.TempDir()
	mediaDir := t.TempDir()
	origMediaDir := services.MediaDir
	services.MediaDir = mediaDir
	t.Cleanup(func() { services.MediaDir = origMediaDir })

	cfg := &config.ApiConfig{
		Server: &config.ServerConfig{
			Host:        "localhost",
			Domain:      "localhost",
			CorsOrigins: []string{},
		},
		Router: &config.RouterConfig{
			ReleaseMode: "test",
			LogOutput:   "stdout",
		},
		Csrf: &config.CsrfConfig{
			Secret:   testCSRFSecret,
			MaxAge:   3600,
			Domain:   "localhost",
			IsSecure: false,
		},
		Auth: &config.AuthConfig{
			Domain:               "localhost",
			AccessSecret:         "test-access-secret",
			RefreshSecret:        "test-refresh-secret",
			AccessExpiration:     3600,
			RefreshExpiration:    86400,
			IsSecure:             false,
			LoginTokenExpiration: 600,
			LoginEmailSubject:    "Test Login",
			LoginEmailTemplate:   "<p>%s %s %d</p>",
		},
		Storage: &config.StorageConfig{
			Adapter:  "local",
			LocalDir: storageDir,
		},
		Email: &config.EmailConfig{
			Host:     "localhost",
			Port:     25,
			From:     "test@example.com",
			Password: "",
		},
	}

	router := routes.Init(db, cfg)
	server := httptest.NewServer(router)

	return server, db, cfg, func() { server.Close() }
}
