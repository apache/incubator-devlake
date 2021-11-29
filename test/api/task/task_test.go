package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/models"
	"github.com/stretchr/testify/mock"
)

func TestNewTask(t *testing.T) {
	r := gin.Default()
	api.RegisterRouter(r)

	type services struct {
		mock.Mock
	}

	// fakeTask := models.Task{}
	testObj := new(services)
	testObj.On("CreateTask").Return(true, nil)

	w := httptest.NewRecorder()
	params := strings.NewReader(`{"name": "hello", "tasks": [[{ "plugin": "jira", "options": { "host": "www.jira.com" } }]]}`)
	req, _ := http.NewRequest("POST", "/pipelines", params)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)
	resp := w.Body.String()
	var pipeline models.Pipeline
	err := json.Unmarshal([]byte(resp), &pipeline)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, pipeline.Name, "hello")

	var tasks [][]*models.NewTask
	err = json.Unmarshal(pipeline.Tasks, &tasks)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tasks[0][0].Plugin, "jira")
}
