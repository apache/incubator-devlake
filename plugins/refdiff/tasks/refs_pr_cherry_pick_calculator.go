package tasks

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm/clause"
)

func CalculatePrCherryPick(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	ctx := taskCtx.GetContext()
	db := taskCtx.GetDb()
	var prTitleRegex *regexp.Regexp
	prTitlePattern := taskCtx.GetConfig("GITHUB_PR_TITLE_PATTERN")

	if len(prTitlePattern) > 0 {
		fmt.Println(prTitlePattern)
		prTitleRegex = regexp.MustCompile(prTitlePattern)
	}

	cursor, err := db.Model(&code.PullRequest{}).
		Joins("left join repos on pull_requests.base_repo_id = repos.id").
		Where("repos.id = ?", repoId).Rows()
	if err != nil {
		return err
	}

	defer cursor.Close()

	pr := &code.PullRequest{}
	var parentPrKeyInt int
	taskCtx.SetProgress(0, -1)

	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = db.ScanRows(cursor, pr)
		if err != nil {
			return err
		}

		parentPrKey := ""
		if prTitleRegex != nil {
			groups := prTitleRegex.FindStringSubmatch(pr.Title)
			if len(groups) > 0 {
				parentPrKey = groups[1]
			}
		}

		if parentPrKeyInt, err = strconv.Atoi(parentPrKey); err != nil {
			continue
		}

		var parentPrId string
		err = db.Model(&code.PullRequest{}).
			Where("key = ? and repo_id = ?", parentPrKeyInt, repoId).
			Pluck("id", &parentPrId).Error
		if err != nil {
			return err
		}
		if len(parentPrId) == 0 {
			continue
		}
		pr.ParentPrId = parentPrId

		err = db.Save(pr).Error
		if err != nil {
			return err
		}
		taskCtx.IncProgress(1)
	}

	cursor2, err := db.Table("pull_requests pr1").
		Joins("left join pull_requests pr2 on pr1.parent_pr_id = pr2.id").Group("pr1.parent_pr_id, pr2.created_date").Where("pr1.parent_pr_id != ''").
		Joins("left join repos on pr2.base_repo_id = repos.id").
		Order("pr2.created_date ASC").
		Select(`pr2.key as parent_pr_key, pr1.parent_pr_id as parent_pr_id, GROUP_CONCAT(pr1.base_ref order by pr1.base_ref ASC) as cherrypick_base_branches, 
			GROUP_CONCAT(pr1.key order by pr1.base_ref ASC) as cherrypick_pr_keys, repos.name as repo_name, 
			concat(repos.url, '/pull/', pr2.key) as parent_pr_url`).Rows()
	if err != nil {
		return err
	}
	defer cursor2.Close()

	refsPrCherryPick := &code.RefsPrCherrypick{}
	for cursor2.Next() {
		err = db.ScanRows(cursor2, refsPrCherryPick)
		if err != nil {
			return err
		}
		err = db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(refsPrCherryPick).Error
		if err != nil {
			return err
		}
	}
	return nil
}

var CalculatePrCherryPickMeta = core.SubTaskMeta{
	Name:             "calculatePrCherryPick",
	EntryPoint:       CalculatePrCherryPick,
	EnabledByDefault: true,
	Description:      "Calculate pr cherry pick",
}
