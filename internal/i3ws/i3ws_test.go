// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		if err := got.Set(test.in); err != nil {
			t.Errorf("Set(%q) failed: %v\n", test.in, err)
			continue
		}
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Set(%q) incorrect, diff (-got, +want):\n%s", test.in, diff)
		}
	}
}
