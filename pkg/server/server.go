package server

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHtmls(dir string) []string {
	fileNames := []string{}
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileParts := strings.Split(file.Name(), ".")
		if fileParts[1] != "html" {
			continue
		}
		fileNames = append(fileNames, fileParts[0])
	}

	return fileNames
}

func handlerFactory(file string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(200, file+".html", gin.H{
			"filename": file,
		})
	}
}

func RunServer(out string) {
	files := GetHtmls(out)
	r := gin.Default()
	r.Static("/css", out+"\\css")

	r.LoadHTMLGlob(out + "/*.html")
	for _, file := range files {
		path := "/" + file
		handler := handlerFactory(file)
		r.GET(path, handler)
	}
	fmt.Println(out)
	r.Run()
}
