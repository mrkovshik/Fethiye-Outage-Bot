package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

func GetConfig() Config {
	if err := ReadConfigYML("config.yml"); err != nil {
		log.Fatalf("Failed init configuration %v", err)
	}
	return GetConfigInstance()
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

type SearchConfig struct {
	LevRatio string `yaml:"levratio"`
	SimRatio string `yaml:"simratio"`
}

type SchedulerConfig struct {
	FetchPeriod string `yaml:"fetch_period"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project         Project         `yaml:"project"`
	Database        Database        `yaml:"database"`
	CrawlersURL     CrawlersURL     `yaml:"crawlersurl"`
	SearchConfig    SearchConfig    `yaml:"searchconfig"`
	SchedulerConfig SchedulerConfig `yaml:"schedulerconfig"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
