package configmap_generator

import (
	"testing"
)


func TestLoadAllVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	allVars := loadAll(baseFolder, "asdf")
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
	if vars["artifactory_user"] != "developer2" {
		t.Error("Expected artifactory_user to be: developer2")
	}
}

func TestLoadVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := LoadVars(baseFolder, "myenv", "asdf")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != "developer2" {
		t.Error("Expected artifactory_user to be: developer ")
	}
}


