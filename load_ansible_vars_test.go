package configmap_generator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func loadedAllVars(t *testing.T, allVars map[string]interface{}  ) {
  if len(allVars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if allVars["var_only_in_all"] != "the_var" {
		t.Error("Expected artifactory_user to be: 'the_var' ")
	}
}

func TestLoadEnvVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv", "asdf")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != "developer1" {
		t.Error("Expected artifactory_user to be: developer1")
	}
	loadedAllVars(t, vars)
}

func TestLoadEnvFromFile(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv2", "asdf")
	assert.NotEmpty(t,vars, "Expected to find entries but did not find any")
}

func TestLoadEnvFromFileShortName(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv3.yml", "asdf")
	assert.NotEmpty(t,vars, "Expected to find entries but did not find any")
}

func TestLoadSecretVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv", "asdf")
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any ")
	assert.Equal(t, "{{ secret1 }}", vars["secret_reference"], "Secret not properly loaded")
}


func TestLoadVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv", "asdf")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != "developer1" {
		 t.Errorf("Expected artifactory_user to be: developer1 not %s ",  vars["artifactory_user"])
	}
}


