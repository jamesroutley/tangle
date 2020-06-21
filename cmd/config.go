package cmd

import (
	"encoding/json"
	"io/ioutil"
)

type Filter struct {
	Language string `json:"language"`
	Regex    string `json:"regex"`
}

type Target struct {
	Outfile string   `json:"outfile"`
	Sources []string `json:"sources"`
	Filters *Filter  `json:"filters"`
	Order   []string `json:"order"`
}

type Config struct {
	Targets []*Target `json:"targets"`
}

func readConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}

func defaultConfig(source string, outfile string) *Config {
	return &Config{
		Targets: []*Target{
			{
				Sources: []string{source},
				Outfile: outfile,
			},
		},
	}
}
