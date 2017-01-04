package configmap_generator

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)


/**
Loads ansible variables for a given environment
 */

// Loads variables from the vimond-ansible project.
// baseFolder: the location of the vimond-ansible project
// env: the name of the environment to use
func LoadVars(baseFolder, env string) (map[string]interface{}) {
	allVars := loadAll(baseFolder)
	envVars := loadEnv(baseFolder, env)
	return combineMaps(allVars, envVars)
}

func loadEnv(baseFolder, env string) (map[string]interface{}){
	env_folder := "/vmp/group_vars/tag_Environment_" + env + "-vpc/general_vars.yml"
	return loadVars(baseFolder, env_folder)
}

func loadAll(baseFolder string) (map[string]interface{}){
	return loadVars(baseFolder, `/vmp/group_vars/all/vars.yaml`)
}

func loadVars(baseFolder, yamlFile string) (map[string]interface{}){
	data, err := ioutil.ReadFile(baseFolder + yamlFile)
	check(err)
	all_vars := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(data), &all_vars)
	check(err)
	return all_vars
}

func combineMaps(allVars, envVars map[string]interface{}) (map[string]interface{}){
	combined := make(map[string]interface{})
	for k,v := range allVars {
		combined[k] = v
	}
	for k,v := range envVars {
		combined[k] = v
	}
	return combined
}
