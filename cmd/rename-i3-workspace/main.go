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

// Show an i3-input prompt and rename the current i3 workspace using the given
// name.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

func main() {
	name, err := focusedWorkspaceName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("name: %q\n", name)

	prefix := getPrefix(name)
	fmt.Printf("prefix: %s\n", prefix)

	i3cmd := fmt.Sprintf("rename workspace to \"%s%%s\"", prefix)
	i3prompt := fmt.Sprintf("workspace name> %s", prefix)

	cmd := exec.Command("i3-input", "-F", i3cmd, "-P", i3prompt)
	fmt.Printf("running: %s\n", cmd)

	if cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func focusedWorkspaceName() (string, error) {
	tree, err := i3.GetTree()
	if err != nil {
		return "", err
	}

	ws := tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Type == i3.WorkspaceNode
	})

	if ws == nil {
		return "", errors.New("could not find active workspace")
	}
	return ws.Name, nil
}

func getPrefix(name string) string {
	if idx := strings.Index(name, ":"); idx >= 0 {
		return name[:idx+1]
	}
	if _, err := strconv.Atoi(name); err == nil {
		return name + ":"
	}
	return name
}
