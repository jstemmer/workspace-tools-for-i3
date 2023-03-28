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
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

var (
	inputCmd = flag.String("input-cmd", "i3-input -F {i3cmd} -P {prompt}{prefix}", "Runs the given command to get user input. The following special strings will be replaced with their actual values: {prompt}, {prefix}, {i3cmd}. Falls back to i3-input if empty.")
	sendToI3 = flag.Bool("x", false, "Send the rename workspace command to i3 after requesting input. Use this when the input-cmd does not rename the workspace.")
)

func main() {
	flag.Parse()
	if len(strings.TrimSpace(*inputCmd)) == 0 {
		fmt.Fprintf(os.Stderr, "-input-cmd is a required flag\n")
		os.Exit(1)
	}

	name, err := focusedWorkspaceName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("name: %q\n", name)

	prefix := getPrefix(name)
	fmt.Printf("prefix: %s\n", prefix)

	prompt := fmt.Sprintf("workspace name> ")
	i3cmd := `rename workspace to "` + prefix + `%s"`

	input, err := runInputCommand(*inputCmd, prompt, prefix, i3cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *sendToI3 {
		if strings.TrimSpace(input) == "" {
			return
		}

		newName := escapeName(input)
		if _, err := i3.RunCommand(fmt.Sprintf(`rename workspace to "%s"`, newName)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
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

func runInputCommand(inputCmd, prompt, prefix, i3cmd string) (string, error) {
	cmd := inputCmd
	var args []string
	if c, a, ok := strings.Cut(inputCmd, " "); ok {
		cmd = c
		args = strings.Split(a, " ")
	}

	for i := range args {
		args[i] = strings.ReplaceAll(args[i], "{prompt}", prompt)
		args[i] = strings.ReplaceAll(args[i], "{prefix}", prefix)
		args[i] = strings.ReplaceAll(args[i], "{i3cmd}", i3cmd)
	}
	fmt.Printf("running: %s %s\n", cmd, strings.Join(args, " "))
	out, err := exec.Command(cmd, args...).Output()
	return strings.TrimSpace(string(out)), err
}

func escapeName(input string) string {
	name := input
	name = strings.ReplaceAll(name, `\`, `\\`)
	name = strings.ReplaceAll(name, `"`, `\"`)
	return name
}
