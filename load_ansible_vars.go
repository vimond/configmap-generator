package configmap_generator

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	"github.com/pbthorste/avtool"
)


/**
Loads ansible variables for a given environment
 */

// Loads variables from the vimond-ansible project.
// baseFolder: the location of the vimond-ansible project
// env: the name of the environment to use
func LoadVars(baseFolder, env, vaultPassword string) (map[string]interface{}) {
	allVars := loadAll(baseFolder, vaultPassword)
	envVars := loadEnv(baseFolder, env, vaultPassword)
	return combineMaps(allVars, envVars)
}

func loadEnv(baseFolder, env, vaultPassword string) (map[string]interface{}){
	return loadVarsInFolder(baseFolder +  "/tag_Environment_" + env + "-vpc", vaultPassword)
}

func loadAll(baseFolder, vaultPassword string) (map[string]interface{}){
	return loadVarsInFolder(baseFolder +  `/all`, vaultPassword)
}


func loadVarsInFolder(folder, vaultPassword string) (map[string]interface{}){
	all_vars := make(map[string]interface{})
	files, err := ioutil.ReadDir(folder)
	check(err)
	for _,v := range files {
		current := folder + "/" +  v.Name()
		data, err := ioutil.ReadFile(current)
		check(err)
		var vars map[string]interface{}
		if checkIfVault(data) {
			vars = decryptVault(current, vaultPassword)
		} else {
			vars = loadPlain(data)
		}

		all_vars = combineMaps(all_vars, vars)
	}

	return all_vars
}

func checkIfVault(fileContents []byte) (bool) {
	contents := string(fileContents)
	return strings.HasPrefix(contents, "$ANSIBLE_VAULT")
}

func loadPlain(fileContents []byte) (map[string]interface{}){
	vars := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(fileContents), &vars)
	check(err)
	return vars
}

func decryptVault(vaultFile, vaultPassword string) (map[string]interface{}){
	result, err := avtool.Decrypt(vaultFile, vaultPassword)
	check(err)
	vars := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(result), &vars)
	check(err)
	return vars
}

func combineMaps(maps ...map[string]interface{}) (map[string]interface{}) {
	combined := make(map[string]interface{})
	for _,m := range maps {
		for k,v := range m {
			combined[k] = v
		}
	}
	return combined
}
