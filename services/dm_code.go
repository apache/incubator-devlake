package services

import "github.com/merico-dev/lake/models/domainlayer/code"

func GetRepos() ([]*code.Repo, int64, error) {
	repos := make([]*code.Repo, 0)
	db := db.Model(repos).Order("id DESC")
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Find(&repos).Error
	if err != nil {
		return nil, count, err
	}
	return repos, count, nil
}
