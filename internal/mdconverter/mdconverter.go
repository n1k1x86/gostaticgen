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

func (m *MdConverter) IsItalicPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`\*{1}([\w\sа-яА-Я^\*]+?)\*{1}`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line), re
}

func (m *MdConverter) IsBoldPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`\*{2}(\w|\s)+\*{2}`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line), re
}

func (m *MdConverter) IsUnOrderedListPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`^-{1}\s(\W|\s|\d)+`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line), re
}

func (m *MdConverter) IsOrderedListPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`^[0-9]+\.\s(\W|\d|\w)+`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line), re
}

func (m *MdConverter) IsLinkPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`\[(\W|\w|\D)+\]\((\W|\w|\D)+\)`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line), re
}

func (m *MdConverter) IsImgPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`\!\[(\W|\w|\D)+\]\((\W|\w|\D)+\)`)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(line), re
}

func (m *MdConverter) ReplaceHeader(line string) string {
	hashCnt := strings.Count(line, "#")
	newLine := strings.Replace(line, strings.Repeat("#", hashCnt)+" ", "<h1 class=\"header-1\">", 1) + "<\\h1>"
	return newLine
}

func (m *MdConverter) ReplaceItalic(line string, re *regexp.Regexp) string {
	var newLine string = ""
	newLine = re.ReplaceAllString(line, "<em class=\"italic-class\">$1</em>")
	return newLine
}

func (m *MdConverter) ConvertToHtml(data string) {
	for _, line := range strings.Split(data, "\r\n") {
		if ok := m.IsHeaderPattern(line); ok {
			newLine := m.ReplaceHeader(line)
			fmt.Println("Line:", newLine, "IsHeaderPattern?:", ok)
		} else if ok, re := m.IsItalicPattern(line); ok {
			newLine := m.ReplaceItalic(line, re)
			fmt.Println("Line:", newLine, "IsItalicPattern?:", ok)
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
