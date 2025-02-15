package yamlconverter

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type YamlProcessorInterface interface {
	ReadYamls()
}

type PreparedConfigs struct {
	Yaml       YamlConfig
	ConfigPath string
}

type YamlProcessor struct {
	ProcessedConfigs []PreparedConfigs
	Configs          []fs.DirEntry
	RootPath         string
}

type Block struct {
	FileName string `yaml:"filename"`
}

type YamlConfig struct {
	Content    []Block `yaml:"content"`
	Navigation string  `yaml:"navigation"`
}

func (y *YamlProcessor) ReadConfig(data []byte, pagePath string) {
	config := YamlConfig{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	prepConfig := PreparedConfigs{Yaml: config, ConfigPath: pagePath}
	y.ProcessedConfigs = append(y.ProcessedConfigs, prepConfig)
}

func (y *YamlProcessor) FindYaml(files []fs.DirEntry) (fs.DirEntry, error) {
	for _, file := range files {
		if strings.Split(file.Name(), ".")[1] == "yaml" {
			return file, nil
		}
	}
	return nil, errors.New("error yaml not found")
}

func (y *YamlProcessor) ReadYamls() {
	for _, pageDir := range y.Configs {
		if pageDir.IsDir() {
			pagePath := y.RootPath + "\\" + pageDir.Name()
			files, err := os.ReadDir(pagePath)
			if err != nil {
				log.Fatal(err)
			}
			file, err := y.FindYaml(files)
			if err != nil {
				log.Fatal(err)
			}
			confPath := pagePath + "\\" + file.Name()
			data, err := os.ReadFile(confPath)
			if err != nil {
				log.Fatal(err)
			}
			y.ReadConfig(data, pagePath)
		}
	}

}
