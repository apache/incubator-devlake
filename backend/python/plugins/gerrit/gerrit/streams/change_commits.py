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

from typing import Iterable

from pydevlake import Substream, DomainType
import pydevlake.domain_layer.code as code
from gerrit.streams.changes import GerritChanges
from gerrit.models import GerritChange, GerritChangeCommit, GerritProject


class GerritChangeCommits(Substream):
    tool_model = GerritChangeCommit
    domain_types = [DomainType.CODE]
    parent_stream = GerritChanges

    def should_run_on(self, scope: GerritProject) -> bool:
        return True

    def collect(self, state, context, parent: GerritChange) -> Iterable[tuple[object, dict]]:
        # project: GerritProject = context.scope
        if parent.status == "MERGED":
            for commit_id, commit_data in parent.revisions.items():
                data = {"commit_id": commit_id, "pull_request_id": parent.domain_id()}
                data.update(commit_data)
                yield data, state

    def convert(self, commit: GerritChangeCommit, context) -> Iterable[code.PullRequestCommit]:
        yield code.PullRequestCommit(
            commit_sha=commit.commit_id,
            pull_request_id=commit.pull_request_id,
            commit_author_name=commit.author_name,
            commit_author_email=commit.author_email,
            commit_authored_date=commit.author_date,
        )
