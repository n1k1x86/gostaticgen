package mdconverter

import (
	conf "core/internal/yamlconverter"
	"fmt"
)

type MdConverterInterface interface {
	StartConverting(mdFiles []string, confPath string)
}

type MdConverter struct {
	Configs []conf.PreparedConfigs
}

func (m *MdConverter) StartConverting() {
	for _, config := range m.Configs {
		fmt.Println(config)
	}
}
