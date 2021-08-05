const mockRawMergeRequests = [
  {
      "id": 109951882,
      "iid": 1,
      "project_id": 28270340,
      "title": "goodbye message",
      "description": "",
      "state": "merged",
      "created_at": "2021-07-27T14:10:11.627Z",
      "updated_at": "2021-07-27T14:56:24.887Z",
      "merged_by": {
          "id": 4126108,
          "name": "Kevin Kline",
          "username": "kevin-kline",
          "state": "active",
          "avatar_url": "https://gitlab.com/uploads/-/system/user/avatar/4126108/avatar.png",
          "web_url": "https://gitlab.com/kevin-kline"
      },
      "merged_at": "2021-07-27T14:56:24.912Z",
      "closed_by": null,
      "closed_at": null,
      "target_branch": "master",
      "source_branch": "test-mr",
      "user_notes_count": 1,
      "upvotes": 0,
      "downvotes": 0,
      "author": {
          "id": 4126108,
          "name": "Kevin Kline",
          "username": "kevin-kline",
          "state": "active",
          "avatar_url": "https://gitlab.com/uploads/-/system/user/avatar/4126108/avatar.png",
          "web_url": "https://gitlab.com/kevin-kline"
      },
      "assignees": [
          {
              "id": 4126108,
              "name": "Kevin Kline",
              "username": "kevin-kline",
              "state": "active",
              "avatar_url": "https://gitlab.com/uploads/-/system/user/avatar/4126108/avatar.png",
              "web_url": "https://gitlab.com/kevin-kline"
          }
      ],
      "assignee": {
          "id": 4126108,
          "name": "Kevin Kline",
          "username": "kevin-kline",
          "state": "active",
          "avatar_url": "https://gitlab.com/uploads/-/system/user/avatar/4126108/avatar.png",
          "web_url": "https://gitlab.com/kevin-kline"
      },
      "reviewers": [
          {
              "id": 4126108,
              "name": "Kevin Kline",
              "username": "kevin-kline",
              "state": "active",
              "avatar_url": "https://gitlab.com/uploads/-/system/user/avatar/4126108/avatar.png",
              "web_url": "https://gitlab.com/kevin-kline"
          }
      ],
      "source_project_id": 28270340,
      "target_project_id": 28270340,
      "labels": [],
      "draft": false,
      "work_in_progress": false,
      "milestone": null,
      "merge_when_pipeline_succeeds": false,
      "merge_status": "can_be_merged",
      "sha": "926ccda073f04f12eb7ab373c84cb73cff4ce238",
      "merge_commit_sha": "1c64834aa4dcb6233bf7c125df6947140d027392",
      "squash_commit_sha": null,
      "discussion_locked": null,
      "should_remove_source_branch": null,
      "force_remove_source_branch": true,
      "reference": "!1",
      "references": {
          "short": "!1",
          "relative": "!1",
          "full": "kevin-kline/test-project!1"
      },
      "web_url": "https://gitlab.com/kevin-kline/test-project/-/merge_requests/1",
      "time_stats": {
          "time_estimate": 0,
          "total_time_spent": 0,
          "human_time_estimate": null,
          "human_total_time_spent": null
      },
      "squash": false,
      "task_completion_status": {
          "count": 0,
          "completed_count": 0
      },
      "has_conflicts": false,
      "blocking_discussions_resolved": true,
      "approvals_before_merge": null
  }
]

module.exports = mockRawMergeRequests
