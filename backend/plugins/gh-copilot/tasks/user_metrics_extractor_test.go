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

package tasks

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruncateToMaxLengthNoTruncation(t *testing.T) {
	value := "copilot-v1.2.3"
	truncated := truncateToMaxLength(value, lastKnownPluginVersionMaxLength)
	require.Equal(t, value, truncated)
}

func TestTruncateToMaxLengthTruncates(t *testing.T) {
	value := strings.Repeat("a", lastKnownPluginVersionMaxLength+10)
	truncated := truncateToMaxLength(value, lastKnownPluginVersionMaxLength)
	require.Len(t, []rune(truncated), lastKnownPluginVersionMaxLength)
	require.Equal(t, strings.Repeat("a", lastKnownPluginVersionMaxLength), truncated)
}

func TestTruncateToMaxLengthTruncatesByRunes(t *testing.T) {
	value := strings.Repeat("é", lastKnownPluginVersionMaxLength+5)
	truncated := truncateToMaxLength(value, lastKnownPluginVersionMaxLength)
	require.Len(t, []rune(truncated), lastKnownPluginVersionMaxLength)
	require.Equal(t, strings.Repeat("é", lastKnownPluginVersionMaxLength), truncated)
}
