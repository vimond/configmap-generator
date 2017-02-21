package configmap_generator

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	appConfig := New("./config/app-config.yml")
	expectedName := "vimond-docker-thumbor"
	actualName := appConfig.Applications[0].Name
	if actualName != expectedName {
		t.Errorf("Expected first name to be: '%v' but was: '%v'", expectedName, actualName)
	}
}
