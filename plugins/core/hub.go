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

package core

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Allowing plugin to know each other

var plugins map[string]PluginMeta
var mutex sync.Mutex

func RegisterPlugin(name string, plugin PluginMeta) error {
	if plugins == nil {
		plugins = make(map[string]PluginMeta)
	}
	mutex.Lock()
	defer mutex.Unlock()
	plugins[name] = plugin
	return nil
}

func GetPlugin(name string) (PluginMeta, error) {
	if plugins == nil {
		return nil, errors.New("RegisterPlugin have never been called.")
	}
	mutex.Lock()
	defer mutex.Unlock()
	if plugin, ok := plugins[name]; ok {
		return plugin, nil
	}
	return nil, fmt.Errorf("Plugin `%s` doesn't exist", name)
}

type PluginCallBack func(name string, plugin PluginMeta) error

func TraversalPlugin(handle PluginCallBack) error {
	for name, plugin := range plugins {
		err := handle(name, plugin)
		if err != nil {
			return err
		}
	}
	return nil
}

func AllPlugins() map[string]PluginMeta {
	return plugins
}

func FindPluginNameBySubPkgPath(subPkgPath string) (string, error) {
	for name, plugin := range plugins {
		if strings.HasPrefix(subPkgPath, plugin.RootPkgPath()) {
			return name, nil
		}
	}
	return "", fmt.Errorf("Unable to find plugin for subPkgPath %s", subPkgPath)
}
