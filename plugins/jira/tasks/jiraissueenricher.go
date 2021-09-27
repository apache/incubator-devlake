package tasks

import (
	"fmt"
	"strings"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jira/models"
)

var storyPointCoefficient float64

func init() {
	storyPointCoefficient = config.V.GetFloat64("JIRA_ISSUE_STORYPOINT_COEFFICIENT")
}

func EnrichIssues(boardId uint64) error {
	jiraIssue := &models.JiraIssue{}

	// select all issues belongs to the board
	cursor, err := lakeModels.Db.Model(jiraIssue).
		Select("jira_issues.*").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_issues.id").
		Where("jira_board_issues.board_id = ?", boardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		if jiraIssue.ResolutionDate.Valid {
			jiraIssue.LeadTime = uint(jiraIssue.ResolutionDate.Time.Unix()-jiraIssue.Created.Unix()) / 60
		}
		jiraIssue.StdStoryPoint = uint(jiraIssue.StoryPoint * storyPointCoefficient)
		jiraIssue.StdType = getStdType(jiraIssue)
		jiraIssue.StdStatus = getStdStatus(jiraIssue)
		err = lakeModels.Db.Save(jiraIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}

var typeMapping map[string]string

func getStdType(jiraIssue *models.JiraIssue) string {
	if typeMapping == nil {
		typeMapping = getStringMapFromConfig("JIRA_ISSUE_TYPE_MAPPING")
	}
	if stdType, ok := typeMapping[jiraIssue.Type]; ok {
		return stdType
	}
	return jiraIssue.Type
}

var statusMappings map[string]map[string]string

func getStdStatus(jiraIssue *models.JiraIssue) string {
	if statusMappings == nil {
		statusMappings = make(map[string]map[string]string)
	}
	if _, ok := statusMappings[jiraIssue.Type]; !ok {
		statusMappings[jiraIssue.Type] = getStringMapFromConfig(
			fmt.Sprintf("JIRA_ISSUE_%v_STATUS_MAPPING", jiraIssue.Type),
		)
	}
	statusMapping := statusMappings[jiraIssue.Type]
	if stdStatus, ok := statusMapping[jiraIssue.StatusName]; ok {
		return stdStatus
	}
	return jiraIssue.StatusName
}

func getStringMapFromConfig(key string) map[string]string {
	mapping := make(map[string]string)
	mappingCfg := strings.TrimSpace(config.V.GetString(key))
	if mappingCfg == "" {
		return mapping
	}
	for _, comp := range strings.Split(mappingCfg, ";") {
		comp := strings.TrimSpace(comp)
		if comp == "" {
			continue
		}
		tmp := strings.Split(comp, ":")
		if len(tmp) != 2 {
			panic(fmt.Errorf("invalid config %v: %v", key, comp))
		}
		std := strings.TrimSpace(tmp[0])
		if std == "" {
			panic(fmt.Errorf("invalid config %v: %v", key, comp))
		}
		orgs := tmp[1]
		for _, org := range strings.Split(orgs, ",") {
			org := strings.TrimSpace(org)
			if org == "" {
				continue
			}
			mapping[org] = std
		}
	}
	return mapping
}
