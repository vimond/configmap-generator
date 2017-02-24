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
func LoadVars(baseFolder, env, vaultPassword string) Variables {
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

func loadEnv(searchPaths []string, vaultPassword string) Variables {

	envVars := make(Variables)
	//var envVars map[string]interface{}
	
	e := make([]error, len(searchPaths))

	//Use walker and visitor instead
	for _, path := range searchPaths {
		filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			fileVars, err := loadVarsInFile(path, vaultPassword)
			if err != nil {
				e = append(e, err)
			}
			//envVars = combineMaps(envVars, fileVars)
			envVars.AddAll(fileVars)
			
			return nil
		})
		
	}
	if envVars == nil {
		checkErrs(e)
	}
	return envVars
}
 
func loadVarsInFile (current, vaultPassword string) (Variables, error) {
	data, err := ioutil.ReadFile(current)
	if err != nil {
		return nil, err
	}
	//var vars Variables
	if checkIfVault(data) {
		//varVals, e := decryptVault2(current, vaultPassword)
		return decryptVault(current, vaultPassword)
		//err = e
		//vars = makeSecret(varVals)
		
	} else {
		//vars, err = loadPlain2(data)
		return  loadPlain(data)
	}
	//return vars, err
}


type Variables map[string]Secret

type PlainVariables map[string]VarVal
type SecVariables map[string]SecretVarVal

type VarVal string
type SecretVarVal string

type Secret interface {
	IsSecret() bool
	String() string
}

func (SecretVarVal) IsSecret() bool{
	return true
}
func (VarVal) IsSecret() bool{
	return false
}
func (s SecretVarVal) String() string{
	return string(s)
}

func (v VarVal) String() string{
	return string(v)
}
func (vars Variables) AddAll(other Variables)  {
		for k, v := range other {
			vars[k] = v
		}
}

//func (item T) compare(other Comparable) int {
//
//}
	

func checkIfVault(fileContents []byte) bool {
	contents := string(fileContents)
	return strings.HasPrefix(contents, "$ANSIBLE_VAULT")
}

//func loadPlain(fileContents []byte) (map[string]interface{}, error) {
//	vars := make(map[string]interface{})
//	err := yaml.Unmarshal([]byte(fileContents), &vars)
//	return vars, err
//}

func loadPlain(fileContents []byte) (Variables, error) {
	var plainVars PlainVariables
	var vars Variables = Variables{}
	err := yaml.Unmarshal([]byte(fileContents), &plainVars)
	for key, value := range plainVars {
		vars[key] = value
	}
	return vars, err
}

func decryptVault(vaultFile, vaultPassword string) (Variables, error) {
	result, err := avtool.Decrypt(vaultFile, vaultPassword)
	if err != nil {
		return nil, errors.Wrapf(err, "Problems decrypting '%v', passsphrase correct?", vaultFile)
	}
	var secVars SecVariables
	var vars Variables = Variables{}
	
	err = yaml.Unmarshal([]byte(result), &secVars)
	for key, value := range secVars {
		vars[key] = value
	}
	return vars, err
}

//func decryptVault(vaultFile, vaultPassword string) (map[string]interface{}, error) {
//	result, err := avtool.Decrypt(vaultFile, vaultPassword)
//	if err != nil {
//		return nil, errors.Wrapf(err, "Problems decrypting '%v', passsphrase correct?", vaultFile)
//	}
//	vars := make(map[string]interface{})
//	err = yaml.Unmarshal([]byte(result), &vars)
//	return vars, err
//}
//
//func combineMaps(maps ...map[string]interface{}) map[string]interface{} {
//	combined := make(map[string]interface{})
//	for _, m := range maps {
//		for k, v := range m {
//			combined[k] = v
//		}
//	}
//	return combined
//}
