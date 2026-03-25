package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setDBEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "pass")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "db")
}

func TestLoadConfig_FromEnv(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")
	t.Setenv("SERVER_ADDR", ":9090")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, "user", cfg.DBConfig.User)
	assert.Equal(t, "pass", cfg.DBConfig.Password)
	assert.Equal(t, "localhost", cfg.DBConfig.Host)
	assert.Equal(t, 5432, cfg.DBConfig.Port)
	assert.Equal(t, "db", cfg.DBConfig.Name)
	assert.Equal(t, "https://api.example.com", cfg.APIBaseURL)
	assert.Equal(t, []string{"http://localhost:3000", "http://localhost:8080"}, cfg.AllowedOrigins)
	assert.Equal(t, ":9090", cfg.ServerAddr)
}

func TestLoadConfig_DefaultServerAddr(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, ":8081", cfg.ServerAddr)
}

func TestLoadConfig_MissingDBUser(t *testing.T) {
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "pass")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_USER")
}

func TestLoadConfig_InvalidDBPort(t *testing.T) {
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "pass")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "99999")
	t.Setenv("DB_NAME", "db")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PORT")
}

func TestLoadConfig_MissingAPIBaseURL(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("SERVER_ADDR", "")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API_BASE_URL")
}

func TestLoadConfig_DefaultLogSettings(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")

	cfg, err := LoadConfig(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, LogLevelInfo, cfg.LogConfig.LogLevel)
	assert.Equal(t, LogOutputStdout, cfg.LogConfig.LogOutput)
}

func TestLoadConfig_InvalidLogLevel(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("LOG_LEVEL", "verbose")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_LEVEL")
}

func TestLoadConfig_InvalidLogFormat(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("LOG_FORMAT", "xml")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_FORMAT")
}

func TestLoadConfig_InvalidLogOutput(t *testing.T) {
	setDBEnv(t)
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	t.Setenv("LOG_OUTPUT", "syslog")

	_, err := LoadConfig(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_OUTPUT")
}
