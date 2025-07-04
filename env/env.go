package env

import (
	"order-placement-system/pkg/load_env"
	"time"
)

var (
	GinMode         string
	ServiceName     string
	AppVersion      string
	LogLevel        string
	Port            string
	ShutdownTimeout time.Duration
)

func LoadEnv() {
	// GinMode = load_env.Require("GIN_MODE")
	GinMode = load_env.Default("GIN_MODE", "release")
	ServiceName = load_env.Default("SERVICE_NAME", "order-placement-system")
	AppVersion = load_env.Default("APP_VERSION", "v1.0.5")
	LogLevel = load_env.Default("LOG_LEVEL", "dev")
	Port = load_env.Default("PORT", "8080")
	ShutdownTimeout, _ = time.ParseDuration(load_env.Default("SHUTDOWN_TIMEOUT", "5s"))
}
