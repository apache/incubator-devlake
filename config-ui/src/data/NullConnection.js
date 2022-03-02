const NullConnection = {
  // id: null,
  ID: null,
  name: null,
  endpoint: null,
  proxy: null,
  token: null,
  username: null,
  password: null,
  basicAuthEncoded: null, // NOTE: we probably want to exclude/null this when exposing this object
  JIRA_ISSUE_TYPE_MAPPING: null,
  JIRA_ISSUE_EPIC_KEY_FIELD: null,
  JIRA_ISSUE_STORYPOINT_FIELD: null,
  JIRA_BOARD_GITLAB_PROJECTS: null,
  JIRA_ISSUE_INCIDENT_STATUS_MAPPING: null,
  JIRA_ISSUE_STORY_STATUS_MAPPING: null,
  createdAt: null,
  updatedAt: null,
}

export {
  NullConnection
}
