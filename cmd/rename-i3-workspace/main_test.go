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

package main

import "testing"

func TestGetPrefix(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"workspace name", "workspace name"},
		{"10", "10:"},
		{"5:workspace name", "5:"},
	}

	for _, test := range tests {
		got := getPrefix(test.in)
		if got != test.want {
			t.Errorf("getPrefix(%q) incorrect, got %q want %q", test.in, got, test.want)
		}
	}
}
