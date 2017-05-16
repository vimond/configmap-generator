package configmap_generator

import (
	"testing"
)


func TestLoadTetsConfig(t *testing.T) {
	appConfig, err := New("../testdata/test-app-config.yml")
	if err != nil {
		t.Error(err)
	}
	expectedName := "test"
	actualName := appConfig.Applications[0].Name
	if actualName != expectedName {
		t.Errorf("Expected first name to be: '%v' but was: '%v'", expectedName, actualName)
	}
}

