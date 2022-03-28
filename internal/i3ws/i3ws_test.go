package i3ws

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWorkspacesFlagString(t *testing.T) {
	tests := []struct {
		flag WorkspacesFlag
		want string
	}{
		{nil, ""},
		{WorkspacesFlag{1}, "1"},
		{WorkspacesFlag{2, 3, 4, 5, 10}, "2,3,4,5,10"},
	}

	for _, test := range tests {
		if got := test.flag.String(); got != test.want {
			t.Errorf("WorkspaceFlag.String() incorrect, got %q want %q", got, test.want)
		}
	}
}

func TestWorkspacesFlagSet(t *testing.T) {
	tests := []struct {
		in   string
		want WorkspacesFlag
	}{
		{"", WorkspacesFlag([]int{})},
		{"3", WorkspacesFlag{3}},
		{"10,1, 7, 2 ,5 ", WorkspacesFlag{1, 2, 5, 7, 10}},
	}

	for _, test := range tests {
		got := NewWorkspacesFlag()
		got.Set(test.in)
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Set(%q) incorrect, diff (-got, +want):\n%s", test.in, diff)
		}
	}
}
