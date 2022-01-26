package store

import (
	"github.com/merico-dev/lake/models/domainlayer/code"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const BathSize = 100

type Database struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{db: db}
}

func (d *Database) RepoCommits(repoCommit *code.RepoCommit) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(repoCommit).Error
}

func (d *Database) Commits(commit *code.Commit) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(commit).Error
}

func (d *Database) Refs(ref *code.Ref) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(ref).Error
}

func (d *Database) CommitFiles(file *code.CommitFile) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(file).Error
}

func (d *Database) CommitParents(pp []*code.CommitParent) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(pp, BathSize).Error
}

func (d Database) Close() error {
	return nil
}
