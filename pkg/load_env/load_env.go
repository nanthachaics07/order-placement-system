package load_env

import (
	"fmt"
	"log"
	"strings"
	"syscall"
)

var panic []string
var warning []string

func WarnIfEmpty(envName string, description ...string) string {
	env, found := syscall.Getenv(envName)
	if !found {
		message := fmt.Sprintf("%s env is empty, it may be needed.", envName)
		message = prependDescription(message, description)
		preWarn(message)
	}
	return env
}

func Default(envName string, defaultValue string) string {
	env, found := syscall.Getenv(envName)
	if !found {
		return defaultValue
	}
	return env
}

func Require(envName string, description ...string) string {
	env, found := syscall.Getenv(envName)
	if !found {
		message := fmt.Sprintf("%s env is required.", envName)
		message = prependDescription(message, description)
		prePanic(message)
	}
	return env
}

func Assert() {
	if len(warning) > 0 {
		log.Println(strings.Join(warning, "\n"))
	}

	if len(panic) > 0 {
		log.Panic(strings.Join(panic, "\n"))
	}

	resetState()
}

func prependDescription(message string, description []string) string {
	if len(description) > 0 {
		message = fmt.Sprintf("%s (%s)", message, description[0])
	}
	return message
}

func prePanic(message string) {
	panic = append(panic, message)
}

func preWarn(message string) {
	warning = append(warning, message)
}

func resetState() {
	panic = nil
	warning = nil
}
