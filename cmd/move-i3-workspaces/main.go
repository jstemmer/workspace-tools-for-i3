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

// Quickly move i3 workspaces to a specified output
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jstemmer/workspace-tools-for-i3/internal/i3ws"

	"go.i3wm.org/i3/v4"
)

type Options struct {
	dryrun     bool
	all        bool
	except     []int
	workspaces []int
	output     string
}

func main() {
	dryrun := flag.Bool("n", false, "dry run; don't move anything, just print what would be moved")
	allWorkspaces := flag.Bool("all", false, "select all workspaces")
	except := i3ws.NewWorkspacesFlag()
	flag.Var(&except, "except", "exclude `workspaces`, for use in combination with --all")
	flag.Parse()

	var chosenWorkspaces string
	if flag.NArg() > 1 {
		chosenWorkspaces = flag.Arg(0)
	}

	output := flag.Arg(flag.NArg() - 1)

	workspaces := i3ws.NewWorkspacesFlag()
	if err := workspaces.Set(chosenWorkspaces); err != nil {
		exitf("invalid list of workspaces: %v\n", err)
	}

	options := Options{
		dryrun:     *dryrun,
		all:        *allWorkspaces,
		except:     []int(except),
		workspaces: workspaces,
		output:     output,
	}
	if err := options.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
	if err := moveWorkspaces(options); err != nil {
		exitf("error: %v\n", err)
	}
}

func (o Options) Validate() error {
	if o.all && len(o.workspaces) > 0 {
		return fmt.Errorf("can't specify both --all and a list of workspaces")
	}
	if !o.all && len(o.workspaces) == 0 {
		return fmt.Errorf("no workspaces specified")
	}
	if len(o.except) > 0 && len(o.workspaces) > 0 {
		return fmt.Errorf("can't specify both --except and a list of workspaces")
	}
	if o.output == "" {
		return fmt.Errorf("no output specified")
	}
	return nil
}

func moveWorkspaces(options Options) error {
	var found bool
	outputs, err := listOutputs()
	if err != nil {
		return fmt.Errorf("error listing outputs: %w", err)
	}
	for _, output := range outputs {
		if output == options.output {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("output %v does not exist", options.output)
	}

	workspaces := options.workspaces
	if options.all {
		except := make(map[int]bool)
		for _, e := range options.except {
			except[e] = true
		}

		ws, err := i3.GetWorkspaces()
		if err != nil {
			return fmt.Errorf("could not get i3 workspaces: %w", err)
		}

		for _, w := range ws {
			if w.Output == options.output {
				continue // skip workspaces that are already on the right output
			}
			num := int(w.Num)
			if _, ok := except[num]; ok {
				continue // skip workspaces in the exception list
			}
			workspaces = append(workspaces, num)
		}
	}

	if options.dryrun {
		fmt.Printf("Dry run:\n")
	}
	if len(workspaces) == 0 {
		fmt.Printf("Nothing to do!\n")
	}
	for _, workspace := range workspaces {
		fmt.Printf("Moving workspace %d to output %s\n", workspace, options.output)
		if !options.dryrun {
			moveWorkspace(workspace, options.output)
		}
	}
	return nil
}

func moveWorkspace(num int, output string) error {
	cmd := fmt.Sprintf("[workspace=%d] move workspace to output %s", num, output)
	_, err := i3.RunCommand(cmd)
	return err
}

func exitf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func listOutputs() ([]string, error) {
	i3outputs, err := i3.GetOutputs()
	if err != nil {
		return nil, err
	}

	var outputs []string
	for _, output := range i3outputs {
		if !output.Active {
			continue
		}
		outputs = append(outputs, output.Name)
	}
	return outputs, nil
}
