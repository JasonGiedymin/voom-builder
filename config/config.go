package config

import (
    "errors"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

type EtcdPaths struct {
    SupervisorConfig string "supervisor_config"
    Services         string "services"
    Supervisors      string "supervisors"
}

type EtcdConfig struct {
    Location             string
    Consistency          string
    Service_ttl          uint64
    RegistrationInterval int       "registration_interval"
    Paths                EtcdPaths "paths"
}

type Config struct {
    Etcd      EtcdConfig
    Workers   int // number of workers to create
    WorkLimit int "worklimit" // number of test work to create
}

func (c *Config) ParseYaml(data []byte) error {
    if err := yaml.Unmarshal(data, &c); err != nil {
        msg := "Could not parse yaml config file." + err.Error()
        return errors.New(msg)
    }

    return nil
}

func ReadConfig(file string) (*Config, error) {
    var config Config

    blob, err := ioutil.ReadFile(file)

    if err != nil {
        return nil, err
    }

    if err := config.ParseYaml(blob); err != nil {
        return nil, err
    } else {
        return &config, nil
    }
}
