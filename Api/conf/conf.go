package conf

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Configuration struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Routes   RoutesConfig   `yaml:"routes"`

	Tables TableConfig `yaml:"tables"`
}

type ServerConfig struct {
	HostPort   string `yaml:"host_port"`
	HealthPort string `yaml:"health_port"`
}

type TableConfig struct {
	TableCapacityArray []int `yaml:"table_capacities"`
	TableCapacityMap   map[int]int
	TableCount         int
	TotalCapacity      int
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Driver   string `yaml:"driver"`
	Timeout  string `yaml:"timeout"`

	SSLMode               string        `yaml:"ssl_mode"`
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

	CountEmptySeatsUri string `yaml:"countEmptySeatsUri"`
}

func GetConfig(env string) (*Configuration, error) {
	fmt.Println("Reading Config")
	filepath := fmt.Sprintf("Api/config/%s.yml", env)
	configFilePath := flag.String("config", filepath, "Path to config file")

	var config *Configuration

	file, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	config.calculateTableValues()

	return config, nil
}

func (config *Configuration) calculateTableValues() {
	tableCount := 0
	totalCapacity := 0
	tableCapacityMap := make(map[int]int)

	for i, tableSize := range config.Tables.TableCapacityArray {
		tableCount++
		totalCapacity += tableSize
		tableCapacityMap[i] = tableSize
	}

	fmt.Println(fmt.Sprintf("There are %d tables, with a total capacity of %d", tableCount, totalCapacity))

	config.Tables.TableCount = tableCount
	config.Tables.TotalCapacity = totalCapacity
	config.Tables.TableCapacityMap = tableCapacityMap
}
