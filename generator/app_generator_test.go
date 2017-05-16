package configmap_generator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConfig(t *testing.T) {
	// Load all variables
	// Load configmap_generator
	// print all variable prefix until a given depth
	// optionally exclude prefixes already covered by app config
	vars, err := SuggestConfig("../testdata/ansible1/vmp", 1)
	
	assert.Contains(t, vars, "artifactory_", "Not empty")
	assert.Contains(t, vars, "base_", "Not empty")
	assert.Contains(t, vars, "var_", "Not empty")
	
	//assert.Contains(t, vars, []string{"artifactory_","base_","var_"}, "Not empty")
	assert.Nil(t, err)
	
}

func TestMatchKey(t *testing.T) {
	assert.Equal(t, "test_", matchKey("test_match_key", 1))
	assert.Equal(t, "test_match_", matchKey("test_match_key", 2))
	assert.Equal(t, "test_match_key", matchKey("test_match_key", 3))
	assert.Equal(t, "test_match_key", matchKey("test_match_key", 4))
}
