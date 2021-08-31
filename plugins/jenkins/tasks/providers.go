package tasks

import (
	"github.com/merico-dev/lake/plugins/jenkins/models"
)

type JenkinsStorage interface {
	SaveJob(job models.JenkinsJobProps) (context interface{}, err error)
	SaveJobs(jobs []models.JenkinsJobProps) (context interface{}, err error)
	SaveBuild(build models.JenkinsBuildProps, ctx interface{}) (context interface{}, err error)
	SaveBuilds(build []models.JenkinsBuildProps, ctx interface{}) (context interface{}, err error)
}
