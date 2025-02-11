package main

import (
	"core/internal/yamlconverter"
	"flag"
	"log"
	"os"
)

func main() {
	configDir := flag.String("markdown_dir", "", "Full path to dir with mardowns")
	// outDir := flag.String("output_dir", "", "Full path where to save htmls")
	configPath := flag.String("yaml", "", "Path to yaml configuration")

	flag.Parse()

	fileInfo, err := os.Stat(*configDir)
	if err != nil {
		log.Fatalf("No such file or directory: %s", *configDir)
	} else {
		if fileInfo.IsDir() {
			log.Println("[INFO] - directory was loaded")
		}
	}

	// reader.MdToHtml(*markdownDir, *outDir)
	yaProc := yamlconverter.YamlProcessor{ConfigsPath: *configPath}
	yaProc.ReadYamls()
	// reader.ReadYaml(*configPath)

	// server.RunServer(*markdownDir, *outDir)
}
