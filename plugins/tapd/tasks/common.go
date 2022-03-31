package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/didgen"
)

const shortForm = "2006-01-02"
const longForm = "2006-01-02 15:04:05"

var UserIdGen *didgen.DomainIdGenerator
var WorkspaceIdGen *didgen.DomainIdGenerator
