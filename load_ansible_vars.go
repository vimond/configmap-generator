package configmap_generator

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	"github.com/pbthorste/avtool"
	"github.com/pkg/errors"
	"fmt"
)


/**
Loads ansible variables for a given environment
 */

// Loads variables from the vimond-ansible project.
// baseFolder: the location of the vimond-ansible project
// env: the name of the environment to use
func LoadVars(baseFolder, env, vaultPassword string) (map[string]interface{}) {
	allVars, _ := loadAll(baseFolder, vaultPassword)
	envVars := loadEnv(baseFolder, env, vaultPassword)
	return combineMaps(allVars, envVars)
}

func createPathAlternatives(baseFolder, env string) []string {
	return []string{
		fmt.Sprintf("%s/tag_Environment_%s-vpc", baseFolder, env),
		fmt.Sprintf("%v/%v", baseFolder, env),
	}
}
type loadVarFn func(path, vaultPassword string)(map[string]interface{}, error)

func createLoadVarsFunctions() ([]loadVarFn)   {
	return []loadVarFn{ loadVarsInFile,loadVarsInFolder }
}

func loadEnv(baseFolder, env, vaultPassword string) (map[string]interface{}) {
	
	var envVars map[string]interface{}
	e := make([]error, 4)
	
	for _, path := range createPathAlternatives(baseFolder, env) {
		for _, loadVarsFn := range createLoadVarsFunctions() {
			vars, err := loadVarsFn(path,vaultPassword)
			if err != nil {
				e = append(e,err)
			}
			if(vars != nil) {
				return vars
			}
		}
	}
	if(envVars == nil) {
		checkErrs(e)
	}
	return envVars
	
}

func loadAll(baseFolder, vaultPassword string) (map[string]interface{}, error) {
	return loadVarsInFolder(baseFolder + `/all`, vaultPassword)
}

var loadVarsInFolder = func (folder, vaultPassword string) (map[string]interface{}, error) {
	all_vars := make(map[string]interface{})
	files, err := ioutil.ReadDir(folder)
	if (err != nil) {
		return nil, err
	}
	for _, v := range files {
		current := folder + "/" + v.Name()
		vars, err := loadVarsInFile(current, vaultPassword)
		if err != nil {
			check(err)
		}
		all_vars = combineMaps(all_vars, vars)
	}
	
	return all_vars, err
}

var loadVarsInFile = func(current, vaultPassword string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(current)
	if err != nil {
		return nil, err
	}
	var vars map[string]interface{}
	if checkIfVault(data) {
		vars, err = decryptVault(current, vaultPassword)
	} else {
		vars, err = loadPlain(data)
	}
	return vars, err
}

func checkIfVault(fileContents []byte) (bool) {
	contents := string(fileContents)
	return strings.HasPrefix(contents, "$ANSIBLE_VAULT")
}

func loadPlain(fileContents []byte) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(fileContents), &vars)
	
	return vars, err
}

func decryptVault(vaultFile, vaultPassword string) (map[string]interface{}, error) {
	result, err := avtool.Decrypt(vaultFile, vaultPassword)
	if err != nil {
		return nil,errors.Wrapf(err, "Problems decrypting '%v', passsphrase correct?",vaultFile)
	}
	vars := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(result), &vars)
	return vars, err
}

func combineMaps(maps ...map[string]interface{}) (map[string]interface{}) {
	combined := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			combined[k] = v
		}
	}
	return combined
}
