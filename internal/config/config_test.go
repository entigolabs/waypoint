package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_FromEnv(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")
	t.Setenv("SERVER_ADDR", ":9090")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, "postgres://user:pass@localhost/db", cfg.DBURI)
	assert.Equal(t, "https://api.example.com", cfg.APIBaseURL)
	assert.Equal(t, []string{"http://localhost:3000", "http://localhost:8080"}, cfg.AllowedOrigins)
	assert.Equal(t, ":9090", cfg.ServerAddr)
}

func TestLoadConfig_DefaultServerAddr(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, ":8081", cfg.ServerAddr)
}

func TestLoadConfig_MissingDBURI(t *testing.T) {
	t.Setenv("DB_URI", "")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_URI")
}

func TestLoadConfig_InvalidDBURI_WrongScheme(t *testing.T) {
	t.Setenv("DB_URI", "mysql://user:pass@localhost:3306/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scheme")
}

func TestLoadConfig_InvalidDBURI_MissingPassword(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user@localhost:5432/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password")
}

func TestLoadConfig_InvalidDBURI_MissingDatabase(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost:5432")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "database")
}

func TestLoadConfig_MissingAPIBaseURL(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost/db")
	t.Setenv("API_BASE_URL", "")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API_BASE_URL")
}

func TestLoadConfig_DefaultLogSettings(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, LogLevelInfo, cfg.LogLevel)
	assert.Equal(t, LogFormatJSON, cfg.LogFormat)
}

func TestLoadConfig_InvalidLogLevel(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("LOG_LEVEL", "verbose")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_LEVEL")
}

func TestLoadConfig_InvalidLogFormat(t *testing.T) {
	t.Setenv("DB_URI", "postgres://user:pass@localhost/db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("LOG_FORMAT", "xml")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_FORMAT")
}
