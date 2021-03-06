package configmap_generator

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pbthorste/avtool"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"os"
)

/**
Loads ansible variables for a given environment
*/

// Loads variables from the vimond-ansible project.
// baseFolder: the location of the vimond-ansible project
// env: the name of the environment to use
func LoadVars(baseFolder, env, vaultPassword string) (map[string]interface{}, error) {
	return loadEnv(createPathAlternatives(baseFolder, env), vaultPassword)
}

func createPathAlternatives(baseFolder, env string) []string {
	return []string{
		fmt.Sprintf("%s/all", baseFolder),
		fmt.Sprintf("%s/all.yml", baseFolder),
		fmt.Sprintf("%s/all.yaml", baseFolder),
		fmt.Sprintf("%s/tag_Environment_%s-vpc", baseFolder, env),
		fmt.Sprintf("%v/%v", baseFolder, env),
		fmt.Sprintf("%v/%v.yml", baseFolder, env),
		fmt.Sprintf("%v/%v.yaml", baseFolder, env),
	}
}


func loadEnv(searchPaths []string, vaultPassword string) (map[string]interface{}, error) {

	var envVars map[string]interface{}
	e := make([]string, len(searchPaths))

	//Use walker and visitor instead
	for _, path := range searchPaths {
		filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			fileVars, err := loadVarsInFile(path, vaultPassword)
			if err != nil {
				e = append(e, fmt.Sprintf("%v", err))
			}
			envVars = combineMaps(envVars, fileVars)
			return nil
		})
		
	}
	if envVars == nil {
		return map[string]interface{}{}, errors.New(strings.Join(e, "\n"))
	}
	return envVars, nil

}
 
func loadVarsInFile (current, vaultPassword string) (map[string]interface{}, error) {
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

func checkIfVault(fileContents []byte) bool {
	contents := string(fileContents)
	return strings.HasPrefix(contents, "$ANSIBLE_VAULT")
}

func loadPlain(fileContents []byte) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(fileContents), &vars)

	return vars, err
}

func decryptVault(vaultFile, vaultPassword string) (map[string]interface{}, error) {
	result, err := avtool.DecryptFile(vaultFile, vaultPassword)
	if err != nil {
		return nil, errors.Wrapf(err, "Problems decrypting '%v', passsphrase correct?", vaultFile)
	}
	vars := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(result), &vars)
	return vars, err
}

func combineMaps(maps ...map[string]interface{}) map[string]interface{} {
	combined := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			combined[k] = v
		}
	}
	return combined
}
