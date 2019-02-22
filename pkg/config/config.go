package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pakerfeldt/lovi/pkg/models"
	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Parse(file string) models.Config {
	dat, err := ioutil.ReadFile(file)
	check(err)
	config := models.Config{}
	err = yaml.Unmarshal([]byte(string(dat)), &config)
	check(err)
	return config
}

func Print(config models.Config) {
	d, err := yaml.Marshal(&config)
	check(err)
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}
