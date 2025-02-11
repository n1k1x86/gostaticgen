package mdconverter

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type MdConverterInterface interface {
	StartConverting(mdFiles []string, confPath string)
}

type ListStack struct {
	IsOrdered bool
	ListStack []string
}

type MdConverter struct {
	ListStack
}

func (m *MdConverter) StartConverting(mdFiles []string, confPath string) {
	for _, file := range mdFiles {
		bytes, err := os.ReadFile(confPath + "\\" + file)
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(bytes), "\\r\\n")
		fmt.Println(lines)
	}
}
