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
