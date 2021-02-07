package conf

import (
	"io/ioutil"
	"time"
)

type Configuration struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Routes   RoutesConfig   `yaml:"routes"`
}

type ServerConfig struct {
	HostPort   string `yaml:"host_port"`
	HealthPort bool   `yaml:"health_port"`
}

type DatabaseConfig struct {
	Host                  string        `yaml:"host"`
	Port                  string        `yaml:"port"`
	Username              string        `yaml:"username"`
	Database              string        `yaml:"database"`
	SSLMode               string        `yaml:"ssl_mode"`
	SSLCert               string        `yaml:"ssl_cert"`
	SSLKey                string        `yaml:"ssl_key"`
	MaxIdleConns          int           `yaml:"max_idle_conns"`
	MaxOpenConns          int           `yaml:"max_open_conns"`
	MaxConnLifeTimeMinute time.Duration `yaml:"max_conn_life_time_minute"`
}

type RoutesConfig struct {
	getGuestListUrl    string `yaml:"getGuestListUri"`
	postGuestListUrl   string `yaml:"postGuestListUri"`
	deleteGuestListUrl string `yaml:"deleteGuestListUri"`
}

func GetConfig(configFile string) (*Configuration, error) {
	var configuration *Configuration

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}
