package mdconverter

import (
	conf "core/internal/yamlconverter"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type MdConverterInterface interface {
	StartConverting(mdFiles []string, confPath string)
}

type MdConverter struct {
	Configs []conf.PreparedConfigs
}

func (m *MdConverter) IsHeaderPattern(line string) bool {
	re, err := regexp.Compile(`^#{1,6}\s(\W|\d|\w)+`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line)
}

func (m *MdConverter) ReplaceHeader(line string) string {
	hashCnt := strings.Count(line, "#")
	newLine := strings.Replace(line, strings.Repeat("#", hashCnt)+" ", "<h1 class=\"header-1\">", 1) + "<\\h1>"
	return newLine
}

func (m *MdConverter) ConvertToHtml(data string) {
	for _, line := range strings.Split(data, "\r\n") {
		if m.IsHeaderPattern(line) {
			newLine := m.ReplaceHeader(line)
			fmt.Println("Line:", newLine, "IsHeaderPattern?:", m.IsHeaderPattern(line))
		} else {
			fmt.Println("Line:", line, "IsHeaderPattern?:", m.IsHeaderPattern(line))
		}

	}
}

func (m *MdConverter) ReadMds(config conf.PreparedConfigs) {
	confPath := config.ConfigPath
	for _, md := range config.Yaml.Content {
		data, err := os.ReadFile(confPath + "\\" + md.FileName + ".md")
		if err != nil {
			log.Fatal(err)
		}
		m.ConvertToHtml(string(data))
	}
}

func (m *MdConverter) StartConverting() {
	for _, config := range m.Configs {
		m.ReadMds(config)
	}
}
