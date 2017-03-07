package configmap_generator

import (
	"testing"
	"strconv"
	"fmt"
	"strings"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestSubstituteVars1(t *testing.T) {
	vars := make(map[string]interface{})
	vars["key1"] = "value1"
	vars["key2"] = "{{ key1 }}"
	
	result := SubstituteVars(vars)
	if result["key2"] != result["key1"] {
		t.Error("Expected vars to be equal")
	}
}

func TestLookupWhenSourcedFromFiles(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := SubstituteVars(LoadVars(baseFolder, "myenv", "asdf"))
	
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any ")
	assert.Equal(t, "secret stuff 1", vars["secret_reference"], "Secret not properly loaded")
	
}

func TestLookupPrecedenceWhenSourcedFromFiles_myenv(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := SubstituteVars(LoadVars(baseFolder, "myenv", "asdf"))
	
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any ")
	assert.Equal(t, "myenv_lookup_secret", vars["general_repeated_lookup"], "Secret not properly loaded")
	
}

func TestLookupPrecedenceWhenSourcedFromFiles_myenv2(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := SubstituteVars(LoadVars(baseFolder, "myenv2", "asdf"))
	
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any ")
	assert.Equal(t, "all_lookup_secret", vars["general_repeated_lookup"], "Secret not properly loaded")
	
}

func TestNonStringsConverted(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := SubstituteVars(LoadVars(baseFolder, "myenv3", "asdf"))

	assert.IsType(t,"string", vars["string_boolean"], "%v", vars["string_boolean"])
	assert.IsType(t,"string", vars["string_arr"], "%v", vars["string_arr"])
	assert.IsType(t,"string", vars["string_dict"], "%v", vars["string_dict"])
	assert.IsType(t,"string", vars["string_int"], "%v", vars["string_int"])

}

func TestLookupSimpleInlineStringVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	loadVars := LoadVars(baseFolder, "myenv3", "asdf")
	substituted := replaceStringVars(loadVars["string_lookup"].(string), loadVars)

	assert.Contains(t, substituted, "developer3")
}

func TestLookupComplexInlineStringVars(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	loadVars := LoadVars(baseFolder, "myenv3", "asdf")
	substituted := replaceStringVars(loadVars["multiline_complex_preformatted"].(string), loadVars)
	
	assert.Contains(t, substituted, "developer3")
}


func TestLookupComplexMultilinePeformatted_myenv3(t *testing.T) {
	baseFolder := "./testdata/ansible1/vmp/group_vars"
	vars := SubstituteVars(LoadVars(baseFolder, "myenv3", "asdf"))
	
	assert.NotEmpty(t, vars, "Expected to find entries but did not find any ")
	assert.Contains(t, vars["multiline_complex_preformatted"], "developer3")
}


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
		t.Error(fmt.Sprintf("Expected : %v, but got: %v", expected, result))
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
