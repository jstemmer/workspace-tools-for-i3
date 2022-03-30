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

// Package i3ws contains common code shared by the various
// workspace-tools-for-i3 commands.
package i3ws

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type WorkspacesFlag []int

func NewWorkspacesFlag() WorkspacesFlag {
	return make(WorkspacesFlag, 0)
}

func (f *WorkspacesFlag) String() string {
	var nums []string
	if f != nil {
		for _, ws := range *f {
			nums = append(nums, strconv.Itoa(ws))
		}
	}
	return strings.Join(nums, ",")
}

func (f *WorkspacesFlag) Set(value string) error {
	workspaces := make(map[int]struct{})
	for _, field := range strings.Split(value, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil {
			return fmt.Errorf("invalid workspace number: %w", err)
		}
		workspaces[n] = struct{}{}
	}
	for n := range workspaces {
		*f = append(*f, n)
	}
	sort.Slice(*f, func(i, j int) bool {
		return (*f)[i] < (*f)[j]
	})
	return nil
}
