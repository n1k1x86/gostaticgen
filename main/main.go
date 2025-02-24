package main

import (
	"core/internal/mdconverter"
	"core/internal/yamlconverter"
	"core/pkg/server"
	"flag"
	"log"
	"os"
)

func main() {
	outDir := flag.String("output_dir", "", "Full path where to save htmls")
	yamls := flag.String("yamls", "", "Path to yaml configuration")

	flag.Parse()

	configDirs, err := os.ReadDir(*yamls)

	if err != nil {
		log.Fatal(err)
	}

	for _, conf := range configDirs {
		fileInfo, err := conf.Info()
		if err != nil {
			log.Fatalf("No such file or directory: %s", conf.Name())
		} else {
			if fileInfo.IsDir() {
				log.Printf("[INFO] - directory %s was loaded", conf.Name())
			}
		}
	}

	yaProc := yamlconverter.YamlProcessor{Configs: configDirs, RootPath: *yamls}
	yaProc.ReadYamls()

	mdConv := mdconverter.MdConverter{Configs: yaProc.ProcessedConfigs}
	mdConv.StartConverting(*outDir + "\\")

	server.RunServer(*outDir)
}
