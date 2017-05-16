package configmap_generator

import (
	"testing"
)


func TestLoadTetsConfig(t *testing.T) {
	appConfig := New("../testdata/test-app-config.yml")
	expectedName := "test"
	actualName := appConfig.Applications[0].Name
	if actualName != expectedName {
		t.Errorf("Expected first name to be: '%v' but was: '%v'", expectedName, actualName)
	}
}

