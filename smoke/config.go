package smoke

import (
	"encoding/json"
	"os"
)


type Config struct {
	Api 	string `json: "api"`
	ApiPort int `json: "api_port"`
	Host	string `json: "host"`
	Port      int `json: "port"`
	Proto	string `json: "proto"`
	TcpEnabled bool `json: "enable_tcp"`
	UdpEnabled bool `json: "enable_udp"`
}

var config *Config


func GetConfig() *Config {
	if config == nil {
		config = loadConfig()
	}
	return config
}

func loadConfig() *Config {
	config := newDefaultConfig()
	loadConfigFromJson(config)
	validateRequiredFields(config)
	return config
}

func newDefaultConfig() *Config {
	return &Config{
		Api: "",
		ApiPort: 80,
		Host: "",
		Port: 2003,
		Proto: "tcp",
		TcpEnabled: true,
		UdpEnabled: false,
	}
}


func validateRequiredFields(config *Config) {
	if config.Api == "" {
		panic("missing configuration 'api'")
	}

	if config.Host == "" {
		panic("missing configuration 'host'")
	}
}

// Loads the config from json into the supplied config object
func loadConfigFromJson(config *Config) {
	path := configPath()

	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}
}

func configPath() string {
	path := os.Getenv("SMOKE_TEST_CONFIG")
	if path == "" {
		panic("Must set $SMOKE_TEST_CONFIG to point to an integration config .json file.")
	}

	return path
}