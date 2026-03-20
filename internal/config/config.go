package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type LogLevel string
type LogFormat string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	ServerAddr     string    `mapstructure:"SERVER_ADDR" required:"false"`
	AllowedOrigins []string  `mapstructure:"ALLOWED_ORIGINS" required:"true"`
	DBURI          string    `mapstructure:"DB_URI" required:"true"`
	APIBaseURL     string    `mapstructure:"API_BASE_URL" required:"true"`
	LogLevel       LogLevel  `mapstructure:"LOG_LEVEL" required:"false"`
	LogFormat      LogFormat `mapstructure:"LOG_FORMAT" required:"false"`
}

// NewLogger creates a slog.Logger configured from the given Config.
func NewLogger(cfg Config) *slog.Logger {
	var level slog.Level
	_ = level.UnmarshalText([]byte(cfg.LogLevel))

	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if cfg.LogFormat == LogFormatText {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

// LoadConfig reads configuration from file or environment variables.
// If a config.local.env file exists, it will be used instead of config.env.
// Exported environment variables will override the values in the config file.
func LoadConfig(path string) (config Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigType("env")
	configFile := ""
	if _, err := os.Stat(path + "/config.local.env"); err == nil {
		configFile = "config.local"
	} else if _, err := os.Stat(path + "/config.env"); err == nil {
		configFile = "config"
	}

	v.SetDefault("SERVER_ADDR", ":8081")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	if configFile != "" {
		v.SetConfigName(configFile)
		if err := v.ReadInConfig(); err != nil {
			return config, fmt.Errorf("error reading config file: %w", err)
		}
	}
	overrideConfig(v, config)
	var result map[string]interface{}
	if err := v.Unmarshal(&result); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
		),
		Result: &config,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return config, fmt.Errorf("error creating decoder: %w", err)
	}

	if err := decoder.Decode(result); err != nil {
		return config, fmt.Errorf("error decoding config: %w", err)
	}

	if err := validate(config); err != nil {
		return config, fmt.Errorf("config validation failed: %w", err)
	}
	return config, nil
}

func overrideConfig(v *viper.Viper, c interface{}) {
	value := reflect.ValueOf(c)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	t := value.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)
		if field.Type.Kind() == reflect.Struct {
			overrideConfig(v, fieldValue.Interface())
			continue
		}
		envName := field.Tag.Get("mapstructure")
		if value := os.Getenv(envName); value != "" {
			_ = v.BindEnv(envName)
			v.Set(envName, value)
		}
	}
}

func validate(c interface{}) error {
	v := reflect.ValueOf(c)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if field.Type.Kind() == reflect.Struct {
			if err := validate(fieldValue.Interface()); err != nil {
				return err
			}
			continue
		}

		value := fieldValue.Interface()
		if required := field.Tag.Get("required"); required != "false" && (value == "" || value == nil) {
			return fmt.Errorf("%s is required", field.Tag.Get("mapstructure"))
		}
		if err := validateField(field.Name, value); err != nil {
			return err
		}
	}
	return nil
}

func validateField(fieldName string, value interface{}) error {
	switch fieldName {
	case "APIBaseURL":
		return validateHTTPURL(fieldName, value.(string))
	case "DBURI":
		return validateDBURI(value.(string))
	case "LogFormat":
		switch value.(LogFormat) {
		case LogFormatJSON, LogFormatText:
		default:
			return fmt.Errorf("LOG_FORMAT must be '%s' or '%s'", LogFormatJSON, LogFormatText)
		}
	case "LogLevel":
		switch value.(LogLevel) {
		case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		default:
			return fmt.Errorf("LOG_LEVEL must be one of: %s, %s, %s, %s", LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
		}
	}
	return nil
}

func validateDBURI(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("DB_URI is not a valid URL: %w", err)
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("DB_URI scheme must be postgres:// or postgresql://")
	}
	if u.User == nil {
		return fmt.Errorf("DB_URI must include a username")
	}
	if _, ok := u.User.Password(); !ok {
		return fmt.Errorf("DB_URI must include a password")
	}
	if u.Host == "" {
		return fmt.Errorf("DB_URI must include a host")
	}
	if u.Path == "" || u.Path == "/" {
		return fmt.Errorf("DB_URI must include a database name")
	}
	return nil
}

func validateHTTPURL(fieldName, value string) error {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return fmt.Errorf("%s must be a valid URL starting with http:// or https://", fieldName)
	}
	return nil
}
