package conf

import (
	"gopkg.in/yaml.v2"
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
	HealthPort string `yaml:"health_port"`
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
	GetGuestListUri    string `yaml:"getGuestListUri"`
	PostGuestListUri   string `yaml:"postGuestListUri"`
	DeleteGuestListUri string `yaml:"deleteGuestListUri"`

	GetGuestsUri    string `yaml:"getGuestsUri"`
	PutGuestsUri    string `yaml:"putGuestsUri"`
	DeleteGuestsUri string `yaml:"deleteGuestsUri"`

	GetInvitationUri string `yaml:"getInvitationUri"`

	GetEmptySeatsUri string `yaml:"getEmptySeats"`
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
