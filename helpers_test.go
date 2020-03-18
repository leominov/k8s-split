package split

import "testing"

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		in  []string
		out string
	}{
		{
			in:  []string{"interspecies", "interstellar", "interstate"},
			out: "inters",
		},
		{
			in:  []string{"throne", "throne"},
			out: "throne",
		},
		{
			in:  []string{"throne", "dungeon"},
			out: "",
		},
		{
			in:  []string{"throne", "", "throne"},
			out: "",
		},
		{
			in:  []string{"cheese"},
			out: "cheese",
		},
		{
			in:  []string{""},
			out: "",
		},
		{
			in:  nil,
			out: "",
		},
	}
	for _, test := range tests {
		out := LongestCommonPrefix(test.in)
		if out != test.out {
			t.Errorf("Must be %q, but got %q", test.out, out)
		}
	}
}
