package utils

import (
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Interval time.Duration
	Postgres PostgresConfig
	Mist     MistConfig
	Grafana  GrafanaConfig
}

type GrafanaConfig struct {
	Url         string
	DashboardId int
	ApiKey      string
}

type PostgresConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

type MistConfig struct {
	MinActivationInterval time.Duration
	MinActiveDuration     time.Duration
	ActivationTemperature int
	MaxDurationPercent    int
	MiApiUrl              string
}

func NewConfig() *Config {
	return &Config{
		Interval: getEnvDuration("INTERVAL", "20s"),
		Mist: MistConfig{
			MiApiUrl:              getEnv("MIAPI_URL", ""),
			MinActivationInterval: getEnvDuration("MIN_ACTIVATION_INTERAVL", "1m"),
			MinActiveDuration:     getEnvDuration("MIN_ACTIVE_DURATION", "20s"),
			ActivationTemperature: getEnvInt("ACTIVATION_TEMPERATURE", "28"),
			MaxDurationPercent:    getEnvInt("MAX_DURATION_PERCENT", "50"),
		},
		Grafana: GrafanaConfig{
			Url:         getEnv("GRAFANA_URL", ""),
			ApiKey:      getEnv("GRAFANA_APIKEY", ""),
			DashboardId: getEnvInt("GRAFANA_DASHBOARD", ""),
		},
		Postgres: PostgresConfig{
			Host: getEnv("POSTGRES_HOST", ""),
			Port: getEnv("POSTGRES_PORT", ""),
			User: getEnv("POSTGRES_USER", ""),
			Pass: getEnv("POSTGRES_PASS", ""),
			DB:   getEnv("POSTGRES_DB", ""),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv("SCHEDULER_" + key); exists {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}
	log.Fatalf("%s is not defined", key)
	return ""
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	duration, err := time.ParseDuration(getEnv(key, defaultValue))
	if err != nil {
		log.Fatalf("Cannot parse duration env variable %s", key)
	}
	return duration
}

func getEnvInt(key string, defaultValue string) int {
	value, err := strconv.Atoi(getEnv(key, defaultValue))
	if err != nil {
		log.Fatalf("Cannot parse duration env variable %s", key)
	}
	return value
}

func getEnvUrl(key string, defaultValue string) *url.URL {
	value, err := url.Parse(getEnv(key, defaultValue))
	if err != nil {
		log.Fatalf("Cannot parse duration env variable %s", key)
	}
	return value
}
