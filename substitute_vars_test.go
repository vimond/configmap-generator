package configmap_generator

import (
	"testing"
	"strconv"
	"fmt"
	"strings"
	"os"
	"k8s.io/kops/_vendor/github.com/docker/docker/pkg/testutil/assert"
)


func TestSubstituteVars1(t *testing.T) {
	vars := make(Variables)
	vars["key1"] = VarVal("value1")
	vars["key2"] = VarVal("{{ key1 }}")

	result := vars.SubstituteVars()
	assert.Equal(t, result["key2"], result["key1"])
	
}

//TODO Test lookup from VarVal to SecretVarVal

func TestExtractLookupString(t *testing.T) {
	testString := `'env', "AWS_ACCESS_KEY_ID"`
	result := extractLookupString(testString)
	if len(result) != 2 {
		t.Error("Expected to find 2 entries but found: " + strconv.Itoa(len(result)))
	}
	if result[0] != "env" {
		t.Error("Should equal 'env', but equals:" + result[0])
	}
}

func TestGetLookup1(t *testing.T) {
	testString := "{{ lookup('env', 'AWS_ACCESS_KEY_ID') }}"
	result, ok := getLookup(testString)
	if !ok {
		t.Error("Should have found a lookup")
	}
	expected := &AnsibleLookup{"env", []string{"AWS_ACCESS_KEY_ID"}}
	if !result.Equals(expected) {
		t.Error("lookups should be equal")
	}
}

func TestReplaceLookups1(t *testing.T) {
	testString := "{{ lookup('env', 'AWS_ACCESS_KEY_ID') }}"

	result := replaceLookups(testString)
	if result != testString {
		t.Error("strings should be equal")
	}
}

func TestReplaceLookups2(t *testing.T) {
	expected := "a-key"
	os.Setenv("AWS_ACCESS_KEY_ID", expected)
	testString := "{{ lookup('env', 'AWS_ACCESS_KEY_ID') }}"

	result := replaceLookups(testString)
	if result != expected {
		t.Error("strings should be equal")
	}
}


func TestGetPlainAnsibleVar(t *testing.T) {
	expected := "my_var"
	toTest := "{{ " + expected + " }}"
	result, ok := getPlainAnsibleVar(toTest)
	if !ok {
		t.Error("Function should have worked")
	}
	if result != expected {
		t.Error(fmt.Sprintf("Expected : %v, but got: %v",expected, result))
	}
}

func TestGetAllAnsibleVars1(t *testing.T) {
	tester := "asdf"
	result := getAllAnsibleVars(tester)
	if len(result) != 0 {
		t.Error("Function should have not returned but did:" +
			 strings.Join(result, ","))
	}
}

func TestGetAllAnsibleVars2(t *testing.T) {
	tester := "a {{ var }}  {{ lookup('type', 'val) }}"
	result := getAllAnsibleVars(tester)
	if len(result) != 2 {
		t.Error("Function should have returned 2 items but returned" +
			strings.Join(result, ","))
	}
}

func TestReplaceRegularVars1(t *testing.T) {
	vars := make(map[string]interface{})
	vars["key1"] = "value1"
	tester := "{{ key1 }}"

	result := replaceRegularVars(tester, vars)
	if result != vars["key1"] {
		t.Error(fmt.Sprintf("expected: %v, but got %v", vars["key1]"], result))
	}
}

func TestReplaceRegularVars2(t *testing.T) {
	vars := make(map[string]interface{})
	vars["key1"] = "value1"
	tester := "{{ key2 }}"

	result := replaceRegularVars(tester, vars)
	if result != tester {
		t.Error(fmt.Sprintf("expected: %v, but got %v", tester, result))
	}
}
