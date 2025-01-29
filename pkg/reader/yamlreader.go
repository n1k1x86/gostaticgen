package reader

type Config struct {
	Js     JsConfig     `yaml:"js"`
	Styles StylesConfig `yaml:"styles"`
}

type JsConfig struct {
	Assets []string `yaml:"assets"`
}

type StylesConfig struct {
	Assets []string `yaml:"assets"`
}
