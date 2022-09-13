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

const (
	// RawMessageType filter by raw messages
	RawMessageType MessageType = iota
	// UserMessageType filter by user messages
	UserMessageType
)

type (
	// Messages alias for messages of an Error
	Messages []*errMessage
	// MessageType the type of message for an Error
	MessageType int

	errMessage struct {
		raw  []string
		user []string
	}
)

func (m *errMessage) addMessage(t *Type, raw string, user string, overrideUserMsg bool) {
	if overrideUserMsg {
		user = raw
	}
	if user != "" && raw != user {
		raw = fmt.Sprintf("%s [%s]", raw, user)
	}
	if raw == "" {
		return
	}
	if t.httpCode != 0 {
		raw = fmt.Sprintf("%s (%d)", raw, t.httpCode)
		user = fmt.Sprintf("%s (%d)", user, t.httpCode)
	}
	m.appendMessage(raw, user)
}

func (m *errMessage) appendMessage(raw string, user string) {
	m.raw = append(m.raw, raw)
	m.user = append(m.user, user)
}

func (m *errMessage) getMessage(messageType MessageType) string {
	f := func(target []string) string {
		if len(target) == 0 {
			return ""
		}
		return strings.Join(target, ",")
	}
	if messageType == RawMessageType {
		return f(m.raw)
	}
	return f(m.user)
}

func (m *errMessage) getPrettifiedMessage(messageType MessageType) string {
	f := func(target []string) string {
		if len(target) == 0 {
			return ""
		}
		if len(target) == 1 {
			return target[0]
		}
		effectiveMsg := strings.Join(target, "\n=====================\n")
		effectiveMsg = "\t" + strings.ReplaceAll(effectiveMsg, "\n", "\n\t")
		return fmt.Sprintf("\ncombined messages: \n{\n%s\n}", effectiveMsg)
	}
	if messageType == RawMessageType {
		return f(m.raw)
	}
	return f(m.user)
}

// Format formats the messages into a single string
func (m Messages) Format(messageType MessageType) string {
	msgs := []string{}
	for _, m := range m {
		if msg := m.getMessage(messageType); msg != "" {
			msgs = append(msgs, msg)
		}
	}
	return strings.Join(msgs, "\ncaused by: ")
}

// Get gets the main (top-level) (or first non-empty message if exists) message of the Messages
func (m Messages) Get(messageType MessageType) string {
	for _, m := range m {
		if msg := m.getMessage(messageType); msg != "" {
			return msg
		}
	}
	return ""
}

// Causes gets the non-main messages of the Messages in causal sequence
func (m Messages) Causes(messageType MessageType) []string {
	if len(m) < 2 {
		return nil
	}
	causes := []string{}
	for _, m := range m[1:] {
		if msg := m.getMessage(messageType); msg != "" {
			causes = append(causes, msg)
		}
	}
	return causes
}
