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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFindEarliestNote(t *testing.T) {
	baseTime, err := time.Parse(time.RFC3339, "2022-01-02T15:04:05Z")
	assert.Nil(t, err)
	// Create some sample notes
	note1 := models.GitlabMrNote{Resolvable: true, GitlabCreatedAt: baseTime.Add(-time.Hour)}
	note2 := models.GitlabMrNote{Resolvable: false, GitlabCreatedAt: baseTime.Add(-time.Minute)}
	note3 := models.GitlabMrNote{Resolvable: true, GitlabCreatedAt: baseTime.Add(-time.Second)}

	// Call the function with the sample notes
	notes := []models.GitlabMrNote{note1, note2, note3}
	earliestNote, err := findEarliestNote(notes)

	// Check that no error was returned
	if err != nil {
		t.Errorf("findEarliestNote returned an error: %v", err)
	}

	// Check that the correct note was returned
	if earliestNote == nil {
		t.Errorf("findEarliestNote returned nil, expected a note")
	}
	if !assert.Equal(t, note1, *earliestNote) {
		t.Errorf("findEarliestNote returned the wrong note: got %v, expected %v", earliestNote, &note1)
	}

	// Modify one of the notes to make it unresolvable
	notes[0].Resolvable = false

	// Call the function again with the modified notes
	earliestNote, err = findEarliestNote(notes)

	// Check that no error was returned
	if err != nil {
		t.Errorf("findEarliestNote returned an error: %v", err)
	}

	// Check that the correct note was returned
	if earliestNote == nil {
		t.Errorf("findEarliestNote returned nil, expected a note")
	}
	if !assert.Equal(t, note3, *earliestNote) {
		t.Errorf("findEarliestNote returned the wrong note: got %v, expected %v", earliestNote, &note3)
	}
}

func TestGetReviewRounds(t *testing.T) {
	baseTime, err := time.Parse(time.RFC3339, "2022-01-02T15:04:05Z")
	assert.Nil(t, err)
	// Test case 1: empty input
	var commits []models.GitlabCommit
	var notes []models.GitlabMrNote
	expected := 1
	if got := getReviewRounds(commits, notes); got != expected {
		t.Errorf("getReviewRounds(%v, %v) = %d, expected %d", commits, notes, got, expected)
	}

	// Test case 2: single comment
	commits = []models.GitlabCommit{
		{AuthoredDate: baseTime.Add(-time.Hour * 2)},
		{AuthoredDate: baseTime.Add(-time.Hour)},
	}
	notes = []models.GitlabMrNote{
		{GitlabCreatedAt: baseTime},
	}
	expected = 1
	if got := getReviewRounds(commits, notes); got != expected {
		t.Errorf("getReviewRounds(%v, %v) = %d, expected %d", commits, notes, got, expected)
	}

	// Test case 3: multiple comments
	commits = []models.GitlabCommit{
		{AuthoredDate: baseTime.Add(-time.Hour * 15)},
		{AuthoredDate: baseTime.Add(-time.Hour * 9)},
		{AuthoredDate: baseTime.Add(-time.Hour * 3)},
		{AuthoredDate: baseTime},
	}
	notes = []models.GitlabMrNote{
		{GitlabCreatedAt: baseTime.Add(-time.Hour * 14)},
		{GitlabCreatedAt: baseTime.Add(-time.Hour * 7)},
		{GitlabCreatedAt: baseTime.Add(-time.Hour * 2)},
	}
	expected = 4
	if got := getReviewRounds(commits, notes); got != expected {
		t.Errorf("getReviewRounds(%v, %v) = %d, expected %d", commits, notes, got, expected)
	}
}
