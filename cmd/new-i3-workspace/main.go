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

// Change to a new i3 workspace
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jstemmer/workspace-tools-for-i3/internal/i3ws"

	"go.i3wm.org/i3/v4"
)

func main() {
	max := flag.Int("max", 20, "maximum workspace `number` to consider")
	reserved := i3ws.NewWorkspacesFlag()
	flag.Var(&reserved, "reserved", "list of reserved workspace numbers")
	name := flag.String("name", "", "name to use for new workspace")
	flag.Parse()

	num, err := nextAvailableWorkspace(reserved, *max)
	if err != nil {
		exitf("could not find next available workspace number: %v\n", err)
	}

	if err := changeWorkspace(num, *name); err != nil {
		exitf("could not create new workspace: %v\n", err)
	}
}

// nextAvailableWorkspace returns the next available unused workspace number
// that is not reserved. Workspace numbers up to and including max are
// considered. If no workspace number is available, an error is returned.
func nextAvailableWorkspace(reserved []int, max int) (int, error) {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return 0, err
	}

	active := make(map[int]bool)
	for _, ws := range workspaces {
		n, err := parseWorkspaceName(ws.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing workspace %q: %v\n", ws.Name, err)
			continue
		}
		active[n] = true
	}

	// treat reserved workspaces as active
	for _, ws := range reserved {
		active[ws] = true
	}

	for i := 1; i <= max; i++ {
		if !active[i] {
			return i, nil
		}
	}
	return 0, errors.New("no more workspaces available")
}

// changeWorkspace switches to a new workspace with the given number and
// optional name.
func changeWorkspace(number int, name string) error {
	var n string
	if name != "" {
		n = fmt.Sprintf("%d:%s", number, strings.ReplaceAll(name, `"`, `\"`))
	} else {
		n = strconv.Itoa(number)
	}
	cmd := fmt.Sprintf(`workspace number "%s"`, n)
	_, err := i3.RunCommand(cmd)
	return err
}

// parseWorkspaceName returns the number prefix of a workspace name.
func parseWorkspaceName(name string) (int, error) {
	number := name
	if idx := strings.Index(name, ":"); idx >= 0 {
		number = name[:idx]
	}
	return strconv.Atoi(number)
}

func exitf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
