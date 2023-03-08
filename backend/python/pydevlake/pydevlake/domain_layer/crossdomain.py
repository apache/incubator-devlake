# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from datetime import datetime

from sqlmodel import Field

from pydevlake.model import DomainModel, NoPKModel


class Account(DomainModel, table=True):
    __tablename__ = "accounts"
    email: str
    full_name: str
    user_name: str
    avatar_url: str
    organization: str
    created_date: datetime
    status: int


class BoardRepo(NoPKModel, table=True):
    __tablename__ = 'board_repos'
    board_id: str = Field(primary_key=True)
    repo_id: str = Field(primary_key=True)


class IssueCommit(NoPKModel, table=True):
    __tablename__ = 'issue_commits'
    issue_id: str = Field(primary_key=True)
    commit_sha: str = Field(primary_key=True)


class IssueRepoCommit(NoPKModel, table=True):
    __tablename__ = 'issue_repo_commits'
    issue_id: str = Field(primary_key=True)
    repo_url: str = Field(primary_key=True)
    commit_sha: str = Field(primary_key=True)


class ProjectIssueMetric(NoPKModel, table=True):
    __tablename__ = "project_issue_metrics"
    project_name: str = Field(primary_key=True)
    deployment_id: str


class ProjectMapping(NoPKModel, table=True):
    __tablename__ = "project_mappings"
    project_name: str = Field(primary_key=True)
    table: str = Field(primary_key=True)
    row_id: str = Field(primary_key=True)


class ProjectPrMetric(DomainModel, table=True):
    __tablename__ = "project_pr_metrics"
    project_name: str = Field(primary_key=True)
    first_commit_sha: str
    pr_coding_time: int
    first_review_id: str
    pr_pick_time: int
    pr_review_time: int
    deployment_id: str
    pr_deploy_time: int
    pr_cycle_time: int


class PullRequestIssue(NoPKModel, table=True):
    __tablename__ = "pull_request_issues"
    pull_request_id: str = Field(primary_key=True)
    issue_id: str = Field(primary_key=True)
    pull_request_key: int
    issue_key: int


class RefsIssuesDiffs(NoPKModel, table=True):
    __tablename__ = "refs_issues_diffs"
    new_ref_id: str = Field(primary_key=True)
    old_ref_id: str = Field(primary_key=True)
    new_ref_commit_sha: str
    old_ref_commit_sha: str
    issue_number: str
    issue_id: str = Field(primary_key=True)


class Team(DomainModel, table=True):
    __tablename__ = "teams"
    name: str
    alias: str
    parent_id: str
    sorting_index: int


class TeamUser(NoPKModel, table=True):
    __tablename__ = "team_users"
    team_id: str = Field(primary_key=True)
    user_id: str = Field(primary_key=True)


class User(DomainModel, table=True):
    __tablename__ = 'users'
    name: str
    email: str


class UserAccount(NoPKModel, table=True):
    __tablename__ = 'user_accounts'
    user_id: str
    account_id: str = Field(primary_key=True)
