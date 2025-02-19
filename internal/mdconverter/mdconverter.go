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

type UnOrderedListInterface interface {
	PushLine(line string)
	ClearQueue()
}

type UnOrderedListQueue struct {
	Lines []*string
}

func (u *UnOrderedListQueue) PushLine(line *string) {
	u.Lines = append(u.Lines, line)
}

func (u *UnOrderedListQueue) ClearQueue() {
	u.Lines = make([]*string, 0)
}

func (u *UnOrderedListQueue) GetQueueLength() int {
	return len(u.Lines)
}

func (u *UnOrderedListQueue) FormList() string {
	resList := "<ul>\n"
	endList := "</ul>"
	for _, line := range u.Lines {
		text := string([]rune(*line)[2:])
		*line = fmt.Sprintf("<li>%s</li>", text)
		resList += *line + "\n"
	}
	u.ClearQueue()
	return resList + endList
}

type MdConverter struct {
	Configs []conf.PreparedConfigs
	UnOrderedListQueue
}

func (m *MdConverter) IsOrderedListPattern(line string) (bool, *regexp.Regexp) {
	re, err := regexp.Compile(`^[0-9]+\.\s(\W|\d|\w)+`)
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

func (m *MdConverter) ReplaceLink(line *string) {
	re, err := regexp.Compile(`\[(\W|\w|\D+?)\]\((\W|\w|\D+?)\)`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) {
		aTag := "<a href=\"$2\" class=\"link-class\">$1</a>"
		*line = re.ReplaceAllString(*line, aTag)
	}
}

func (m *MdConverter) ReplaceImg(line *string) {
	re, err := regexp.Compile(`\!\[([\W|\w|\D]+?)\]\((\W|\w|\D+?)\)`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) {
		imgTag := "<img src=\"$2\" class=\"img-class\" alt=\"$1\" />"
		*line = re.ReplaceAllString(*line, imgTag)
	}
}

func (m *MdConverter) ReplaceUnOrderedList(line *string, fin bool) (bool, string) {
	re, err := regexp.Compile(`^-{1}\s[(\W|\s|\d|\w^\-+?)]`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) && !fin {
		m.PushLine(line)
		return true, ""
	} else if m.GetQueueLength() != 0 {
		resList := m.FormList()
		return false, resList
	} else {
		return false, *line
	}
}

func (m *MdConverter) ConvertToHtml(data string) {
	htmlContent := ""
	rows := strings.Split(data, "\r\n")
	for ind, line := range rows {
		fin := ind == len(rows)-1
		m.ReplaceHeader(&line)
		m.ReplaceBold(&line)
		m.ReplaceItalic(&line)
		m.ReplaceImg(&line)
		m.ReplaceLink(&line)
		if ok, res := m.ReplaceUnOrderedList(&line, fin); ok {
			continue
		} else {
			htmlContent += res + "\n"
		}
	}
	fmt.Println(htmlContent)
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
