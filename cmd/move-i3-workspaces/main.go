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
	"sort"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

type Options struct {
	dryrun     bool
	all        bool
	except     []int64
	workspaces []int64
	output     string
}

func main() {
	dryrun := flag.Bool("n", false, "dry run; don't move anything, just print what would be moved")
	allWorkspaces := flag.Bool("all", false, "select all workspaces")
	exceptWorkspaces := flag.String("except", "", "exclude workspaces, for use in combination with --all")
	flag.Parse()

	except, err := parseWorkspaces(*exceptWorkspaces)
	if err != nil {
		exitf("invalid value for --except: %v\n", err)
	}

	var chosenWorkspaces string
	if flag.NArg() > 1 {
		chosenWorkspaces = flag.Arg(0)
	}

	output := flag.Arg(flag.NArg() - 1)

	workspaces, err := parseWorkspaces(chosenWorkspaces)
	if err != nil {
		exitf("invalid value for --workspaces: %v\n", err)
	}

	options := Options{
		dryrun:     *dryrun,
		all:        *allWorkspaces,
		except:     except,
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
		except := make(map[int64]bool)
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
			if _, ok := except[w.Num]; ok {
				continue // skip workspaces in the exception list
			}
			workspaces = append(workspaces, w.Num)
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

func moveWorkspace(num int64, output string) error {
	cmd := fmt.Sprintf("[workspace=%d] move workspace to output %s", num, output)
	_, err := i3.RunCommand(cmd)
	return err
}

func parseWorkspaces(in string) ([]int64, error) {
	if in == "" {
		return nil, nil
	}

	workspaces := make(map[int64]struct{})
	for _, field := range strings.Split(in, ",") {
		n, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid workspace number: %w", err)
		}
		workspaces[n] = struct{}{}
	}
	var ws []int64
	for n := range workspaces {
		ws = append(ws, n)
	}
	sort.Slice(ws, func(i, j int) bool {
		return ws[i] < ws[j]
	})
	return ws, nil
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
