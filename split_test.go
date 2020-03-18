package split

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestMultiByEntries(t *testing.T) {
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
		res, err := MultiByEntries(bytes.NewReader([]byte(test.input)))
		if err != nil {
			t.Error(err)
		}
		if len(res) != test.count {
			t.Errorf("Must be %d, but got %d", test.count, len(res))
		}
	}
	_, err := MultiByEntries(bytes.NewReader([]byte(`
	`)))
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
	dir, err := ioutil.TempDir("", "k8s-split")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tests := []string{
		"test_data/correct_single.yaml",
		"test_data/correct_multi.yaml",
		"test_data/correct_list.yaml",
		"test_data/correct_single_with_items.yaml",
	}
	for _, test := range tests {
		if err := Process(test, dir); err != nil {
			t.Error(err)
		}
	}
	Quiet = true
	defer func() {
		Quiet = false
	}()
	err = os.Chmod(dir, 0444)
	if err != nil {
		t.Error(err)
	}
	if err := Process("test_data/correct_single.yaml", dir); err == nil {
		t.Error("Must be an error, but got nil")
	}
	tests = []string{
		"test_data/incorrect_not_found.yaml",
		"test_data/incorrect_1.yaml",
		"test_data/incorrect_2.yaml",
		"test_data/incorrect_3.yaml",
		"test_data/incorrect_list.yaml",
	}
	for _, test := range tests {
		if err := Process(test, dir); err == nil {
			t.Error("Must be an error, but got nil")
		}
	}

	err = Process("-", dir)
	// Empty is not an error
	if err != nil {
		t.Error(err)
	}
}

func TestProcess_Prefix(t *testing.T) {
	Prefix = true
	defer func() {
		Prefix = false
	}()

	dir1, err := ioutil.TempDir("", "k8s-split-prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir1)
	err = Process("test_data/correct_multi_prefix.yaml", dir1)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(path.Join(dir1, "application"))
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(path.Join(dir1, "application", "application.Pod.yaml"))
	if err != nil {
		t.Error(err)
	}

	dir2, err := ioutil.TempDir("", "k8s-split-prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir2)
	err = Process("test_data/correct_multi_prefix_empty.yaml", dir2)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(path.Join(dir2, "application.Pod.yaml"))
	if err != nil {
		t.Error(err)
	}

	dir3, err := ioutil.TempDir("", "k8s-split-prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir3)
	_, err = os.Create(path.Join(dir3, "application"))
	if err != nil {
		t.Fatal(err)
	}
	err = Process("test_data/correct_multi_prefix.yaml", dir3)
	if err == nil {
		t.Error("Must be an error, but got nil")
	}
}
