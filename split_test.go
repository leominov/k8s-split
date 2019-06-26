package split

import (
	"os"
	"testing"
)

func TestSplitByEntries(t *testing.T) {
	tests := []struct {
		input string
		count int
	}{
		{
			input: `document: 1`,
			count: 1,
		},
		{
			input: `---
document: 1
---
document: 2
`,
			count: 2,
		},
	}
	for _, test := range tests {
		res, err := SplitByEntries([]byte(test.input))
		if err != nil {
			t.Error(err)
		}
		if len(res) != test.count {
			t.Errorf("Must be %d, but got %d", test.count, len(res))
		}
	}
	_, err := SplitByEntries([]byte(`
	`))
	if err == nil {
		t.Error("Must be an error, but got nil")
	}
}

func TestGetNameAndKind(t *testing.T) {
	tests := []struct {
		val        interface{}
		name, kind string
	}{
		{
			val: "foobar",
		},
		{
			val: map[string]interface{}{
				"kind": "kind",
			},
		},
		{
			val: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name": "name",
				},
			},
		},
	}
	for _, test := range tests {
		_, _, err := GetNameAndKind(test.val)
		if err == nil {
			t.Error("Must be an error, but got nil")
		}
	}
	successTest := map[string]interface{}{
		"kind": "kind",
		"metadata": map[string]interface{}{
			"name": "name",
		},
	}
	kind, name, err := GetNameAndKind(successTest)
	if err != nil {
		t.Error(err)
	}
	if kind != "kind" {
		t.Errorf("Must be kind, but got %s", kind)
	}
	if name != "name" {
		t.Errorf("Must be name, but got %s", name)
	}
}

func TestProcess(t *testing.T) {
	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll("tmp")
	tests := []string{
		"test_data/correct_single.yaml",
		"test_data/correct_multi.yaml",
	}
	for _, test := range tests {
		if err := Process(test, "tmp"); err != nil {
			t.Error(err)
		}
	}
	tests = []string{
		"test_data/incorrect_not_found.yaml",
		"test_data/incurrect_1.yaml",
		"test_data/incurrect_2.yaml",
		"test_data/incurrect_3.yaml",
	}
	for _, test := range tests {
		if err := Process(test, "tmp"); err == nil {
			t.Error("Must be an error, but got nil")
		}
	}
}
