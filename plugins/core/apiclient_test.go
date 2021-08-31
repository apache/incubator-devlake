package core

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestStruct - our JSON unmarshalling struct
type TestStruct struct {
	Name string
	Age  int `json:"realAge"`
}

// TestRunner - struct we can keep updating with new test arguments
type TestRunner struct {
	Container *TestStruct
	Response  *http.Response
}

// updateArgs - Updates the arguments for the next test run
func (tr *TestRunner) updateArgs(newTestInput string) {
	tr.Response.Body = io.NopCloser(strings.NewReader(newTestInput))
	tr.Container = &TestStruct{}
}

// TestUnmarshalResponse - Tests for UnmarshalResponse
// Since UnmarshalResponse basically just wraps the built-in java decoder
// this test is not super necessary. That being said, it's good to have
// in case we decide to change the implementation of UnmarshalResponse
// in the future
func TestUnmarshalResponse(t *testing.T) {

	tr := &TestRunner{
		Container: &TestStruct{},
		Response:  &http.Response{},
	}

	// Try with valid JSON, we shouldn't get an error
	tr.updateArgs(`{"Name": "Jane"}`)
	err := UnmarshalResponse(tr.Response, tr.Container)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Try with invalid JSON, we should get an error
	tr.updateArgs(`{"Name": "}`)
	err = UnmarshalResponse(tr.Response, tr.Container)
	if err == nil {
		t.Errorf(err.Error())
	}

	// Try with empty string, should fail
	tr.updateArgs(``)
	err = UnmarshalResponse(tr.Response, tr.Container)
	if err == nil {
		t.Errorf(err.Error())
	}

	// Try with second value in json object
	tr.updateArgs(`{"Name": "Jane", "realAge": 30}`)
	err = UnmarshalResponse(tr.Response, tr.Container)
	if err != nil {
		t.Errorf(err.Error())
	}
	if tr.Container.Age != 30 {
		t.Errorf("Age is not %d", 30)
	}
}
