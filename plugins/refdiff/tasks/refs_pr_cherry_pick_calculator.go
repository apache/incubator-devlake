package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
)

var prTitleRegex *regexp.Regexp
var prTitlePattern string

func init() {
	prTitlePattern = config.GetConfig().GetString("GITHUB_PR_TITLE_PATTERN")
}

func CalculatePrCherryPick(ctx context.Context, pairs []RefPair, progress chan<- float32, repoId string) error {
	if len(prTitlePattern) > 0 {
		fmt.Println(prTitlePattern)
		prTitleRegex = regexp.MustCompile(prTitlePattern)
	}

	cursor, err := lakeModels.Db.Model(&code.PullRequest{}).
		Joins("left join repos on pull_requests.repo_id = repos.id").
		Where("repos.id = ?", repoId).Rows()
	if err != nil {
		return err
	}

	defer cursor.Close()

	pr := &code.PullRequest{}
	var parentPrKeyInt int

	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = lakeModels.Db.ScanRows(cursor, pr)
		if err != nil {
			return err
		}

		parentPrKey := getParentPrKey(pr.Title)

		if parentPrKeyInt, err = strconv.Atoi(parentPrKey); err != nil {
			continue
		}

		var parentPrId string
		err = lakeModels.Db.Model(&code.PullRequest{}).
			Where("`key` = ? and repo_id = ?", parentPrKeyInt, repoId).
			Pluck("id", &parentPrId).Error
		if err != nil {
			return err
		}
		if len(parentPrId) == 0 {
			continue
		}
		pr.ParentPrId = parentPrId

		err = lakeModels.Db.Save(pr).Error
		if err != nil {
			return err
		}
	}

	progress <- 0.90

	cursor2, err := lakeModels.Db.Table("pull_requests pr1").
		Joins("left join pull_requests pr2 on pr1.parent_pr_id = pr2.id").Group("pr1.parent_pr_id, pr2.created_date").Where("pr1.parent_pr_id != ''").
		Joins("left join repos on pr2.repo_id = repos.id").
		Order("pr2.created_date ASC").
		Select("pr2.`key` as parent_pr_key, pr1.parent_pr_id as parent_pr_id, GROUP_CONCAT(pr1.base_ref order by pr1.base_ref ASC) as cherrypick_base_branches, " +
			"GROUP_CONCAT(pr1.`key` order by pr1.base_ref ASC) as cherrypick_pr_keys, repos.`name` as repo_name, " +
			"concat(repos.url, '/pull/', pr2.`key`) as parent_pr_url").Rows()
	defer cursor2.Close()

	refsPrCherryPick := &code.RefsPrCherrypick{}
	for cursor2.Next() {
		err = lakeModels.Db.ScanRows(cursor2, refsPrCherryPick)
		if err != nil {
			return err
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(refsPrCherryPick).Error
		if err != nil {
			return err
		}
	}
	progress <- 1

	return nil
}

func getParentPrKey(title string) string {
	if prTitleRegex != nil {
		groups := prTitleRegex.FindStringSubmatch(title)
		if len(groups) > 0 {
			return groups[1]
		}
	}
	return ""
}
