package helper

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	dalMocks "github.com/apache/incubator-devlake/plugins/core/dal/mocks"
	coreMocks "github.com/apache/incubator-devlake/plugins/core/mocks"
	"github.com/apache/incubator-devlake/plugins/helper/common"
	helperMocks "github.com/apache/incubator-devlake/plugins/helper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchPageUndetermined(t *testing.T) {
	mockDal := new(dalMocks.Dal)
	mockDal.On("AutoMigrate", mock.Anything, mock.Anything).Return(nil).Once()
	mockDal.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockDal.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	mockLog := new(coreMocks.Logger)
	mockLog.On("Info", mock.Anything, mock.Anything)
	mockLog.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLog.On("Debug", mock.Anything, mock.Anything, mock.Anything).Maybe()

	mockCtx := new(coreMocks.SubTaskContext)
	mockCtx.On("GetDal").Return(mockDal)
	mockCtx.On("GetLogger").Return(mockLog)
	mockCtx.On("SetProgress", mock.Anything, mock.Anything)
	mockCtx.On("IncProgress", mock.Anything, mock.Anything)
	mockCtx.On("GetName").Return("test")

	mockInput := new(helperMocks.Iterator)
	mockInput.On("HasNext").Return(true).Once()
	mockInput.On("HasNext").Return(false).Once()
	mockInput.On("Fetch").Return(nil, nil).Once()
	mockInput.On("Close").Return(nil)

	// simulate fetching all pages of jira changelogs for 1 issue id with 1 concurrency,
	// assuming api doesn't return total number of pages.
	// then, we are expecting 2 calls for GetAsync and NextTick each, otherwise, deadlock happens
	getAsyncCounter := 0
	mockApi := new(helperMocks.RateLimitedApiClient)
	mockApi.On("GetAsync", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// fake records for first page, no records for second page
		body := "[1,2,3]"
		if getAsyncCounter > 0 {
			body = "[]"
		}
		getAsyncCounter += 1
		res := &http.Response{
			Request: &http.Request{
				URL: &url.URL{},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		handler := args.Get(3).(common.ApiAsyncCallback)
		handler(res)
	}).Twice()
	mockApi.On("NextTick", mock.Anything).Run(func(args mock.Arguments) {
		handler := args.Get(0).(func() error)
		assert.Nil(t, handler())
	}).Twice()
	mockApi.On("HasError").Return(false)
	mockApi.On("WaitAsync").Return(nil)

	params := struct {
		Name string
	}{Name: "testparams"}

	collector, err := NewApiCollector(ApiCollectorArgs{
		RawDataSubTaskArgs: RawDataSubTaskArgs{
			Ctx:    mockCtx,
			Table:  "whatever rawtable",
			Params: params,
		},
		ApiClient:      mockApi,
		Input:          mockInput,
		UrlTemplate:    "whatever url",
		Concurrency:    1,
		PageSize:       3,
		ResponseParser: GetRawMessageArrayFromResponse,
	})

	assert.Nil(t, err)
	assert.Nil(t, collector.Execute())

	mockDal.AssertExpectations(t)
}
