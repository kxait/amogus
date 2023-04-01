package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Mode string

const (
	// for lines in a file, attempt to crack every sha512 hash
	Sha512 Mode = "sha512"
	// for lines in a file, attempt to crack every sha256 hash
	Sha256 Mode = "sha256"
	// for lines in a file, attempt to crack every shadow-coded line, including algo and salt
	Shadow Mode = "shadow"
)

type AmogusConfig struct {
	LengthStart int64  `yaml:"length_start"`
	LengthEnd   int64  `yaml:"length_end"`
	Characters  string `yaml:"characters"`
	Mode        Mode   `yaml:"mode"`
	Slaves      int64  `yaml:"slaves"`
}

func (conf *AmogusConfig) Base() int {
	return len(conf.Characters)
}

func GetConfig(filename string) (*AmogusConfig, error) {
	conf, err := readConf(filename)

	if err != nil {
		return nil, err
	}

	err = validateConf(conf)

	if err != nil {
		return nil, err
	}

	return conf, err
}

func readConf(filename string) (*AmogusConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &AmogusConfig{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}

func validateConf(conf *AmogusConfig) error {
	if !containsMode([]Mode{Sha512, Sha256, Shadow}, conf.Mode) {
		return fmt.Errorf("unsupported cracking mode '%s'", conf.Mode)
	}

	if conf.LengthStart < 1 || conf.LengthStart > conf.LengthEnd {
		return fmt.Errorf("unsupported password lenghts: start %d end %d", conf.LengthStart, conf.LengthEnd)
	}

	if conf.Slaves < 1 {
		return fmt.Errorf("can't have %d (<1) slaves", conf.Slaves)
	}

	return nil
}

func containsMode(s []Mode, str Mode) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
