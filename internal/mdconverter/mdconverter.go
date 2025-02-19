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

type HTMLListInterface interface {
	PushLine(line string)
	ClearQueue()
}

type HTMLListQueue struct {
	UnLines []*string
	OrLines []*string
}

func (h *HTMLListQueue) PushLineUn(line *string) {
	h.UnLines = append(h.UnLines, line)
}

func (h *HTMLListQueue) ClearQueueUn() {
	h.UnLines = make([]*string, 0)
}

func (h *HTMLListQueue) GetQueueLengthUn() int {
	return len(h.UnLines)
}

func (h *HTMLListQueue) PushLineOr(line *string) {
	h.OrLines = append(h.OrLines, line)
}

func (h *HTMLListQueue) ClearQueueOr() {
	h.OrLines = make([]*string, 0)
}

func (h *HTMLListQueue) GetQueueLengthOr() int {
	return len(h.OrLines)
}

func (h *HTMLListQueue) FormList() string {
	resList := "<ul>\n"
	endList := "</ul>"
	for _, line := range h.UnLines {
		text := string([]rune(*line)[2:])
		*line = fmt.Sprintf("<li>%s</li>", text)
		resList += *line + "\n"
	}
	h.ClearQueueUn()
	return resList + endList
}

func (h *HTMLListQueue) FormOrderedList() string {
	resList := "<ol>\n"
	endList := "</ol>"
	for _, line := range h.OrLines {
		text := string([]rune(*line)[3:])
		*line = fmt.Sprintf("<li>%s</li>", text)
		resList += *line + "\n"
	}
	h.ClearQueueOr()
	return resList + endList
}

type MdConverter struct {
	Configs []conf.PreparedConfigs
	HTMLListQueue
}

func (m *MdConverter) ReplaceHeader(line *string) {
	re, err := regexp.Compile(`^#{1,6}\s(\W|\d|\w)+`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) {
		hashCnt := strings.Count(*line, "#")
		hTagStart := fmt.Sprintf("<h1 class=\"header-%d\">", hashCnt)
		*line = strings.Replace(*line, strings.Repeat("#", hashCnt)+" ", hTagStart, 1) + "</h1>"
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

func (m *MdConverter) ReplaceUnOrderedList(line *string, fin bool) (string, bool) {
	re, err := regexp.Compile(`^-{1}\s[(\W|\s|\d|\w^\-+?)]`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) && !fin {
		m.PushLineUn(line)
		return "", true
	} else if m.GetQueueLengthUn() != 0 {
		resList := m.FormList()
		return resList, false
	} else {
		return "", false
	}
}

func (m *MdConverter) ReplaceOrderedList(line *string, fin bool) (string, bool) {
	re, err := regexp.Compile(`^[0-9]+\.\s[(\W|\d|\w)]+`)
	if err != nil {
		log.Fatal(err)
	}
	if re.MatchString(*line) && !fin {
		m.PushLineOr(line)
		return "", true
	} else if m.GetQueueLengthOr() != 0 {
		resList := m.FormOrderedList()
		return resList, false
	} else {
		return "", false
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
		res, ok := m.ReplaceUnOrderedList(&line, fin)
		if ok {
			continue
		} else if res != "" {
			htmlContent += res + "\n"
		}

		res, ok = m.ReplaceOrderedList(&line, fin)
		if ok {
			continue
		} else if res != "" {
			htmlContent += res + "\n"
		}

		htmlContent += line + "\n"

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
