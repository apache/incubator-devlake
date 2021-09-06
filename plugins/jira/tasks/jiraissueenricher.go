package tasks

import (
	"fmt"
	"strings"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jira/models"
)

var workloadCoefficient float64

func init() {
	workloadCoefficient = config.V.GetFloat64("JIRA_WORKLOAD_COEFFICIENT")
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

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		if jiraIssue.Fields.ResolutionDate.IsZero() {
			jiraIssue.LeadTime = uint(jiraIssue.Fields.ResolutionDate.Unix()-jiraIssue.Fields.Created.Unix()) / 60
		}
		jiraIssue.StdWorkload = uint(jiraIssue.Workload * workloadCoefficient)
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
	if stdType, ok := typeMapping[jiraIssue.Fields.Issuetype.Name]; ok {
		return stdType
	}
	return jiraIssue.Fields.Issuetype.Name
}

var statusMappings map[string]map[string]string

func getStdStatus(jiraIssue *models.JiraIssue) string {
	if statusMappings == nil {
		statusMappings = make(map[string]map[string]string)
	}
	if _, ok := statusMappings[jiraIssue.Fields.Issuetype.Name]; !ok {
		statusMappings[jiraIssue.Fields.Issuetype.Name] = getStringMapFromConfig(
			fmt.Sprintf("JIRA_ISSUE_%v_STATUS_MAPPING", jiraIssue.Fields.Issuetype.Name),
		)
	}
	statusMapping := statusMappings[jiraIssue.Fields.Issuetype.Name]
	if stdStatus, ok := statusMapping[jiraIssue.Fields.Status.Name]; ok {
		return stdStatus
	}
	return jiraIssue.Fields.Status.Name
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
