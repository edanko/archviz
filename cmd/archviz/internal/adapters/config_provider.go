package adapters

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-faster/errors"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the configuration for the application.
type Config struct {
	Logger LoggerConfig `yaml:"logger" env-required:"true"`
	App    AppConfig    `yaml:"app"    env-required:"true"`
	Debug  DebugConfig  `yaml:"debug"  env-required:"true"`
	HTTP   HTTPConfig   `yaml:"http"   env-required:"true"`
	Render RenderConfig `yaml:"render" env-required:"true"`
}

// Validate checks if the config is valid.
func (c *Config) Validate() error {
	return nil
}

// ReadConfig reads the config file and returns a Config object.
func ReadConfig(path string) (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		help, _ := cleanenv.GetDescription(config, nil)
		const format = `error reading config file: %s\n%s`
		return nil, errors.Errorf(format, err, help)
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoggerConfig is a logger config.
type LoggerConfig struct {
	Level slog.Level `yaml:"level" env:"LOG_LEVEL" env-default:"INFO"`
}

// AppConfig is an application config.
type AppConfig struct {
	Environment     string        `yaml:"environment"     env:"APP_ENV"              env-default:"development"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout" env:"APP_SHUTDOWN_TIMEOUT" env-default:"5s"`
}

// DebugConfig is a debug config.
type DebugConfig struct {
	Host string `yaml:"host" env:"DEBUG_HOST" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env:"DEBUG_PORT" env-default:"4000"`
}

// Address returns the address of the Debug instance.
func (d *DebugConfig) Address() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}

// HTTPConfig is an HTTPConfig config.
type HTTPConfig struct {
	Host         string        `yaml:"host"         env:"HTTP_HOST"`
	Port         int           `yaml:"port"         env:"HTTP_PORT"                    env-default:"3000" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"readTimeout"  env:"HTTP_SERVER_READ_TIMEOUT"     env-default:"5s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env:"HTTP_SERVER_WRITE_TIMEOUT"    env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"  env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" env-default:"120s"`
}

// Address returns the address of the HTTP server.
func (h *HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

// RenderConfig is a render config.
type RenderConfig struct {
	Direction    string              `yaml:"direction"    env:"RENDER_DIRECTION"     env-default:"down"`
	LayoutEngine string              `yaml:"layoutEngine" env:"RENDER_LAYOUT_ENGINE" env-default:"dagre"`
	Minify       bool                `yaml:"minify"       env:"RENDER_MINIFY"        env-default:"true"`
	Theme        string              `yaml:"theme"        env:"RENDER_THEME"         env-default:"default"`
	ImageBundler bool                `yaml:"imageBundler" env:"RENDER_IMAGE_BUNDLER" env-default:"true"`
	Classes      []map[string]string `yaml:"classes"`
}
