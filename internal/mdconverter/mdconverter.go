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

func (m *MdConverter) ReplaceHeader(line *string) {
	re, err := regexp.Compile(`^#{1,6}\s(\W|\d|\w)+`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) {
		hashCnt := strings.Count(*line, "#")
		hTagStart := fmt.Sprintf("<h1 class=\"header-%d\">", hashCnt)
		*line = strings.Replace(*line, strings.Repeat("#", hashCnt)+" ", hTagStart, 1) + "<\\h1>"
	}
}

func (m *MdConverter) ReplaceItalic(line *string) {
	re, err := regexp.Compile(`\*{1}([\w\sа-яА-Я^\*]+?)\*{1}`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) {
		*line = re.ReplaceAllString(*line, "<em class=\"italic-class\">$1</em>")
	}
}

func (m *MdConverter) ReplaceBold(line *string) {
	re, err := regexp.Compile(`\*{2}([\w\sа-яА-Я^\*]+?)\*{2}`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) {
		*line = re.ReplaceAllString(*line, "<b class=\"bold-class\">$1</b>")
	}
}

func (m *MdConverter) ConvertToHtml(data string) {
	for _, line := range strings.Split(data, "\r\n") {
		m.ReplaceHeader(&line)
		m.ReplaceBold(&line)
		m.ReplaceItalic(&line)
		fmt.Println("Line:", line)
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
