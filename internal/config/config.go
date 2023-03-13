package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

func GetConfig() (Config, error) {
	if err := ReadConfigYML("config.yml"); err != nil {
		err = errors.Wrap(err, "Failed init configuration")
		return Config{}, err
	}
	return GetConfigInstance(), nil
}

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}
	return Config{}
}

// Database - contains all parameters database connection.
type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Migrations string `yaml:"migrations"`
	Name       string `yaml:"name"`
	SslMode    string `yaml:"sslmode"`
	Driver     string `yaml:"driver"`
}

type CrawlersURL struct {
	Muski string `yaml:"muski"`
	Aydem string `yaml:"aydem"`
}

type SchedulerConfig struct {
	FetchPeriod        string `yaml:"fetch_period"`
	StateCleanUpPeriod string `yaml:"map_clean_period"`
	AlertSendPeriod    string `yaml:"alert_send_period"`
}

type BotConfig struct {
	UserStateLifeTime int `yaml:"user_state_life_time"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
	CommitHash  string
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project         Project         `yaml:"project"`
	Database        Database        `yaml:"database"`
	CrawlersURL     CrawlersURL     `yaml:"crawlersurl"`
	SchedulerConfig SchedulerConfig `yaml:"schedulerconfig"`
	LoggerConfig    zap.Config      `yaml:"loggerconfig"`
	BotConfig       BotConfig       `yaml:"botconfig"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		err = errors.Wrap(err, "Error opening config file")
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		err = errors.Wrap(err, "Error decoding config file")
		return err
	}

	cfg.Project.CommitHash = commitHash

	return nil
}
