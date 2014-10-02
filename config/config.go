package config

import (
    "errors"
    "io/ioutil"

    "github.com/coreos/go-etcd/etcd"
    "gopkg.in/yaml.v1"
)

type EtcdConfig struct {
    Location    string
    Consistency string
}

type Config struct {
    Etcd    EtcdConfig
    Workers int
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

func NewEtcdClient(configFile *Config) *etcd.Client {
    client := etcd.NewClient([]string{configFile.Etcd.Location})

    consistency := configFile.Etcd.Consistency
    if consistency != "" {
        client.SetConsistency(consistency)
    }

    return client
}
