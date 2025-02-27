package mdconverter

import (
	conf "core/internal/yamlconverter"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

func (m *MdConverter) ConvertToHtml(data string) string {
	blockContent := "<div class=\"block-class\">"
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
			blockContent += res + "\n"
		}

		res, ok = m.ReplaceOrderedList(&line, fin)
		if ok {
			continue
		} else if res != "" {
			blockContent += res + "\n"
		}

		blockContent += line + "\n"

	}
	return blockContent + "\n</div>\n"
}

func (m *MdConverter) FinishPage(title string, body string) string {
	finishPage := fmt.Sprintf(`<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>%s</title>
		<link rel="stylesheet" href="/css/styles.css">
	</head>
	<body>
		<header class="site-header">
			<div class="container">
				<div class="logo">
					<a href="#">GoWiki</a>
				</div>
			</div>
		</header>
		%s
	</body>
	</html>`, title, body)
	return finishPage
}

func (m *MdConverter) BuildHtml(config conf.PreparedConfigs) string {
	htmlContent := ""
	confPath := config.ConfigPath
	for _, md := range config.Yaml.Content {
		data, err := os.ReadFile(confPath + "\\" + md.FileName + ".md")
		if err != nil {
			log.Fatal(err)
		}
		htmlContent += m.ConvertToHtml(string(data))
	}
	return htmlContent
}

func (m *MdConverter) SaveHtml(htmlContent string, outDir string, fileName string) error {
	err := os.WriteFile(outDir+fileName+".html", []byte(htmlContent), fs.FileMode(os.O_RDWR))
	if err != nil {
		return err
	}
	return nil
}

func (m *MdConverter) CreateCss(outDir string) error {
	css := `/* Общие стили для body */
	body {
		font-family: 'Arial', sans-serif;
		line-height: 1.6;
		color: #333;
		background-color: #f9f9f9;
		margin: 0;
		padding: 20px;
	}
	
	/* Стили для блока с контентом */
	.block-class {
		background-color: #ffffff;
		border-radius: 8px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		padding: 20px;
		margin: 20px 0;
		max-width: 800px;
		margin-left: auto;
		margin-right: auto;
	}
	
	/* Стили для изображения */
	.img-class {
		max-width: 100%;
		height: auto;
		border-radius: 8px;
		display: block;
		margin: 20px auto;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}
	
	/* Стили для ссылки */
	.link-class {
		color: #007bff;
		text-decoration: none;
		font-weight: bold;
		transition: color 0.3s ease;
	}
	
	.link-class:hover {
		color: #0056b3;
		text-decoration: underline;
	}
	
	/* Стили для жирного текста */
	.bold-class {
		font-weight: bold;
		color: #222;
	}
	
	/* Стили для курсива */
	.italic-class {
		font-style: italic;
		color: #555;
	}
	
	/* Стили для заголовков */
	.header-1 {
		font-size: 2.5rem;
		color: #222;
		margin-bottom: 20px;
		font-weight: bold;
	}
	
	.header-2 {
		font-size: 2rem;
		color: #222;
		margin-bottom: 18px;
		font-weight: bold;
	}
	
	.header-3 {
		font-size: 1.75rem;
		color: #222;
		margin-bottom: 16px;
		font-weight: bold;
	}
	
	.header-4 {
		font-size: 1.5rem;
		color: #222;
		margin-bottom: 14px;
		font-weight: bold;
	}
	
	.header-5 {
		font-size: 1.25rem;
		color: #222;
		margin-bottom: 12px;
		font-weight: bold;
	}
	
	.header-6 {
		font-size: 1rem;
		color: #222;
		margin-bottom: 10px;
		font-weight: bold;
	}
	/* General Styles */

	.container {
		width: 90%;
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	/* Header Styles */
	.site-header {
		background-color: #333;
		color: #fff;
		padding: 20px 0;
		position: sticky;
		top: 0;
		z-index: 1000;
	}

	.logo a {
		color: #fff;
		text-decoration: none;
		font-size: 24px;
		font-weight: bold;
	}

	.nav-links {
		list-style: none;
		display: flex;
		margin: 0;
		padding: 0;
	}

	.nav-links li {
		margin-left: 20px;
	}

	.nav-links a {
		color: #fff;
		text-decoration: none;
		font-size: 16px;
		transition: color 0.3s ease;
	}

	.nav-links a:hover {
		color: #007bff;
	}

	.cta .btn {
		background-color: #007bff;
		color: #fff;
		padding: 10px 20px;
		border-radius: 5px;
		text-decoration: none;
		font-size: 16px;
		transition: background-color 0.3s ease;
	}

	.cta .btn:hover {
		background-color: #0056b3;
	}`
	err := os.Mkdir(outDir+"\\css", fs.FileMode(os.O_RDWR))
	if err != nil {
		return err
	}
	err = os.WriteFile(outDir+"\\css\\styles.css", []byte(css), fs.FileMode(os.O_RDWR))
	if err != nil {
		return err
	}
	return nil
}

func (m *MdConverter) ClearOutDir(dir *os.File, dirPath string) {
	names, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range names {
		filename := filepath.Join(dirPath, name)
		err := os.RemoveAll(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Print("[INFO] - output dir was cleared successfully")
}

func (m *MdConverter) IsDirExist(outDir string, createOut bool) {
	dir, err := os.Open(outDir)
	if err != nil && !createOut {
		log.Fatal(err)
	} else if createOut && err != nil {
		err := os.Mkdir(outDir, fs.FileMode(os.O_RDWR))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer dir.Close()
	if dir != nil {
		m.ClearOutDir(dir, outDir)
	}
}

func (m *MdConverter) StartConverting(outDir string, createOut bool) {
	m.IsDirExist(outDir, createOut)
	for _, config := range m.Configs {
		fileNameWords := strings.Split(config.ConfigPath, "\\")
		fileName := fileNameWords[len(fileNameWords)-1]
		content := m.BuildHtml(config)
		htmlPage := m.FinishPage(fileName, content)
		err := m.SaveHtml(htmlPage, outDir, fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := m.CreateCss(outDir)
	if err != nil {
		log.Fatal(err)
	}
}
