package configmap_generator

import (
	"testing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//Test to mess around with some marshalling, not testing anything real
func TestName(t *testing.T) {

	bytes, _ := ioutil.ReadFile("./testdata/ansible1/vmp/group_vars/myenv3.yml")
	vars := make(map[string]interface{})
	yaml.Unmarshal(bytes, &vars)
	for key, value := range vars {
		t.Logf("%v: %v\n", key, value)
	}
	bytes, _ = yaml.Marshal(vars)
	newYaml := string(bytes)
	t.Log(newYaml)

	vars2 := make(map[string]interface{})
	yaml.Unmarshal(bytes, &vars2)
	t.Logf("%v",vars2["string_with_json"])


}