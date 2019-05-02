package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pakerfeldt/lovi/pkg/models"
	"gopkg.in/yaml.v2"
)

const defaultAckTimeout int = 15
const autoResolveAfterSeconds int = 900

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Settings() models.Settings {
	config := os.Getenv("CONFIG")
	if config == "" {
		config = "config.yaml"
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	return models.Settings{ConfigFile: config, Port: port}
}

func Parse(file string) models.Config {
	dat, err := ioutil.ReadFile(file)
	check(err)
	config := models.Config{}
	err = yaml.Unmarshal([]byte(string(dat)), &config)
	check(err)
	for i := 0; i < len(config.Policies); i++ {
		policy := &config.Policies[i]
		if policy.AckTimeoutSeconds == 0 {
			policy.AckTimeoutSeconds = defaultAckTimeout
		}
		if policy.AutoResolveAfterSeconds == 0 {
			policy.AutoResolveAfterSeconds = autoResolveAfterSeconds
		}
	}

	for _, policy := range config.Policies {
		fmt.Printf("AckTimeout: %d\n", policy.AckTimeoutSeconds)
	}

	return config
}

func Print(config models.Config) {
	d, err := yaml.Marshal(&config)
	check(err)
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}
