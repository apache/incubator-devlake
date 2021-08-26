package task

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/utils"
)

func TestNewTask(t *testing.T) {
	r := gin.Default()
	api.RegisterRouter(r)

	w := httptest.NewRecorder()
	params := strings.NewReader(`{ "plugin": "jira", "options": { "host": "www.jira.com" } }`)
	req, _ := http.NewRequest("POST", "/task", params)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)
	resp := w.Body.String()
	task, err := utils.JsonToMap(resp)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, task["Plugin"], "jira")
}
