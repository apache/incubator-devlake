package store

import (
	"github.com/merico-dev/lake/models/domainlayer/code"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const BathSize = 100

type Database struct {
	db            *gorm.DB
	repoCommits   []*code.RepoCommit
	commits       []*code.Commit
	refs          []*code.Ref
	commitFiles   []*code.CommitFile
	commitParents []*code.CommitParent
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{db: db}
}

func (d *Database) RepoCommits(repoCommit *code.RepoCommit) error {
	d.repoCommits = append(d.repoCommits, repoCommit)
	if len(d.repoCommits) < BathSize {
		return nil
	}
	defer func() { d.repoCommits = make([]*code.RepoCommit, 0, BathSize) }()
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.repoCommits).Error
}

func (d *Database) Commits(commit *code.Commit) error {
	d.commits = append(d.commits, commit)
	if len(d.commits) < BathSize {
		return nil
	}
	defer func() { d.commits = make([]*code.Commit, 0, BathSize) }()
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.commits).Error
}

func (d *Database) Refs(ref *code.Ref) error {
	d.refs = append(d.refs, ref)
	if len(d.refs) < BathSize {
		return nil
	}
	defer func() { d.refs = make([]*code.Ref, 0, BathSize) }()
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.refs).Error
}

func (d *Database) CommitFiles(file *code.CommitFile) error {
	d.commitFiles = append(d.commitFiles, file)
	if len(d.commitFiles) < BathSize {
		return nil
	}
	defer func() { d.commitFiles = make([]*code.CommitFile, 0, BathSize) }()
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.commitFiles).Error
}

func (d *Database) CommitParents(pp []*code.CommitParent) error {
	d.commitParents = append(d.commitParents, pp...)
	if len(d.commitParents) < BathSize {
		return nil
	}
	defer func() { d.commitParents = make([]*code.CommitParent, 0, BathSize) }()
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.commitParents).Error
}

func (d *Database) Flush() error {
	var err error
	if len(d.repoCommits) > 0 {
		err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.repoCommits).Error
		if err != nil {
			return err
		}
	}
	d.repoCommits = make([]*code.RepoCommit, 0, BathSize)
	if len(d.commits) > 0 {
		err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.commits).Error
		if err != nil {
			return err
		}
	}
	d.commits = make([]*code.Commit, 0, BathSize)
	if len(d.refs) > 0 {
		err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.refs).Error
		if err != nil {
			return err
		}
	}
	d.refs = make([]*code.Ref, 0, BathSize)
	if len(d.commitFiles) > 0 {
		err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.commitFiles).Error
		if err != nil {
			return err
		}
	}
	d.commitFiles = make([]*code.CommitFile, 0, BathSize)
	if len(d.commitParents) > 0 {
		err = d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(d.commitParents).Error
		if err != nil {
			return err
		}
	}
	d.commitParents = make([]*code.CommitParent, 0, BathSize)
	return nil
}

func (d *Database) Close() error {
	return nil
}
