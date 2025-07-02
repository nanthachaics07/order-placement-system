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
	GinMode = load_env.Require("GIN_MODE")
	ServiceName = load_env.Require("SERVICE_NAME")
	AppVersion = load_env.Require("APP_VERSION")
	LogLevel = load_env.Require("LOG_LEVEL")
	Port = load_env.Require("PORT")
	ShutdownTimeout, _ = time.ParseDuration(load_env.Default("SHUTDOWN_TIMEOUT", "5s"))
}
