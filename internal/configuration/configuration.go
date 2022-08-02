package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
	"gopkg.in/yaml.v3"
)

type peopleConfiguration struct {
	BaseURL  string `yaml:"base-url"`
	Distance int    `yaml:"default-distance"`
}

type City struct {
	Latitude  string `yaml:"lat"`
	Longitude string `yaml:"lon"`
}

type Configuration struct {
	Port                string              `yaml:"port"`
	LoggingLevel        logging.Level       `yaml:"logging-level"`
	PeopleConfiguration peopleConfiguration `yaml:"people"`
	Cities              map[string]City
}

func LoadConfiguration(filename string) (Configuration, error) {
	var config Configuration

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("loadConfiguration error: unable to read configuration file - %s: %w", filename, err)
	}

	configFile := string(b)

	configFile = os.ExpandEnv(configFile)

	configFile = replaceSubstitutions(configFile)

	if err = yaml.Unmarshal([]byte(configFile), &config); err != nil {
		return config, fmt.Errorf("loadConfiguration error: unable to parse configuration file - %s: %w", filename, err)
	}

	return config, nil
}

var r = regexp.MustCompile(`([\w\d:/\-.]*):-([\w\d:/\-.]*)`)

func replaceSubstitutions(configFile string) string {
	submatch := r.FindAllStringSubmatch(configFile, -1)

	for _, match := range submatch {
		if match[1] != "" {
			configFile = strings.Replace(configFile, match[0], match[1], 1)
		} else {
			configFile = strings.Replace(configFile, match[0], match[2], 1)
		}
	}

	return configFile
}
