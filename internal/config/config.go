package config

import (
    "io/ioutil"

    "github.com/eujoy/data-dict/pkg"
    "gopkg.in/yaml.v2"
)

// Config describes the configuration of the service.
type Config struct {
    Application application `yaml:"application"`
}

// application describes the main details of the service.
type application struct {
    Authors []author `yaml:"authors"`
    Name    string   `yaml:"name"`
    Usage   string   `yaml:"usage"`
    Version string   `yaml:"version"`
}

// author describes the details of the authors of the tool.
type author struct {
    Name  string `yaml:"name"`
    Email string `yaml:"email"`
}

// New creates and returns a configuration object for the service.
func New(configFile string) (*Config, *pkg.Error) {
    var config *Config

    yamlBytes, err := ioutil.ReadFile(configFile)
    if err != nil {
        return nil, &pkg.Error{Err: err}
    }

    err = yaml.Unmarshal(yamlBytes, &config)
    if err != nil {
        return nil, &pkg.Error{Err: err}
    }

    return config, nil
}
