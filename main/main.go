package main

import (
	"core/pkg/reader"
	"core/pkg/server"
	"flag"
	"log"
	"os"
)

func main() {
	markdownDir := flag.String("markdown_dir", "", "Full path to dir with mardowns")
	outDir := flag.String("output_dir", "", "Full path where to save htmls")

	flag.Parse()

	fileInfo, err := os.Stat(*markdownDir)
	if err != nil {
		log.Fatalf("No such file or directory: %s", *markdownDir)
	} else {
		if fileInfo.IsDir() {
			log.Println("[INFO] - directory was loaded")
		}
	}

	reader.MdToHtml(*markdownDir, *outDir)

	server.RunServer(*markdownDir, *outDir)
}
