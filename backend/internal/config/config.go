package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	HealthCheckPath = "/health"
	MetricsPath     = "/metrics"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	ServerAddr     string    `mapstructure:"SERVER_ADDR" required:"false"`
	AllowedOrigins []string  `mapstructure:"ALLOWED_ORIGINS" required:"true"`
	APIBaseURL     string    `mapstructure:"API_BASE_URL" required:"true"`
	UserAgent      string    `mapstructure:"USER_AGENT" required:"false"`
	DBConfig       DBConfig  `mapstructure:",squash"`
	LogConfig      LogConfig `mapstructure:",squash"`
}

type DBConfig struct {
	User     string `mapstructure:"DB_USER" required:"true"`
	Password string `mapstructure:"DB_PASSWORD" required:"true"`
	Host     string `mapstructure:"DB_HOST" required:"true"`
	Port     int    `mapstructure:"DB_PORT" required:"true"`
	Name     string `mapstructure:"DB_NAME" required:"true"`
}

type LogConfig struct {
	LogLevel  LogLevel  `mapstructure:"LOG_LEVEL" required:"false"`
	LogFormat LogFormat `mapstructure:"LOG_FORMAT" required:"false"`
	LogOutput LogOutput `mapstructure:"LOG_OUTPUT" required:"false"`
	LogPath   string    `mapstructure:"LOG_PATH" required:"false"`
}

func (c *Config) Validate() error {
	if c.LogConfig.LogOutput == LogOutputFile && c.LogConfig.LogPath == "" {
		return fmt.Errorf("log output file path is required")
	}
	return nil
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
	v.SetDefault("LOG_LEVEL", LogLevelInfo)
	v.SetDefault("LOG_FORMAT", LogFormatOTEL)
	v.SetDefault("LOG_OUTPUT", LogOutputStdout)

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
			StringToIntHookFunc(),
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
	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("config validation failed: %w", err)
	}
	return config, nil
}

func StringToIntHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int {
			return data, nil
		}
		return strconv.Atoi(data.(string))
	}
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
	case "LogOutput":
		switch value.(LogOutput) {
		case LogOutputFile, LogOutputStdout:
		default:
			return fmt.Errorf("LOG_OUTPUT must be '%s' or '%s'", LogOutputFile, LogOutputStdout)
		}
	case "LogLevel":
		switch value.(LogLevel) {
		case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		default:
			return fmt.Errorf("LOG_LEVEL must be one of: %s, %s, %s, %s", LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
		}
	case "LogFormat":
		switch value.(LogFormat) {
		case LogFormatText, LogFormatJSON, LogFormatOTEL:
		default:
			return fmt.Errorf("LOG_FORMAT must be one of: %s, %s, %s", LogFormatText, LogFormatJSON, LogFormatOTEL)
		}
	case "Port":
		port, ok := value.(int)
		if !ok {
			return fmt.Errorf("DB_PORT must be an integer")
		}
		if port < 0 || port > 65535 {
			return fmt.Errorf("DB_PORT must be between 0 and 65535")
		}
	}
	return nil
}

func validateHTTPURL(fieldName, value string) error {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return fmt.Errorf("%s must be a valid URL starting with http:// or https://", fieldName)
	}
	return nil
}
