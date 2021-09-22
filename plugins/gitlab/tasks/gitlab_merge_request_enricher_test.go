package tasks

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/merico-dev/lake/plugins/gitlab/models"
)

func TestGetReviewRoundsOneReview(t *testing.T) {
	// Handles merge with no comments or commits which makes one review rounds
	commits := []models.GitlabMergeRequestCommit{}
	notes := []models.GitlabMergeRequestNote{}

	numberOfReviewRounds := GetReviewRounds(commits, notes)
	assert.Equal(t, numberOfReviewRounds, 1)
}
func TestGetReviewRoundsTwoReviews(t *testing.T) {
	// Handles one review round plus merge which makes two review rounds
	commits := []models.GitlabMergeRequestCommit{
		{
			AuthoredDate: "a",
		},
	}
	notes := []models.GitlabMergeRequestNote{
		{
			GitlabCreatedAt: "a",
		},
	}

	numberOfReviewRounds := GetReviewRounds(commits, notes)
	assert.Equal(t, numberOfReviewRounds, 2)
}

func TestGetReviewRoundsManyReviews(t *testing.T) {
	// Handles one review round plus merge which makes two review rounds
	// 1. Merge Request created
	// 2. Commit added by dev
	// 3. Feedback given as note -- Review Round
	// 4. Two more commits added
	// 5. Two more notes added -- Review Round
	// 6. Commit added by dev
	// 7. Merged by reviewer -- Review Round
	commits := []models.GitlabMergeRequestCommit{
		{
			AuthoredDate: "a",
		},
		{
			AuthoredDate: "c",
		},
		{
			AuthoredDate: "d",
		},
		{
			AuthoredDate: "g",
		},
	}
	notes := []models.GitlabMergeRequestNote{
		{
			GitlabCreatedAt: "b",
		},
		{
			GitlabCreatedAt: "e",
		},
		{
			GitlabCreatedAt: "f",
		},
	}

	numberOfReviewRounds := GetReviewRounds(commits, notes)
	assert.Equal(t, numberOfReviewRounds, 3)
}
