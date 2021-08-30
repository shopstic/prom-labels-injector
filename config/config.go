package config

import (
	"log"
	"os"
	"shopstic/prom-labels-injector/util"
	"strings"
)

type HttpServerConfig struct {
	Port uint16
}

type PrometheusTargetConfig struct {
	Address string
}

type LabelConfig struct {
	Label string
	Value string
}

type LabelInjectorConfig struct {
	Labels []LabelConfig
}

type Config struct {
	Server           HttpServerConfig
	PrometheusTarget PrometheusTargetConfig
	LabelInjector    LabelInjectorConfig
}

func LoadFromEnvVars() Config {
	return Config{
		Server:           LoadHttpServerConfigFromEnvVars("SERVER_"),
		PrometheusTarget: LoadPrometheusTargetConfigFromEnvVars("PROMETHEUS_TARGET_"),
		LabelInjector:    LoadLabelInjectorConfigEnvVars("LABEL_INJECTOR_"),
	}
}

func LoadHttpServerConfigFromEnvVars(prefix string) HttpServerConfig {
	return HttpServerConfig{
		Port: LoadUint16(prefix+"PORT", "8080"),
	}
}

func LoadPrometheusTargetConfigFromEnvVars(prefix string) PrometheusTargetConfig {
	return PrometheusTargetConfig{
		Address: GetEnvOrFail(prefix + "ADDRESS"),
	}
}

func LoadUint16(key string, fallback string) uint16 {
	value := GetEnv(key, fallback)
	port, err := util.StringToUint16(value)
	if err != nil {
		log.Fatalf("Couldn't convert provided string [%v] to uint16 for env var [%v]", value, EnvKey(key))
	}
	return port
}

func LoadLabelInjectorConfigEnvVars(prefix string) LabelInjectorConfig {
	staticLabels := GetEnvArray(prefix + "LABELS")
	envVarValues := GetEnvArray(prefix + "ENV_VAR_VALUES")
	if len(staticLabels) != len(envVarValues) {
		log.Fatalf("Length of the provided static labels differs from the provided env var values. [LABELS=%+v] [ENV_VAR_VALUES=%+v]", staticLabels, envVarValues)
	}
	var labels []LabelConfig
	for i := range staticLabels {
		r := LabelConfig{
			Label: staticLabels[i],
			Value: GetRawEnvOrFail(envVarValues[i]),
		}
		labels = append(labels, r)
	}
	return LabelInjectorConfig{
		Labels: labels,
	}
}

func GetEnvArray(key string) []string {
	value := os.Getenv(EnvKey(key))
	return util.RemoveEmpty(strings.Split(value, " "))
}

// GetEnv is wrapping os.getenv with a fallback
func GetEnv(key string, fallback string) string {
	value := os.Getenv(EnvKey(key))
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvOrFail(key string) string {
	return GetRawEnvOrFail(EnvKey(key))
}

func GetRawEnvOrFail(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("Environment variable %s is required, but it was not defined", key)
	}
	return value
}

func EnvKey(key string) string {
	return "PROM_LABELS_INJECTOR_" + key
}
