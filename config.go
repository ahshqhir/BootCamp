package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerConfig ServerConfig `yaml:"serverConfig"`
	SQLConfig    SQLConfig    `yaml:"SQLConfig"`
}

type ServerConfig struct {
	RuleAddress   string `yaml:"ruleAddress"`
	TicketAddress string `yaml:"ticketAddress"`
	Port          int    `yaml:"port"`
}

type SQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DataBase string `yaml:"dataBase"`
	Schema   string `yaml:"schema"`
	Table    string `yaml:"table"`
}

func loadConfig(file string) (Config, error) {
	code, err := os.ReadFile(file)
	config := Config{}
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(code, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
