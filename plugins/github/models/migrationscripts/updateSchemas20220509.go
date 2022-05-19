package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"time"

	"gorm.io/gorm"
)

type GithubPullRequest20220509 struct {
	GithubId        int    `gorm:"primaryKey"`
	RepoId          int    `gorm:"index"`
	Number          int    `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	State           string `gorm:"type:varchar(255)"`
	Title           string `gorm:"type:varchar(255)"`
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	ClosedAt        *time.Time
	// In order to get the following fields, we need to collect PRs individually from GitHub
	Additions      int
	Deletions      int
	Comments       int
	Commits        int
	ReviewComments int
	Merged         bool
	MergedAt       *time.Time
	Body           string
	Type           string `gorm:"type:varchar(255)"`
	Component      string `gorm:"type:varchar(255)"`
	MergeCommitSha string `gorm:"type:varchar(40)"`
	HeadRef        string `gorm:"type:varchar(255)"`
	BaseRef        string `gorm:"type:varchar(255)"`
	BaseCommitSha  string `gorm:"type:varchar(255)"`
	HeadCommitSha  string `gorm:"type:varchar(255)"`
	Url            string `gorm:"type:varchar(255)"`
	AuthorName     string `gorm:"type:varchar(100)"`
	AuthorId       int
	archived.NoPKModel
}

func (GithubPullRequest20220509) TableName() string {
	return "_tool_github_pull_requests"
}

type UpdateSchemas20220509 struct{}

func (*UpdateSchemas20220509) Up(ctx context.Context, db *gorm.DB) error {
	if db.Name() == "postgres" {
		err := db.Exec("ALTER TABLE _tool_github_pull_requests ALTER COLUMN author_id TYPE bigint USING author_id::bigint").Error
		if err != nil {
			return err
		}
		return nil
	}
	err := db.Migrator().AlterColumn(GithubPullRequest20220509{}, "author_id")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220509) Version() uint64 {
	return 20220509212344
}

func (*UpdateSchemas20220509) Owner() string {
	return "Github"
}

func (*UpdateSchemas20220509) Name() string {
	return "Alter pr column `author_id` to int"
}
