package configmap_generator

import (
	"testing"
)

func TestLoadAllVars(t *testing.T) {
	baseFolder := "./testdata/ansible1"
	allVars := loadAll(baseFolder)
	if len(allVars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if allVars["artifactory_user"] != "developer" {
		t.Error("Expected artifactory_user to be: developer ")
	}
}

func TestLoadEnvVars(t *testing.T) {
	baseFolder := "./testdata/ansible1"
	vars := loadEnv(baseFolder, "myenv")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != "developer2" {
		t.Error("Expected artifactory_user to be: developer2")
	}
}

func TestLoadVars(t *testing.T) {
	baseFolder := "./testdata/ansible1"
	vars := LoadVars(baseFolder, "myenv")
	if len(vars) == 0 {
		t.Error("Expected to find entries but did not find any: ")
	}
	if vars["artifactory_user"] != "developer2" {
		t.Error("Expected artifactory_user to be: developer ")
	}
}