package reader

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func ReadMdFiles(dir string, output string) {
	log.Println("[INFO] - start reading files and transform to html")
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		mdContent, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		htmlConvert := blackfriday.Run(mdContent)
		fileName := strings.Split(file.Name(), ".")[0]
		err = os.WriteFile(filepath.Join(output, fileName+".html"), htmlConvert, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
