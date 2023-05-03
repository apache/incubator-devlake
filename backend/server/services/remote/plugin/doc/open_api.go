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

package doc

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/core/config"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/server/services/remote/models"
)

func GenerateOpenApiSpec(pluginInfo *models.PluginInfo) (*string, errors.Error) {
	connectionSchema, err := json.Marshal(pluginInfo.ConnectionModelInfo.JsonSchema)
	if err != nil {
		return nil, errors.Default.Wrap(err, "connection schema is not valid JSON")
	}
	scopeSchema, err := json.Marshal(pluginInfo.ScopeModelInfo.JsonSchema)
	if err != nil {
		return nil, errors.Default.Wrap(err, "scope schema is not valid JSON")
	}
	txRuleSchema, err := json.Marshal(pluginInfo.TransformationRuleModelInfo.JsonSchema)
	if err != nil {
		return nil, errors.Default.Wrap(err, "transformation rule schema is not valid JSON")
	}
	specTemplate, tmplErr := specTemplate()
	if tmplErr != nil {
		return nil, tmplErr
	}
	writer := &strings.Builder{}
	err = specTemplate.Execute(writer, map[string]interface{}{
		"PluginName":               pluginInfo.Name,
		"ConnectionSchema":         string(connectionSchema),
		"ScopeSchema":              string(scopeSchema),
		"TransformationRuleSchema": string(txRuleSchema),
	})
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not execute swagger doc template")
	}
	doc := writer.String()
	return &doc, nil
}

func specTemplate() (*template.Template, errors.Error) {
	path := config.GetConfig().GetString("SWAGGER_DOCS_DIR")
	if path == "" {
		return nil, errors.Default.New("path for Swagger docs resources is not set")
	}
	file, err := os.Open(filepath.Join(path, "open_api_spec.json.tmpl"))
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not open swagger doc template")
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not read swagger doc template")
	}
	specTemplate, err := template.New("doc").Parse(string(contents))
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not parse swagger doc template")
	}
	return specTemplate, nil
}
