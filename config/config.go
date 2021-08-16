package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Cluster struct {
		Owner string `yaml:"owner"`
		Nodes struct {
			Master []Node
			Slave  []Node
		}
	}
}

type Node struct {
	Address struct {
		External string `yaml:"external"`
		Internal string `yaml:"internal"`
	}
}

func Yaml() (*Config, error) {
	f, err := ioutil.ReadFile("./cluster.yml")
	if err != nil {
		return nil, fmt.Errorf("cluster.yml read failed")
	}
	cfg := &Config{}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cluster.yml unmarshal failed")
	}

	return cfg, nil
}
