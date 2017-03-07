package configmap_generator

import (
	"testing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

//Test to mess around with some marshalling, not testing anything real
func TestName(t *testing.T) {
	bytes, _ := ioutil.ReadFile("./testdata/ansible1/vmp/group_vars/myenv3.yml")
	vars := make(map[string]interface{})
	yaml.Unmarshal(bytes, &vars)
	for key, value := range vars {
		fmt.Printf("%v: %v\n", key, value)
	}
	bytes, _ = yaml.Marshal(vars)
	fmt.Print(string(bytes))

}