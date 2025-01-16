package config

import (
	"encoding/json"
	"os"
	"path"
)

const configDirName = "gator"
const configFileName = "config.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(path.Join(dir, configDirName, configFileName))
	if err != nil {
		return Config{}, err
	}
	var conf Config
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}

func Write(config Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, configDirName, configFileName), data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (conf *Config) SetUser(user string) error {
	conf.CurrentUserName = user
	return Write(*conf)
}
