package main

import (
	"context"

	"github.com/merico-dev/lake/plugins/core"
)

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
	/* TODO: adopt new interface
	var op tasks.RefdiffOptions
	var err error
	progress <- 0.00
	// decode options
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return fmt.Errorf("failed to parse option: %v", err)
	}

	logger := helper.NewDefaultTaskLogger(nil, "refdiff")

	tasksToRun := make(map[string]bool, len(op.Tasks))
	if len(op.Tasks) == 0 {
		tasksToRun = map[string]bool{
			"calculateCommitsDiff":  true,
			"calculateIssuesDiff":   true,
			"calculatePrCherryPick": true,
		}
	} else {
		for _, task := range op.Tasks {
			tasksToRun[task] = true
		}
	}

	taskData := &tasks.RefdiffTaskData{
		Options: &op,
	}

	taskCtx := helper.NewDefaultTaskContext("refdiff", ctx, logger, taskData, tasksToRun)
	newTasks := []struct {
		name       string
		entryPoint core.SubTaskEntryPoint
	}{
		{name: "calculateCommitsDiff", entryPoint: tasks.CalculateCommitsDiff},
		{name: "calculateIssuesDiff", entryPoint: tasks.CalculateIssuesDiff},
		{name: "calculatePrCherryPick", entryPoint: tasks.CalculatePrCherryPick},
	}
	for _, t := range newTasks {
		c, err := taskCtx.SubTaskContext(t.name)
		if err != nil {
			return err
		}
		if c != nil {
			err = t.entryPoint(c)
			if err != nil {
				return &errors.SubTaskError{
					SubTaskName: t.name,
					Message:     err.Error(),
				}
			}
		}
	}
	*/
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
	/* TODO: adopt new method
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
					//"calculateCommitsDiff",
					//"calculateIssuesDiff",
					"calculatePrCherryPick",
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
	*/
}
