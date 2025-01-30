package reader

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Js     JsConfig     `yaml:"js"`
	Styles StylesConfig `yaml:"styles"`
}

type JsConfig struct {
	Assets []string `yaml:"assets"`
}

type StylesConfig struct {
	Assets []string `yaml:"assets"`
}

func ReadYaml(path string) {
	conf := Config{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(bytes, &conf)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}
