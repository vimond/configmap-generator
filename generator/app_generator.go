package configmap_generator

import (
	"strings"
	"k8s.io/kubernetes/pkg/util/slice"
	"errors"
)

const sep = "_"

func matchKey(key string, levels int) string {
	
	splits := strings.Split(key, sep)
	idx := levels
	postfix := sep
	if(idx >= len(splits)) {
		idx = len(splits)
		postfix = ""
	}
	prefix := strings.Join(splits[:idx], sep)
	return  prefix + postfix
}

func extractVarPrefixes(env map[string]interface{}, levels int) []string {
	
	keyMap := make(map[string]bool)
	
	for key := range env {
		match := matchKey(key, levels)
		keyMap[match] = true
	}
	
	keys := make([]string, 0, len(keyMap))
	for k := range keyMap {
		keys = append(keys, k)
	}
	
	return slice.SortStrings(keys)
}

func SuggestConfig(baseDir string, levels int) ([]string, error) {
	if (levels <= 0 ) {
		return nil, errors.New("Levels must be gt 0")
	}
	
	env, err := loadEnv([]string{baseDir}, "skipping-vault-loading")
	if err != nil {
		return []string{}, err
	}
	
	return extractVarPrefixes(env, levels), nil
	
}
