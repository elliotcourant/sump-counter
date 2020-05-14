package service

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	configFilePath = "CONFIG_PATH"
)

type (
	Configuration struct {
		SpreadsheetId    string       `yaml:"spreadsheetId"`
		SheetName        string       `yaml:"sheetName"`
		DateFormat       string       `yaml:"dateFormat"`
		TimeZone         string       `yaml:"timezone"`
		ValueInputOption string       `yaml:"valueInputOption"`
		LogLevel         logrus.Level `yaml:"logLevel"`
	}
)

func GetConfigurationFromEnv() (*Configuration, error) {
	path := os.Getenv(configFilePath)
	if len(path) == 0 {
		return nil, errors.Errorf("%s environment variable is not set", configFilePath)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read config file: %s", path)
	}

	var config Configuration
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to parse config file: %s", path)
	}

	return &config, nil
}
