package configmap_generator

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	"github.com/pbthorste/avtool"
	"fmt"
)


/**
Loads ansible variables for a given environment
 */

// Loads variables from the vimond-ansible project.
// baseFolder: the location of the vimond-ansible project
// env: the name of the environment to use
func LoadVars(baseFolder, env, vaultPassword string) (map[string]interface{}) {
	allVars,_ := loadAll(baseFolder, vaultPassword)
	envVars := loadEnv(baseFolder, env, vaultPassword)
	return combineMaps(allVars, envVars)
}

func loadEnv(baseFolder, env, vaultPassword string) (map[string]interface{}){
	dirFolder := fmt.Sprintf("%s/tag_Environment_%s-vpc", baseFolder, env)
	e := make([]error,3)
	var vars map[string]interface{}
	
	vars, e[0] = loadVarsInFolder(dirFolder, vaultPassword)
	
	if e[0] != nil {
		vars, e[1] = loadVarsInFile(dirFolder, vaultPassword)
		if e[1] != nil {
			vars, e[2] = loadVarsInFile(baseFolder+"/"+env, vaultPassword)
			if e[2] != nil {
				checkErrs(e)
			}
		}
	}
	return vars
}

func loadAll(baseFolder, vaultPassword string) (map[string]interface{}, error){
	return loadVarsInFolder(baseFolder +  `/all`, vaultPassword)
}


func loadVarsInFolder(folder, vaultPassword string) (map[string]interface{}, error){
	all_vars := make(map[string]interface{})
	files, err := ioutil.ReadDir(folder)
	if (err != nil) {
		return nil, err
	}
	
	//check(err)
	for _,v := range files {
		current := folder + "/" +  v.Name()
		vars, err := loadVarsInFile(current, vaultPassword)
		if err != nil {
			check(err)
		}
		all_vars = combineMaps(all_vars, vars)
	}

	return all_vars, err
}

func loadVarsInFile(current, vaultPassword string) (map[string]interface{}, error){
	data, err := ioutil.ReadFile(current)
	if err != nil {
		return nil, err
	}
	var vars map[string]interface{}
	if checkIfVault(data) {
		vars = decryptVault(current, vaultPassword)
	} else {
		vars = loadPlain(data)
	}
	return vars, err
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
