package main // must be main for plugin entry point

import (
	"context"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/api"
)

var _ core.Plugin = (*Jira)(nil)

// plugin interface
type Jira string

/*
func (plugin Jira) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	err := db.AutoMigrate(
		&models.JiraProject{},
		&models.JiraUser{},
		&models.JiraIssue{},
		&models.JiraBoard{},
		&models.JiraBoardIssue{},
		&models.JiraChangelog{},
		&models.JiraChangelogItem{},
		&models.JiraRemotelink{},
		&models.JiraIssueCommit{},
		&models.JiraSource{},
		&models.JiraIssueTypeMapping{},
		&models.JiraIssueStatusMapping{},
		&models.JiraSprint{},
		&models.JiraBoardSprint{},
		&models.JiraSprintIssue{},
		&models.JiraWorklog{},
	)
	if err != nil {
		panic(err)
	}
}
*/
func (plugin Jira) Init() {

}

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	/* TODO: adopt new interface
	// process options
	var op tasks.JiraOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}
	if op.SourceId == 0 {
		return fmt.Errorf("sourceId is invalid")
	}
	source := &models.JiraSource{}
	err = lakeModels.Db.Find(source, op.SourceId).Error
	if err != nil {
		return err
	}

	var since time.Time
	if op.Since != "" {
		since, err = time.Parse("2006-01-02T15:04:05Z", op.Since)
		if err != nil {
			return fmt.Errorf("invalid value for `since`: %w", err)
		}
	}

	tasksToRun := map[string]bool{
		"collectProjects":     true,
		"extractProjects":     true,
		"collectBoard":        true,
		"extractBoard":        true,
		"collectIssues":       true,
		"extractIssues":       true,
		"collectChangelogs":   true,
		"extractChangelogs":   true,
		"collectRemotelinks":  true,
		"extractRemotelinks":  true,
		"collectSprints":      true,
		"extractSprints":      true,
		"convertBoard":        true,
		"convertIssues":       true,
		"collectWorklogs":     true,
		"extractWorklogs":     true,
		"convertWorklogs":     true,
		"convertChangelogs":   true,
		"convertSprints":      true,
		"convertIssueCommits": true,
	}
	if len(op.Tasks) > 0 {
		// set all to false
		for task := range tasksToRun {
			tasksToRun[task] = false
		}
		// set those specified tasks to true
		for _, task := range op.Tasks {
			tasksToRun[task] = true
		}
	}

	rateLimit := source.RateLimit
	if rateLimit <= 0 {
		rateLimit = 50
	}
	scheduler, err := utils.NewWorkerScheduler(rateLimit, rateLimit, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	// prepare contextual variables
	logger := helper.NewDefaultTaskLogger(nil, "jira")
	jiraApiClient := tasks.NewJiraApiClient(
		source.Endpoint,
		source.BasicAuthEncoded,
		source.Proxy,
		scheduler,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create jira api client: %w", err)
	}

	taskData := &tasks.JiraTaskData{
		Options:   &op,
		ApiClient: &jiraApiClient.ApiClient,
		Source:    source,
	}
	if !since.IsZero() {
		taskData.Since = &since
	}
	info, code, err := jiraApiClient.GetJiraServerInfo()
	if err != nil || code != http.StatusOK {
		return fmt.Errorf("fail to get server info")
	}

	if info.DeploymentType == models.DeploymentServer {
		if versions := info.VersionNumbers; len(versions) == 3 && versions[0] == 8 {
			tasksToRun["collectChangelogs"] = false
			tasksToRun["extractChangelogs"] = false
		} else {
			return fmt.Errorf("not support version:%s", info.Version)
		}
	}
	taskCtx := helper.NewDefaultTaskContext("jira", ctx, logger, taskData, tasksToRun)
	newTasks := []struct {
		name       string
		entryPoint core.SubTaskEntryPoint
	}{
		{name: "collectProjects", entryPoint: tasks.CollectProjects},
		{name: "extractProjects", entryPoint: tasks.ExtractProjects},

		{name: "collectBoard", entryPoint: tasks.CollectBoard},
		{name: "extractBoard", entryPoint: tasks.ExtractBoard},

		{name: "collectIssues", entryPoint: tasks.CollectIssues},
		{name: "extractIssues", entryPoint: tasks.ExtractIssues},

		{name: "collectChangelogs", entryPoint: tasks.CollectChangelogs},
		{name: "extractChangelogs", entryPoint: tasks.ExtractChangelogs},

		{name: "collectWorklogs", entryPoint: tasks.CollectWorklogs},
		{name: "extractWorklogs", entryPoint: tasks.ExtractWorklogs},

		{name: "collectRemotelinks", entryPoint: tasks.CollectRemotelinks},
		{name: "extractRemotelinks", entryPoint: tasks.ExtractRemotelinks},

		{name: "collectSprints", entryPoint: tasks.CollectSprints},
		{name: "extractSprints", entryPoint: tasks.ExtractSprints},

		{name: "convertBoard", entryPoint: tasks.ConvertBoard},

		{name: "convertIssues", entryPoint: tasks.ConvertIssues},

		{name: "convertWorklogs", entryPoint: tasks.ConvertWorklogs},

		{name: "convertChangelogs", entryPoint: tasks.ConvertChangelogs},

		{name: "convertSprints", entryPoint: tasks.ConvertSprints},

		{name: "convertIssueCommits", entryPoint: tasks.ConvertIssueCommits},
	}
	for _, t := range newTasks {
		c, err := taskCtx.SubTaskContext(t.name)
		if err != nil {
			return err
		}
		fmt.Println(c)
		if c != nil {
			logger.Info("enter")
			err = t.entryPoint(c)
			if err != nil {
				return &errors.SubTaskError{
					SubTaskName: t.name,
					Message:     err.Error(),
				}
			}
		}
	}

<<<<<<< HEAD
=======
	// run tasks
	var collector tasks.Collector
	info, code, err := jiraApiClient.GetJiraServerInfo()
	if err != nil || code != http.StatusOK {
		return fmt.Errorf("fail to get server info")
	}
	if info.DeploymentType == models.DeploymentCloud {
		collector = tasks.NewCloud()
	}
	if info.DeploymentType == models.DeploymentServer {
		if versions := info.VersionNumbers; len(versions) == 3 && versions[0] == 8 {
			collector = tasks.NewServerVersion8(lakeModels.Db, jiraApiClient)
		}
	}
	if collector == nil {
		return fmt.Errorf("Jira server %s is not supported", info.Version)
	}

	logger.Info("start plugin execution")
	if tasksToRun["collectProjects"] {
		err := collector.CollectProjects(jiraApiClient, op.SourceId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectProjects",
				Message:     err.Error(),
			}
		}
	}
	if tasksToRun["collectUsers"] {
		err := collector.CollectUsers(jiraApiClient, op.SourceId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectUsers",
				Message:     err.Error(),
			}
		}
	}
	if tasksToRun["collectBoard"] {
		err := collector.CollectBoard(jiraApiClient, source, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectBoard",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.01
	if tasksToRun["collectIssues"] {
		err = collector.CollectIssues(jiraApiClient, source, boardId, since, source.RateLimit, ctx)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectIssues",
				Message:     err.Error(),
			}
		}
	}

	progress <- 0.1
	if tasksToRun["collectChangelogs"] {
		err = collector.CollectChangelogs(jiraApiClient, source, boardId, rateLimit, ctx)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectChangelogs",
				Message:     err.Error(),
			}
		}
	}
	if tasksToRun["collectRemotelinks"] {
		err = collector.CollectRemoteLinks(jiraApiClient, source, boardId, rateLimit, ctx)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectRemotelinks",
				Message:     err.Error(),
			}
		}
	}
	// TODO: Remove this subtask, we can collect user while we are collecting issues
	if tasksToRun["collectUsers"] {
		err := collector.CollectUsers(jiraApiClient, op.SourceId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectUsers",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.2
	if tasksToRun["enrichIssues"] {
		err = tasks.EnrichIssues(source, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "enrichIssues",
				Message:     err.Error(),
			}
		}
	}
	if tasksToRun["enrichRemotelinks"] {
		err = tasks.EnrichRemotelinks(source, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "enrichRemotelinks",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.3
	if tasksToRun["collectSprints"] {
		err = collector.CollectSprint(jiraApiClient, source, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectSprints",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.4
	if tasksToRun["convertBoard"] {
		err := tasks.ConvertBoard(op.SourceId, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertBoard",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.5
	if tasksToRun["convertUsers"] {
		err := tasks.ConvertUsers(op.SourceId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertUsers",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.7
	if tasksToRun["convertWorklogs"] {
		err = tasks.ConvertWorklog(op.SourceId, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertWorklogs",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.8
	if tasksToRun["convertChangelogs"] {
		err = tasks.ConvertChangelogs(op.SourceId, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertChangelogs",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.85
	if tasksToRun["convertSprints"] {
		err = tasks.ConvertSprint(op.SourceId, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertSprints",
				Message:     err.Error(),
			}
		}
	}
	progress <- 0.9
	if tasksToRun["convertIssueCommits"] {
		err = tasks.ConvertIssueCommits(op.SourceId, boardId)
		if err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertIssueCommits",
				Message:     err.Error(),
			}
		}
	}
	progress <- 1.0
	logger.Info("end plugin execution")
	*/
