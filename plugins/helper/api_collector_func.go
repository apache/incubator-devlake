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

package helper

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"io"
	"net/http"
)

// GetRawMessageDirectFromResponse FIXME ...
func GetRawMessageDirectFromResponse(res *http.Response) ([]json.RawMessage, error) {
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return []json.RawMessage{body}, nil
}

// GetRawMessageArrayFromResponse FIXME ...
func GetRawMessageArrayFromResponse(res *http.Response) ([]json.RawMessage, error) {
	rawMessages := []json.RawMessage{}

	if res == nil {
		return nil, errors.Default.New("res is nil")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error reading response body of %s", res.Request.URL.String()))
	}

	err = json.Unmarshal(resBody, &rawMessages)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error decoding response of %s: raw response: %s", res.Request.URL.String(), string(resBody)))
	}

	return rawMessages, nil
}
