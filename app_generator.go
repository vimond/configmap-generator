package configmap_generator

import (
	"strings"
	"k8s.io/kops/_vendor/github.com/pkg/errors"
	"k8s.io/kubernetes/pkg/util/slice"
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

func extractVarPrefixes(env Variables, levels int) []string {
	
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
	
	env := loadEnv([]string{baseDir}, "skipping-vault-loading")
	
	return extractVarPrefixes(env, levels), nil
	
}
