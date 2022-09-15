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

package errors

import (
	"fmt"
	"strings"
)

type (
	// Messages alias for messages of an Error
	Messages []*errMessage
	// MessageType the type of message for an Error
	MessageType int

	// errMessage an abstraction around a given Error's message
	errMessage struct {
		// all the messages associated with an Error. The size will be > 1 if the Error is created using Type.Combine
		msgs []string
	}
)

func (m *errMessage) addMessage(t *Type, msg string) {
	if msg == "" {
		return
	}
	if t.httpCode != 0 {
		msg = fmt.Sprintf("%s (%d)", msg, t.httpCode)
	}
	m.appendMessage(msg)
}

func (m *errMessage) appendMessage(msg string) {
	m.msgs = append(m.msgs, msg)
}

func (m *errMessage) getMessage() string {
	if len(m.msgs) == 0 {
		return ""
	}
	return strings.Join(m.msgs, ",")
}

func (m *errMessage) getPrettifiedMessage() string {
	if len(m.msgs) == 0 {
		return ""
	}
	if len(m.msgs) == 1 {
		return m.msgs[0]
	}
	effectiveMsg := strings.Join(m.msgs, "\n=====================\n")
	effectiveMsg = "\t" + strings.ReplaceAll(effectiveMsg, "\n", "\n\t")
	return fmt.Sprintf("\ncombined messages: \n{\n%s\n}", effectiveMsg)
}

// Format formats the messages into a single string
func (m Messages) Format() string {
	msgs := []string{}
	for _, m := range m {
		if msg := m.getMessage(); msg != "" {
			msgs = append(msgs, msg)
		}
	}
	return strings.Join(msgs, "\ncaused by: ")
}

// Get gets the main (top-level) (or first non-empty message if exists) message of the Messages
func (m Messages) Get() string {
	for _, m := range m {
		if msg := m.getMessage(); msg != "" {
			return msg
		}
	}
	return ""
}

// Causes gets the non-main messages of the Messages in causal sequence
func (m Messages) Causes() []string {
	if len(m) < 2 {
		return nil
	}
	causes := []string{}
	for _, m := range m[1:] {
		if msg := m.getMessage(); msg != "" {
			causes = append(causes, msg)
		}
	}
	return causes
}
