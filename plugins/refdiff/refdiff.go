package main

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/refdiff/tasks"
	"github.com/mitchellh/mapstructure"
	"os"
)

type RefDiffOptions struct {
	RepoId string
	Pairs  []tasks.RefPair
	Tasks  []string
}

// make sure interface is implemented
var _ core.Plugin = (*RefDiff)(nil)

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry RefDiff //nolint

type RefDiff string

func (rd RefDiff) Description() string {
	return "Calculate commits diff for specified ref pairs based on `commits` and `commit_parents` tables"
}

func (rd RefDiff) Init() {
}

func (rd RefDiff) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var op RefDiffOptions
	var err error
	progress <- 0.00
	// decode options
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return fmt.Errorf("failed to parse option: %v", err)
	}
	tasksToRun := make(map[string]bool, len(op.Tasks))
	if len(op.Tasks) == 0 {
		tasksToRun = map[string]bool{
			"calculateRefDiff": true,
			"creatRefBugStats": true,
		}
	} else {
		for _, task := range op.Tasks {
			tasksToRun[task] = true
		}
	}
	// validation
	if op.RepoId == "" {
		return fmt.Errorf("repoId is required")
	}
	if tasksToRun["calculateRefDiff"] {
		progress <- 0.1
		fmt.Println("INFO >>> starting CalculateRefDiff")
		err = tasks.CalculateRefDiff(ctx, op.Pairs, op.RepoId, progress)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not calculate refdiff: %v", err).Error(),
				SubTaskName: "calculateRefDiff",
			}
		}
	}
	if tasksToRun["creatRefBugStats"] {
		progress <- 0.5
		fmt.Println("INFO >>> starting CreatRefBugStats")
		err = tasks.CreatRefBugStats(ctx, progress, op.RepoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not creat ref bug stats: %v", err).Error(),
				SubTaskName: "creatRefBugStats",
			}
		}
	}

	return nil
}

// PkgPath information lost when compiled as plugin(.so)
func (rd RefDiff) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/refdiff"
}

func (rd RefDiff) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return nil
}

// standalone mode for debugging
func main() {
	var err error

	args := os.Args[1:]
	if len(args) < 2 {
		panic(fmt.Errorf("Usage: refdiff <repo_id> <new_ref_name> <old_ref_name>"))
	}
	repoId, newRefName, oldRefName := args[0], args[1], args[2]

	err = core.RegisterPlugin("refdiff", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"repoId": repoId,
				"pairs": []map[string]string{
					{
						"NewRef": newRefName,
						"OldRef": oldRefName,
					},
				},
				"tasks": []string{
					//"calculateRefDiff",
					"creatRefBugStats",
				},
			},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
		close(progress)
	}()
	for p := range progress {
		fmt.Println(p)
	}
}
