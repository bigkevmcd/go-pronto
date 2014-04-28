package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/v1/yaml"
	"launchpad.net/goose/identity"
)

type config struct {
	Credentials struct {
		AuthUrl    string `yaml:"auth-url"`
		TenantName string `yaml:"tenant-name"`
		Region     string
		Username   string
		Password   string
	}
	Container string
	Port      string
}

// CredentialsFromEnv creates and initializes the credentials from the
// environment variables.
func ConfigFromYaml(filename string) (*config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %v", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %v", err)
	}
	conf := new(config)
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, fmt.Errorf("processing config file: %v", err)
	}
	return conf, nil
}

func CredentialsFromConfig(conf *config) *identity.Credentials {
	return &identity.Credentials{
		URL:        conf.Credentials.AuthUrl,
		User:       conf.Credentials.Username,
		Secrets:    conf.Credentials.Password,
		Region:     conf.Credentials.Region,
		TenantName: conf.Credentials.TenantName,
	}
}
