package tasks

//import (
//	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
//	"time"
//)
//
//func getStage(t time.Time, sprintStart, sprintComplete *time.Time) *string {
//	if sprintStart == nil {
//		return &ticket.BeforeSprint
//	}
//	if sprintStart.After(t) {
//		return &ticket.BeforeSprint
//	}
//	if sprintComplete == nil {
//		return &ticket.DuringSprint
//	}
//	if sprintComplete.Before(t) {
//		return &ticket.AfterSprint
//	}
//	return &ticket.DuringSprint
//}
