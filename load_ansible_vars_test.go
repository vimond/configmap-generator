package configmap_generator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestLoadAllVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	allVars,_ := loadAll(baseFolder, "asdf")
	if len(allVars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if allVars["artifactory_user"] != "developer" {
		t.Error("Expected artifactory_user to be: developer ")
	}
}

func TestLoadEnvVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := loadEnv(baseFolder, "myenv", "asdf")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != "developer1" {
		t.Error("Expected artifactory_user to be: developer1")
	}
}


func TestLoadEnvFromFile(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := loadEnv(baseFolder, "myenv2", "asdf")
	assert.NotEmpty(t,vars, "Expected to find entries but did not find any")
}

func TestLoadEnvFromFileShortName(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := loadEnv(baseFolder, "myenv3.yml", "asdf")
	assert.NotEmpty(t,vars, "Expected to find entries but did not find any")
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


