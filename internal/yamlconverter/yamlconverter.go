package yamlconverter

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlProcessorInterface interface {
	ReadYamls()
}

type YamlProcessor struct {
	ConfigsPath string
}

type Block struct {
	FileName string `yaml:"filename"`
}

type YamlConfig struct {
	Content    []Block `yaml:"content"`
	Navigation string  `yaml:"navigation"`
}

func (y *YamlProcessor) ReadConfig(data []byte) YamlConfig {
	config := YamlConfig{}
	fmt.Println(string(data))
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (y *YamlProcessor) ReadYamls() {
	files, err := os.ReadDir(y.ConfigsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		confPath := y.ConfigsPath + "\\" + file.Name()
		fmt.Println(confPath)
		data, err := os.ReadFile(confPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(y.ReadConfig(data))
	}
}
