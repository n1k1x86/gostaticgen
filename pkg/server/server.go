package server

import (
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
		fileNames = append(fileNames, strings.Split(file.Name(), ".")[0])
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
	r.LoadHTMLGlob(out + "/*")
	for _, file := range files {
		path := "/" + file
		handler := handlerFactory(file)
		r.GET(path, handler)
	}
	r.Run()
}
