package mdconverter

import (
	conf "core/internal/yamlconverter"
	"fmt"
	"log"
	"os"
)

type MdConverterInterface interface {
	StartConverting(mdFiles []string, confPath string)
}

type MdConverter struct {
	Configs []conf.PreparedConfigs
}

func (m *MdConverter) ReadMds(config conf.PreparedConfigs) {
	confPath := config.ConfigPath
	for _, md := range config.Yaml.Content {
		data, err := os.ReadFile(confPath + "\\" + md.FileName + ".md")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	}
}

func (m *MdConverter) StartConverting() {
	for _, config := range m.Configs {
		m.ReadMds(config)
	}
}
