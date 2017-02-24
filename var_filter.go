package configmap_generator

import (
	"strings"
	"fmt"
	"log"
)
/**
Filters variables for application
 */



func FilterVariables(appConfig *AppConfig, ansibleVars Variables, appName string) (Variables) {
	prefixVars := filterByPrefix(appConfig, ansibleVars, appName)
	return prefixVars
}

func filterByPrefix(appConfig *AppConfig, ansibleVars Variables, appName string) (Variables) {
	vars := make(Variables)
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

func getVarsByPrefix(ansibleVars Variables, prefix string) (Variables) {
	vars := make(Variables)
	for k,v := range ansibleVars {
		if strings.HasPrefix(k, prefix) {
			vars[k] = v
		}
	}
	return vars
}