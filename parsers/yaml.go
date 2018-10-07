package parsers

import (
	YamlUtil "github.com/smallfish/simpleyaml"
	"io/ioutil"
)

type YamlParser struct {
}

func (parser YamlParser) Parse(filePath string) *YamlUtil.Yaml {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("unable to read file: " + filePath)
	}

	yaml, err := YamlUtil.NewYaml(content)
	if err != nil {
		panic(err)
	}

	return yaml
}
