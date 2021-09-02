package main // must be main for plugin entry point

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/compound/models"
)

// plugin interface
type Compound string

func (plugin Compound) Init() {
	err := lakeModels.Db.AutoMigrate(
		&models.JiraBoardGitlabProject{},
	)
	if err != nil {
		panic(err)
	}
	text := config.V.GetString("JIRA_BOARD_GITLAB_PROJECTS")
	if text == "" {
		return
	}
	fmt.Printf("text: %v\n", text)
	rows := make([]*models.JiraBoardGitlabProject, 0)
	for _, comp := range strings.Split(text, ";") {
		if comp == "" {
			continue
		}
		fmt.Printf("comp: %v\n", comp)
		tmp := strings.Split(comp, ":")
		if len(tmp) != 2 {
			panic(fmt.Errorf("[compound] invalid config %v", text))
		}
		// board id
		boardId, _ := strconv.ParseUint(tmp[0], 10, 64)
		if boardId == 0 {
			panic(fmt.Errorf("[compound] invalid boardId %v", boardId))
		}
		// project ids
		projectIds := strings.Split(tmp[1], ",")
		if len(projectIds) == 0 {
			panic(fmt.Errorf("[compound] invalid config %v", text))
		}
		for _, pid := range projectIds {
			projectId, _ := strconv.ParseUint(pid, 10, 64)
			fmt.Printf("boardId: %v, projectId: %v\n", boardId, projectId)
			if projectId == 0 {
				panic(fmt.Errorf("[compound] invalid projectId %v", pid))
			}
			rows = append(rows, &models.JiraBoardGitlabProject{JiraBoardId: boardId, GitlabProjectId: projectId})
		}
	}
	err = lakeModels.Db.Exec("truncate table jira_board_gitlab_projects").Error
	if err != nil {
		panic(fmt.Errorf("[compound] failed to truncate jira_board_gitlab_projects"))
	}
	for _, row := range rows {
		fmt.Printf("row: %v\n", row)
		if row == nil {
			continue
		}
		err = lakeModels.Db.Create(row).Error
		if err != nil {
			panic(fmt.Errorf("[compound] failed to insert %v %v", row, err))
		}
	}
}

func (plugin Compound) Description() string {
	return "To relate jira board and gitlab projects"
}

func (plugin Compound) Execute(options map[string]interface{}, progress chan<- float32) {
	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Compound //nolint
