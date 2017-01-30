package configmap_generator

import (
	"gopkg.in/yaml.v2"
	"github.com/pbthorste/avtool"
	"regexp"
)

func DecryptVaultVars(baseFolder, env, vaultPassword string) (map[string]interface{}) {
	allVars := loadSecretsAll(baseFolder, vaultPassword)
	envVars := loadSecretsEnv(baseFolder, vaultPassword, env)
	return combineMaps(allVars, envVars)
}

func loadSecretsEnv(baseFolder, vaultPassword, env string) (map[string]interface{}){
	env_folder := "/vmp/group_vars/tag_Environment_" + env + "-vpc/secrets.yml"
	return decryptVars(baseFolder, env_folder, vaultPassword)
}

func loadSecretsAll(baseFolder, vaultPassword string) (map[string]interface{}){
	return decryptVars(baseFolder, `/vmp/group_vars/all/secrets.yaml`, vaultPassword)
}

func decryptVars(baseFolder, yamlFile, vaultPassword string) (map[string]interface{}){
	result, err := avtool.Decrypt(baseFolder + yamlFile, vaultPassword)
	check(err)
	vars := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(result), &vars)
	check(err)
	return vars
}

func SubstituteVaultVars(vars, vaultVars map[string]interface{}) (map[string]interface{}) {
	checkIfVault := regexp.MustCompile(`^\{\{\s*?([\w]*?)\s*?}}`)
	for k,v := range vars {
		if checkIfVault.MatchString(v.(string)) {
			key := checkIfVault.FindAllStringSubmatch(v.(string), -1)[0][1]
			vars[k] = vaultVars[key]
		}
	}
	return vars
}