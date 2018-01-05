package configmap_generator

import (
	"fmt"
	"errors"
	"strings"
	"gopkg.in/yaml.v2"
)

type GeneratorConfig struct {
	AppName			string
	Environment		string
	GroupVars		string
	VaultPassword	string
	AppConfig		*AppConfig
}


func (config *GeneratorConfig) GenerateConfigMap() (string, error){

	allVars, err := LoadVars(config.GroupVars, config.Environment, config.VaultPassword)
	if err != nil {
		return "", err
	}

	var result string
	if config.AppName != "all" {
		result, err = getConfigMap(config.AppName, allVars, config.AppConfig)
		if err != nil {
			return "", err
		}
	} else {
		result, err = getAllConfigMaps(allVars, config.AppConfig)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

func (config *GeneratorConfig) GenerateConfigMapAsMap() (map[string]interface{}, error) {
	allVars, err := LoadVars(config.GroupVars, config.Environment, config.VaultPassword)
	if err != nil {
		return map[string]interface{}{}, err
	}
	allVars["service_name"] = config.AppName
	allVars = SubstituteVars(allVars)
	return FilterVariables(config.AppConfig, allVars, config.AppName), nil
}

func getConfigMap(name string, allVars map[string]interface{}, appConfig *AppConfig) (string, error) {
	allVars["service_name"] = name
	allVars = SubstituteVars(allVars)
	vars := FilterVariables(appConfig, allVars, name)
	vars2,_ := yaml.Marshal(vars)
	app := ConfigMapData{
		AppName: name,
		Data: string(vars2[:]),
	}
	return Generate(app)
}

func getAllConfigMaps(allVars map[string]interface{}, appConfig *AppConfig) (string, error) {
	var err error
	errs := make([]string, 0)
	configMaps := make([]string, len(appConfig.Applications))

	for i, v := range appConfig.Applications {
		configMaps[i], err = getConfigMap(v.Name, allVars, appConfig)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%v", err))
		}
	}
	if len(errs) > 0 {
		return "", errors.New(strings.Join(errs, "\n"))
	}
	return strings.Join(configMaps, "\n"), nil
}