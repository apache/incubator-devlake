package tasks

import (
	"sync"

	"github.com/merico-dev/lake/plugins/gitlab/models"
)

type ConcurrentCommits struct {
	sync.RWMutex
	commits []models.GitlabCommit
}
type ConcurrentProjectCommits struct {
	sync.RWMutex
	projectCommits []models.GitlabProjectCommit
}
type ConcurrentUsers struct {
	sync.RWMutex
	users []models.GitlabUser
}

func (cs *ConcurrentCommits) Append(commit *models.GitlabCommit) {
	cs.Lock()
	defer cs.Unlock()

	cs.commits = append(cs.commits, *commit)
}
func (cs *ConcurrentProjectCommits) Append(pjCommit *models.GitlabProjectCommit) {
	cs.Lock()
	defer cs.Unlock()

	cs.projectCommits = append(cs.projectCommits, *pjCommit)
}
func (cs *ConcurrentUsers) Append(user *models.GitlabUser) {
	cs.Lock()
	defer cs.Unlock()

	cs.users = append(cs.users, *user)
}
