package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	HomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return HomeDir + "/" + configFileName, nil
}

func write(c Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configFilePath, jsonData, 0666); err != nil {
		return err
	}
	return nil
}

func Read() Config {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}
	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}
	}
	var config Config
	if err := json.Unmarshal(contents, &config); err != nil {
		return Config{}
	}
	return config
}

func (c Config) SetUser(UserName string) error {
	c.CurrentUserName = UserName
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}
