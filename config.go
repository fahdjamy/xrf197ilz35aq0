package xrf197ilz35aq0

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"sync"
	"time"
)

const (
	maxRetries     = 3
	retryAfter     = 2 * time.Second
	configFilename = "configs/%s/config.yml"
)

var mutex = sync.Mutex{}
var configurations *Config

type Log struct {
	Level  string `yaml:"level"`
	Logger string `yaml:"logger"`
}

type MongoConfig struct {
	AppName          string `yaml:"appName"`
	RetryWrites      bool   `yaml:"retryWrites"`
	Uri              string `yaml:"uri"`
	Acknowledgment   string `yaml:"w"`
	DatabaseName     string `yaml:"databaseName"`
	DirectConnection bool   `yaml:"directConnection"`
}

type ApplicationConfig struct {
	Port int `yaml:"port"`

	IdleTimeoutSecs     time.Duration `yaml:"idleTimeoutSecs"`
	ReadTimeout         time.Duration `yaml:"readTimeoutSecs"`
	WriteTimeout        time.Duration `yaml:"writeTimeoutSecs"`
	GracefulTimeoutSecs time.Duration `yaml:"gracefulTimeoutSecs"`
}

type Database struct {
	Mongo MongoConfig `yaml:"mongo"`
}

type Config struct {
	Environment string            `yaml:"environment"`
	Log         Log               `yaml:"log"`
	Database    Database          `yaml:"database"`
	Application ApplicationConfig `yaml:"application"`
}

func NewConfig(env string) (Config, error) {
	if configurations == nil {
		// Acquire the lock to ensure strict singleton but only when creating a new config
		mutex.Lock()
		defer mutex.Unlock()

		yamlFile, err := readFromFile(fmt.Sprintf(configFilename, env))
		if err != nil {
			return Config{}, err
		}
		configurations, err := readConfiguration(yamlFile)
		if err != nil {
			return Config{}, err
		}
		return *configurations, nil
	}
	return *configurations, nil
}

// Open the YAML file
func readFromFile(filePath string) (io.ReadCloser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readConfiguration(file io.ReadCloser) (*Config, error) {
	defer func() {
		if err := CloseFileWithRetry(file, maxRetries, retryAfter); err != nil {
			fmt.Println(err)
		}
	}()

	// Decode the YAML into a struct
	var config Config

	// NewDecoder returns a new decoder that reads from r (file)
	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
