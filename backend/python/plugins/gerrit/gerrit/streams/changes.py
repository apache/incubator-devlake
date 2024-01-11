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
from time import time
from typing import Iterable
import json

from pydevlake.model import ToolModel
from pydevlake import Stream, DomainType
from pydevlake.context import Context
import pydevlake.domain_layer.code as code
from gerrit.models import GerritChange, GerritProject
from gerrit.api import GerritApi


class GerritChanges(Stream):
    tool_model = GerritChange
    domain_types = [DomainType.CODE]

    def should_run_on(self, scope: GerritProject) -> bool:
        return True

    def collect(self, state, context) -> Iterable[tuple[object, dict]]:
        api = GerritApi(context.connection)
        project: GerritProject = context.scope
        response = api.changes(project.name)
        for raw_change in response.json:
            yield raw_change, state

    def extract(self, raw_data: dict) -> ToolModel:
        def get_localtime(utctime):
            ts = time()
            ts_now = datetime.fromtimestamp(ts)
            ts_utc_now = datetime.utcfromtimestamp(ts)
            offset = ts_now - ts_utc_now
            return utctime + offset

        def get_time_from_text(text):
            return datetime.strptime(text, "%Y-%m-%d %H:%M:%S.%f000")

        utc_datetime = get_time_from_text(raw_data["created"])
        raw_data["created_date"] = utc_datetime

        if raw_data["status"] == "MERGED":
            utc_datetime = get_time_from_text(raw_data["updated"])
            raw_data["merged_date"] = utc_datetime
        if raw_data["status"] == "ABANDONED":
            utc_datetime = get_time_from_text(raw_data["updated"])
            raw_data["closed_date"] = utc_datetime
        revisions = raw_data.get("revisions", {})
        saved_revisions_data = {}
        # we only need few fields from revisions
        for commit_id, data in revisions.items():
            saved_revisions_data[commit_id] = {
                "author_name": data["commit"]["author"]["name"],
                "author_email": data["commit"]["author"]["email"],
                "author_date": data["commit"]["author"]["date"],
                "parent_commit_id": data["commit"]["parents"][0]["commit"],
            }
        raw_data["revisions_json"] = json.dumps(saved_revisions_data)
        return super().extract(raw_data)

    def convert(self, change: GerritChange, ctx: Context):
        def get_status():
            if change.status == "MERGED":
                return "MERGED"
            elif change.status == "ABANDONED":
                return "CLOSED"
            return "OPEN"

        project: GerritProject = ctx.scope
        status = get_status()
        repo_id = project.domain_id()
        base_repo_id = repo_id

        merge_commit_sha = None
        if change.status == "MERGED":
            merge_commit_sha = change.current_revision

        yield code.PullRequest(
            base_repo_id=base_repo_id,
            head_repo_id=base_repo_id,
            status=status,
            original_status=change.status,
            title=change.subject,
            description=change.subject,
            url=f"{ctx.connection.url}c/{project.name}/+/{change.change_number}",
            author_name=change.owner_name,
            author_id=change.owner_email,
            pull_request_key=change.change_number,
            created_date=change.created_date,
            closed_date=change.closed_date,
            merged_date=change.merged_date,
            type=None,
            component=None,
            base_ref=change.branch,
            merge_commit_sha=merge_commit_sha,
        )