>>>>>>> feat: decouple global variables
	return nil
}

func (plugin Jira) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jira"
}

func (plugin Jira) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"echo": {
			"POST": func(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
				return &core.ApiResourceOutput{Body: input.Body}, nil
			},
		},
		"sources": {
			"POST": api.PostSources,
			"GET":  api.ListSources,
		},
		"sources/:sourceId": {
			"PUT":    api.PutSource,
			"DELETE": api.DeleteSource,
			"GET":    api.GetSource,
		},
		"sources/:sourceId/epics": {
			"GET": api.GetEpicsBySourceId,
		},
		"sources/:sourceId/granularities": {
			"GET": api.GetGranularitiesBySourceId,
		},
		"sources/:sourceId/boards": {
			"GET": api.GetBoardsBySourceId,
		},
		"sources/:sourceId/type-mappings": {
			"POST": api.PostIssueTypeMappings,
			"GET":  api.ListIssueTypeMappings,
		},
		"sources/:sourceId/type-mappings/:userType": {
			"PUT":    api.PutIssueTypeMapping,
			"DELETE": api.DeleteIssueTypeMapping,
		},
		"sources/:sourceId/type-mappings/:userType/status-mappings": {
			"POST": api.PostIssueStatusMappings,
			"GET":  api.ListIssueStatusMappings,
		},
		"sources/:sourceId/type-mappings/:userType/status-mappings/:userStatus": {
			"PUT":    api.PutIssueStatusMapping,
			"DELETE": api.DeleteIssueStatusMapping,
		},
		"sources/:sourceId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira //nolint

// standalone mode for debugging
func main() {
	/* TODO: adopt new method
	args := os.Args[1:]
	if len(args) < 2 {
		panic(fmt.Errorf("Usage: jira <source_id> <board_id>"))
	}
	sourceId, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		panic(fmt.Errorf("error paring source_id: %w", err))
	}
	boardId, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("error paring board_id: %w", err))
	}

	err = core.RegisterPlugin("jira", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"sourceId": sourceId,
				"boardId":  boardId,
				"tasks": []string{
					//"collectBoard",
					//"collectProjects",
					//"collectIssues",
					"collectApiIssues",
					"extractApiIssues",
					//"collectChangelogs",
					//"collectApiChangelogs",
					//"collectRemotelinks",
					//"enrichIssues",
					//"enrichRemotelinks",
					//"collectSprints",
					//"collectUsers",
					//"convertBoard",
					"convertIssues",
					//"convertWorklogs",
					//"convertChangelogs",
					//"convertUsers",
					//"convertSprints",
					//"convertIssueCommits",
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
