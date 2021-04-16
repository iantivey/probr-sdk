package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// NewConfigDecoder reads the provided into a new yaml decoder
func NewConfigDecoder(path string) (*yaml.Decoder, error) {
	file, err := openConfigPath(path)
	if err != nil {
		return nil, err // err if path can not be validated
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file) // Init new YAML decode
	return decoder, nil
}

func openConfigPath(path string) (file *os.File, err error) {
	s, err := os.Stat(path)
	if err == nil && s.IsDir() {
		err = fmt.Errorf("'%s' is a directory; expected a normal file", path)
		return
	}
	file, err = os.Open(path)
	if err != nil {
		return
	}
	return
}
