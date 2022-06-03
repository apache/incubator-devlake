package unithelper

import (
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/stretchr/testify/mock"
)

func DummySubTaskContext(db dal.Dal) *mocks.SubTaskContext {
	mockCtx := new(mocks.SubTaskContext)
	mockCtx.On("GetDal").Return(db)
	mockCtx.On("GetLogger").Return(DummyLogger())
	mockCtx.On("SetProgress", mock.Anything, mock.Anything)
	mockCtx.On("IncProgress", mock.Anything, mock.Anything)
	mockCtx.On("GetName").Return("test")
	return mockCtx
}
