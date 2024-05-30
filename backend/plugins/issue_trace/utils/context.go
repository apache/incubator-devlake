package utils

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

func CheckCancel(taskCtx plugin.SubTaskContext) errors.Error {
	ctx := taskCtx.GetContext()
	select {
	case <-ctx.Done():
		return errors.Convert(ctx.Err())
	default:
	}
	return nil
}
