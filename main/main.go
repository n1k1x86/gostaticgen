package main

import (
	"core/internal/yamlconverter"
	"flag"
	"log"
	"os"
)

func main() {
	// outDir := flag.String("output_dir", "", "Full path where to save htmls")
	yamls := flag.String("yamls", "", "Path to yaml configuration")

	flag.Parse()

	configDirs, err := os.ReadDir(*yamls)

	if err != nil {
		log.Fatal(err)
	}

	for _, conf := range configDirs {
		fileInfo, err := conf.Info()
		if err != nil {
			log.Fatalf("No such file or directory: %s", conf)
		} else {
			if fileInfo.IsDir() {
				log.Printf("[INFO] - directory %s was loaded", conf)
			}
		}
	}

	yaProc := yamlconverter.YamlProcessor{Configs: configDirs, RootPath: *yamls}
	yaProc.ReadYamls()

	// server.RunServer(*markdownDir, *outDir)
}
