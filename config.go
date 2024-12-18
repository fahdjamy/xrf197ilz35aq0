package xrf197ilz35aq0

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"sync"
	"time"
	"xrf197ilz35aq0/internal"
)

const (
	maxRetries     = 3
	retryAfter     = 2 * time.Second
	configFilename = "configs/%s/config.yml"
)

var mutex = sync.Mutex{}
var configurations *Config

type Log struct {
	Level    string `yaml:"level"`
	Logger   string `yaml:"logger"`
	Filename string `yaml:"filename"`
}

type PasswordConfig struct {
	Time   uint8  `yaml:"time"`
	Thread uint8  `yaml:"thread"`
	Memory uint32 `yaml:"memory"`
}

type Security struct {
	PasswordConfig PasswordConfig `yaml:"passwordHash"`
}

type MongoConfig struct {
	AppName          string `yaml:"appName"`
	RetryWrites      bool   `yaml:"retryWrites"`
	Uri              string `yaml:"uri"`
	Acknowledgment   string `yaml:"w"`
	CloudUri         string `yaml:"cloudUri"`
	DatabaseName     string `yaml:"databaseName"`
	DirectConnection bool   `yaml:"directConnection"`
}

type ApplicationConfig struct {
	Port int `yaml:"port"`

	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	GracefulTimeout time.Duration `yaml:"gracefulTimeout"`
}

type Database struct {
	Mongo MongoConfig `yaml:"mongo"`
}

type Config struct {
	Environment string            `yaml:"environment"`
	Log         Log               `yaml:"log"`
	Database    Database          `yaml:"database"`
	Application ApplicationConfig `yaml:"application"`
	Security    Security          `yaml:"security"`
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
		if err := internal.CloseFileWithRetry(file, maxRetries, retryAfter); err != nil {
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
