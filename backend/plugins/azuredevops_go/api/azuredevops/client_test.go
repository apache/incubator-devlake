/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azuredevops

import (
	"bytes"
	"fmt"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/rogpeppe/go-internal/txtar"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type response struct {
	Body   string
	Header http.Header
	Status int
}

type responseMap map[string]response

func TestRetrieveUserProfile(t *testing.T) {
	files, _ := filepath.Glob("testdata/test.txt")
	if len(files) != 1 {
		t.Fatalf("no testdata")
	}

	a, err := txtar.ParseFile(files[0])
	if err != nil {
		t.Fatalf("testdata/test.txt, failed to parse txtar archive. Err: %s", err)
	}

	if len(a.Files) != 2 {
		t.Fatalf("%s, want two files (request & response), found %d", files, len(a.Files))
	}

	responses := buildResponses(a.Files[1].Data)

	ts := buildMockServer(t, responses)
	defer ts.Close()

	lines := bytes.Split(a.Files[0].Data, []byte("\n"))
	for lineno, line := range lines {
		lineno++

		if strings.TrimSpace(string(line)) == "" {
			continue
		}

		l := strings.Fields(string(line))
		if len(l) != 3 {
			t.Errorf("wrong field count at line %d", lineno)
			continue
		}
		token, res := l[0], l[2]
		code, err := strconv.Atoi(l[1])
		if err != nil {
			t.Errorf("failed to parse status code %s from string to int. Err: %s", l[1], err)
		}

		t.Run(fmt.Sprintf("test Azure DevOps connection with %s", token), func(t *testing.T) {
			conn := &models.AzuredevopsConnection{
				BaseConnection: api.BaseConnection{},
				AzuredevopsConn: models.AzuredevopsConn{
					AzuredevopsAccessToken: models.AzuredevopsAccessToken{
						Token: token,
					},
				},
			}

			client := NewClient(conn, nil, ts.URL)
			p, err := client.GetUserProfile()
			if err != nil && err.GetType().GetHttpCode() != code {
				t.Errorf("User Profile API Response = %d; want: %d", err.GetType().GetHttpCode(), code)
			}

			if code == 200 && p.Id != res {
				t.Errorf("User Profile Id = %q; want %q", p.Id, res)
			}
		})

	}
}

func buildMockServer(t *testing.T, responses responseMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerAuth := r.Header.Get("Authorization")
		var res response
		var ok bool

		if res, ok = responses[headerAuth]; !ok {
			t.Fatalf("no response found for authorization header %q", headerAuth)
			return
		}

		w.WriteHeader(res.Status)
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		_, err := io.WriteString(w, res.Body)
		if err != nil {
			t.Errorf("failed write mock server response. Err: %s", err)
		}

	}))
}

func buildResponses(text []byte) responseMap {
	lines := bytes.Split(text, []byte("\n"))
	res := make(responseMap)
	var key string
	var body *response
	for _, line := range lines {
		line = bytes.TrimRight(line, "\t")

		if len(line) == 0 && body == nil {
			continue
		} else if len(line) == 0 && body != nil {
			res[key] = *body
			body = nil
			continue
		}

		if !bytes.HasPrefix(line, []byte("\t")) {
			key = string(line)
			continue
		}

		if body == nil {
			body = &response{
				Header: make(http.Header),
			}
		}

		line = bytes.TrimPrefix(line, []byte("\t"))
		parts := bytes.Split(line, []byte(":"))

		if len(parts) < 2 {
			continue
		}

		id := string(parts[0])
		value := string(bytes.Join(parts[1:], []byte(":")))
		value = strings.TrimSpace(value)
		switch id {
		case "StatusCode":
			body.Status, _ = strconv.Atoi(value)
		case "Body":
			body.Body = value
		default:
			body.Header.Add(id, value)
		}
	}

	return res
}
