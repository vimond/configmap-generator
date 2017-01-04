package configmap_generator

import (
	"testing"
	"fmt"
)

func TestLoadConfig(t *testing.T) {
	appConfig := LoadConfig()
	fmt.Println(appConfig)
	expectedName := "vimond-docker-thumbor"
	actualName := appConfig.Applications[0].Name
	if actualName != expectedName {
		t.Errorf("Expected first name to be: " + expectedName + " but was: " +
		 actualName)
	}
}
