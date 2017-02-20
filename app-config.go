package configmap_generator

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
	"sort"
	"os"
)

func check(err error) {
	if(err != nil) {
		log.Printf("error: %v", err)
		panic(err)
	}
}

func checkErrs(errs []error) {
	if errs != nil && len(errs) > 0{
		for _, error := range errs {
			log.Printf("error: %v", error)
		}
		os.Exit(1)
	}
}
type Application struct {
	Name string
	Prefixes []string
}
type AppConfig struct {
	Applications []Application
}

// List of all the applications names
func (a *AppConfig) AppNames() []string {
	names := make([]string, len(a.Applications))
	for key,entry := range a.Applications {
		names[key] = entry.Name
	}
	sort.Strings(names)
	return names
}

func LoadConfig() (*AppConfig) {
	data, err := ioutil.ReadFile(`./config/app-config.yml`)
	check(err)
	a := AppConfig{}
	err = yaml.Unmarshal([]byte(data), &a)
	check(err)
	a.SanityCheck()
	return &a
}

func (a *AppConfig) SanityCheck() {
	names := a.AppNames()
	for i := 0; i < len(names) - 1; i++ {
		if names[i] == names[i+1] {
			error := "Duplicate application in config: " + names[i]
			panic(error)
		}
	}
}
