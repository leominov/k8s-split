package split

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"reflect"
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

func TestGetNameAndKindAndPartof(t *testing.T) {
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
		_, _, _, err := GetNameAndKindAndPartof(test.val)
		if err == nil {
			t.Error("Must be an error, but got nil")
		}
	}
	successTest := map[string]interface{}{
		"kind": "kind",
		"metadata": map[string]interface{}{
			"name": "name",
			"labels": map[string]interface{}{
				"app.kubernetes.io/part-of": "bar",
			},
		},
	}
	kind, name, partof, err := GetNameAndKindAndPartof(successTest)
	if err != nil {
		t.Error(err)
	}
	if kind != "kind" {
		t.Errorf("Must be kind, but got %s", kind)
	}
	if name != "name" {
		t.Errorf("Must be name, but got %s", name)
	}
	if partof != "bar" {
		t.Errorf("Must be bar, but got %s", partof)
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

func TestProcess_Tag(t *testing.T) {
	Tag = true
	defer func() {
		Tag = false
	}()

	dir4, err := ioutil.TempDir("", "k8s-split-prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir4)
	err = Process("test_data/correct_multi_prefix.yaml", dir4)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(path.Join(dir4, "bar"))
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(path.Join(dir4, "bar", "application.Pod.yaml"))
	if err != nil {
		t.Error(err)
	}
}

func TestFindUniqueLabelValues(t *testing.T) {
	t.Run("Object with app.kubernetes.io/part-of label", func(t *testing.T) {
		testObj := []map[string]interface{}{
			{
				"kind": "kind",
				"metadata": map[string]interface{}{
					"name": "name",
					"labels": map[string]interface{}{
						"app.kubernetes.io/part-of": "foo",
					},
				},
			},
			{
				"kind": "kind",
				"metadata": map[string]interface{}{
					"name": "name",
					"labels": map[string]interface{}{
						"app.kubernetes.io/part-of": "bar",
					},
				},
			},
		}

		want := []string{"foo", "bar"}
		got, _ := FindUniqueLabelValues(testObj)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	},
	)

	t.Run("Object with duplicate app.kubernetes.io/part-of label", func(t *testing.T) {
		testObj := []map[string]interface{}{
			{
				"kind": "kind",
				"metadata": map[string]interface{}{
					"name": "name",
					"labels": map[string]interface{}{
						"app.kubernetes.io/part-of": "foo",
					},
				},
			},
			{
				"kind": "kind",
				"metadata": map[string]interface{}{
					"name": "name2",
					"labels": map[string]interface{}{
						"app.kubernetes.io/part-of": "foo",
					},
				},
			},
		}

		want := []string{"foo"}
		got, _ := FindUniqueLabelValues(testObj)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	},
	)
}

func TestPreparePrefixedDirectory(t *testing.T) {

	testObj := []map[string]interface{}{
		{
			"kind": "Pod",
			"metadata": map[string]interface{}{
				"name": "application",
				"labels": map[string]interface{}{
					"app.kubernetes.io/part-of": "foo",
				},
			},
		},
		{
			"kind": "Service",
			"metadata": map[string]interface{}{
				"name": "application",
				"labels": map[string]interface{}{
					"app.kubernetes.io/part-of": "bar",
				},
			},
		},
		{
			"kind": "Job",
			"metadata": map[string]interface{}{
				"name": "application",
				"labels": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}

	TestCases := []struct {
		Description  string
		Prefix       bool
		Tag          bool
		val          []map[string]interface{}
		ExpectedDirs []string
	}{
		{
			"Split by prefix",
			true,
			false,
			testObj,
			[]string{"application"},
		},
		{
			"Split by tag",
			false,
			true,
			testObj,
			[]string{"foo", "bar"},
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.Description, func(t *testing.T) {
			Prefix = tc.Prefix
			Tag = tc.Tag
			defer func() {
				Prefix = false
				Tag = false
			}()

			maindir, err := ioutil.TempDir("", "output")
			if err != nil {
				t.Fatal(err)
			}

			defer os.RemoveAll(maindir)

			_, err = preparePrefixedDirectory(tc.val, maindir)
			if err != nil {
				t.Fatal(err)
			}

			for _, dir := range tc.ExpectedDirs {
				_, err = os.Stat(path.Join(maindir, dir))
				if err != nil {
					t.Error(err)
				}
			}
		})

	}
}
