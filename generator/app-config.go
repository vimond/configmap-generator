package configmap_generator

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"sort"
	"github.com/pkg/errors"
	"strings"
)

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

func (a *AppConfig) CheckNameExists(name string) (bool) {
	for _, appName := range a.AppNames() {
		if (name == appName) {
			return true
		}
	}
	return false
}


func New(configPath string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &AppConfig{}, err
	}
	a := AppConfig{}
	err = yaml.Unmarshal([]byte(data), &a)
	if err != nil {
		return &AppConfig{}, err
	}
	a.SanityCheck()
	if err != nil {
		return &AppConfig{}, err
	}
	return &a, nil
}

func (a *AppConfig) SanityCheck() (error) {
	errs := make([]string, 0)
	names := a.AppNames()
	for i := 0; i < len(names) - 1; i++ {
		if names[i] == names[i+1] {
			errs = append(errs, "Duplicate application in config: " + names[i])
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
