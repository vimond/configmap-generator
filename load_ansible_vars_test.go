package configmap_generator

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"

)

func loadedAllVars(t *testing.T, allVars Variables) {
	if len(allVars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if allVars["var_only_in_all"] != VarVal("the_var") {
		t.Error("Expected artifactory_user to be: 'the_var' ")
	}
}

func TestLoadEnvVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv", "asdf")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	
	//val := -
	actual := vars["artifactory_user"].String()
	assert.Equal(t, actual, string(SecretVarVal("developer1")))
	//if vars["artifactory_user"] != VarVal("devloper1") {
	//	t.Error("Expected artifactory_user to be: developer1")
	//}
	loadedAllVars(t, vars)
}

func TestLoadVarVals(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/ansible1/vmp/group_vars/all/vars.yaml")
	assert.NoError(t, err)
	vars, err := loadPlain(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, vars)
	
	t.Log("%v", vars)
	for _, value := range vars {
		assert.False(t,value.IsSecret())
	}
}


func TestLoadSecretVarVals(t *testing.T) {
	vars, err := decryptVault("./testdata/ansible1/vmp/group_vars/tag_Environment_myenv-vpc/secrets.yaml", "asdf")
	assert.NoError(t, err)
	assert.NotEmpty(t, vars)
	
	for _, value := range vars {
		assert.True(t,value.IsSecret())
	}
}


func TestLoadEnvFromFile(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv2", "asdf")
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any")
}

func TestLoadEnvFromFileShortName(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv3.yml", "asdf")
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any")
}

func TestLoadVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv", "asdf")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != VarVal("developer1") {
		t.Errorf("Expected artifactory_user to be: developer1 not %s ", vars["artifactory_user"])
	}
}


