package configmap_generator

import (
	"reflect"
	"strings"
	"regexp"
	"os"
)

/*
Some basic variable substitution for the ansible map
 */

func SubstituteVars(vars map[string]interface{}) (map[string]interface{}) {
	for k,v := range vars {
		if reflect.TypeOf(v) != reflect.TypeOf("string") {
			// only supports strings at this time
			// could support maps later
			continue
		}
		value := v.(string)
		if !strings.Contains(value, "{{") {
			// no ansible var found
			continue
		}

		if isPure, ok := getPlainAnsibleVar(value); ok {
			if replace, ok := vars[isPure]; ok {
				vars[k] = replace
			} else {
				os.Stderr.WriteString("Var not found: " + value + "\n")
			}
		} else {
			replaceStringVars(value, vars)
		}
	}
	return vars

}


func getPlainAnsibleVar(value string) (string, bool) {
	pureVal := regexp.MustCompile(`^\{\{\s*?([\w]*?)\s*?}}$`)
	if !pureVal.MatchString(value) {
		return "", false
	} else {
		return pureVal.FindAllStringSubmatch(value, -1)[0][1], true
	}
}

// Gets all ansible vars of the type {{ something }} contained in a string.
// Returns an empty list if nothing is found
func getAllAnsibleVars(value string) ([]string) {
	ansibleVarFinder := regexp.MustCompile(`\{\{.*?}}`)
	return ansibleVarFinder.FindAllString(value, -1)
}

func replaceStringVars(value string, vars map[string]interface{}) (string) {
	updated := value
	updated = replaceRegularVars(updated, vars)
	updated = replaceLookups(updated)
	return updated
}

/*
Replace simple ansible vars in the string
 */
func replaceRegularVars(value string, vars map[string]interface{}) (string) {
	updated := value
	ansibleVars := getAllAnsibleVars(value)
	for _,item := range ansibleVars {
		if plainVar, ok := getPlainAnsibleVar(item); ok {
			if replace, ok := vars[plainVar]; ok {
				updated = strings.Replace(updated, item, replace.(string), 1)
			} else {
				os.Stderr.WriteString("Var not found: " + item + "\n")
			}
		}
	}
	return updated
}

type AnsibleLookup struct {
	Type string
	Data []string
}

func (local *AnsibleLookup) Equals(other *AnsibleLookup) bool {
	if other == nil {
		return false
	}
	if local.Type != other.Type {
		return false
	}
	if len(local.Data) != len(other.Data) {
		return false
	}
	for i := range local.Data {
		if local.Data[i] != other.Data[i] {
			return false
		}
	}
	return true
}


func getLookup(variable string) (*AnsibleLookup, bool) {
	lookupFinder := regexp.MustCompile(`\{\{\s*?lookup\((.*)?\)\s*?}}$`)
	result := lookupFinder.FindAllStringSubmatch(variable, -1)
	if len(result) != 0 {
		data := result[0][1]
		parts := extractLookupString(data)
		return &AnsibleLookup{parts[0], parts[1:]}, true
	}
	return nil, false
}

// Extracts the contents of a lookup. Supports single and double quotes
// returns items of the lookup in order. The first item is the type of lookup
func extractLookupString(lookup string) ([]string) {
	lookupFinder := regexp.MustCompile(`'[^\\']*?'|"[^\\"]*?"|\S+`)
	items := lookupFinder.FindAllString(lookup, -1)
	indexes := make([]int, 0)
	for i,v := range items {
		if strings.HasPrefix(v, "'") || strings.HasPrefix(v, `"`) {
			items[i] = v[1:len(v)-1]
			indexes = append(indexes, i)
		}
	}
	result := make([]string, len(indexes))
	for i,v := range indexes {
		result[i] = items[v]
	}
	return result
}

func replaceLookups(value string) (string) {
	updated := value
	ansibleVars := getAllAnsibleVars(value)
	for _,item := range ansibleVars {
		if lookup, ok := getLookup(item); ok {
			if lookup.Type != "env" {
				os.Stderr.WriteString("Unsupported lookup: " + lookup.Type + "\n")
			} else {
				if envVar, ok := getEnvLookup(lookup.Data[0]); ok {
					updated = strings.Replace(updated, item, envVar, 1)
				}
			}
		}
	}
	return updated
}

func getEnvLookup(envVar string) (string, bool) {
	result := os.Getenv(envVar)
	if result == "" {
		os.Stderr.WriteString("Environment var not set: " + envVar + "\n")
		return "", false
	} else {
		return result, true
	}
}