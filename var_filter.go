package configmap_generator

import (
	"strings"
	"fmt"
	"log"
)
/**
Filters variables for application
 */



func FilterVariables(appConfig *AppConfig, ansibleVars map[string]interface{}, appName string) (map[string]interface{}) {
	prefixVars := filterByPrefix(appConfig, ansibleVars, appName)
	return prefixVars
}

func filterByPrefix(appConfig *AppConfig, ansibleVars map[string]interface{}, appName string) (map[string]interface{}) {
	vars := make(map[string]interface{})
	app := getApplication(appConfig, appName)
	for _,v := range app.Prefixes {
		prefixMap := getVarsByPrefix(ansibleVars, v)
		for key, val := range prefixMap {
			vars[key] = val
		}
	}
	return vars
}

func getApplication(appConfig *AppConfig, appName string) (*Application) {
	for _,entry := range appConfig.Applications {
		if entry.Name == appName {
			return &entry
		}
	}
	log.Fatalf(fmt.Sprintf("Error - could not find application: %v", appName))
	panic(fmt.Sprintf("Error - could not find application: %v", appName))
}

func getVarsByPrefix(ansibleVars map[string]interface{}, prefix string) (map[string]interface{}) {
	vars := make(map[string]interface{})
	for k,v := range ansibleVars {
		if strings.HasPrefix(k, prefix) {
			vars[k] = v
		}
	}
	return vars
}