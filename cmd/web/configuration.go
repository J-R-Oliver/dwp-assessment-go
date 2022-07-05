package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
)

type peopleConfiguration struct {
	BaseURL  string `yaml:"base-url"`
	Distance int    `yaml:"default-distance"`
}

type city struct {
	Name      string
	Latitude  string `yaml:"lat"`
	Longitude string `yaml:"lon"`
}

type configuration struct {
	Port                string
	PeopleConfiguration peopleConfiguration
	Cities              map[string]city
}

func loadConfiguration() configuration {
	configFile, err := ioutil.ReadFile("configuration.yaml")
	if err != nil {
		log.Fatal(err)
	}

	configFile = []byte(os.ExpandEnv(string(configFile)))

	r := regexp.MustCompile(`([\w\d:/\-.]*):-([\w\d:/\-.]*)`) // ToDo - This should be at the top of the file

	submatch := r.FindAllStringSubmatch(string(configFile), -1)
	_ = submatch

	fixedString := string(configFile)

	for _, match := range submatch {
		if match[1] != "" {
			fixedString = strings.Replace(fixedString, match[0], match[1], 1)
		} else {
			fixedString = strings.Replace(fixedString, match[0], match[2], 1)
		}
	}

	var config configuration

	if err = yaml.Unmarshal([]byte(fixedString), &config); err != nil {
		log.Fatal(err)
	}

	return config
}
