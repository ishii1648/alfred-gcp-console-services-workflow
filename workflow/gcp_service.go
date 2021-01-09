package workflow

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GcpService struct {
	Id          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	ShortName   string       `yaml:"short_name"`
	Description string       `yaml:"description"`
	Url         string       `yaml:"url"`
	SubServices []GcpService `yaml:"sub_services"`
}

func ParseConsoleServicesYml(ymlPath string) []GcpService {
	gcpServices := []GcpService{}
	yamlFile, err := ioutil.ReadFile(ymlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &gcpServices)
	if err != nil {
		log.Fatal(err)
	}
	return gcpServices
}

func (g *GcpService) GetName() string {
	if g.ShortName != "" {
		return g.ShortName + " – " + g.Name
	}
	return g.Name
}
