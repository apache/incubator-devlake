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

package tap

import (
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
)

// SingerTapArgs the args needed to instantiate tap.Tap for singer-taps
type SingerTapArgs struct {
	// Name of the env variable that expands to the tap binary path
	TapClass string
	// The name of the properties/catalog JSON file of the tap
	StreamPropertiesFile string
	// Optional - Use for any extra tweaking of the required stream at plugin runtime. Usually not needed.
	StreamModifier func(stream *SingerTapStream)
	// IsLegacy - set to true if this is an old tap that uses the "--properties" flag
	IsLegacy bool
}

// NewSingerTapClient returns an instance of tap.Tap for singer-taps
func NewSingerTapClient(args *SingerTapArgs) (Tap, errors.Error) {
	env := config.GetConfig()
	cmd := env.GetString(args.TapClass)
	if cmd == "" {
		return nil, errors.Default.New("singer tap command not provided")
	}
	return NewSingerTap(&SingerTapConfig{
		Cmd:                  cmd,
		StreamPropertiesFile: args.StreamPropertiesFile,
		IsLegacy:             args.IsLegacy,
		// This function is called for the selected streams at runtime.
		TapStreamModifier: func(stream *SingerTapStream) {
			// default behavior
			for _, meta := range stream.Metadata {
				innerMeta := meta["metadata"].(map[string]any)
				innerMeta["selected"] = true
			}
			if args.StreamModifier != nil {
				args.StreamModifier(stream)
			}
		},
	})
}
