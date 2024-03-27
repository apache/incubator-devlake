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
import pytest

from pydevlake.testing.testing import assert_stream_convert
from pydevlake.testing import ContextBuilder
import pydevlake.domain_layer.code as code
from gerrit.main import GerritPlugin


@pytest.fixture
def context():
    return (
        ContextBuilder(GerritPlugin())
        .with_connection(endpoint="https://gerrit.onap.org/r/")
        .with_scope_config()
        .with_scope(name="ccsdk/oran", url="https://gerrit.onap.org/r/ccsdk/oran")
        .build()
    )


@pytest.mark.parametrize(
    "raw, expected",
    [
        (
            {
                "id": "ccsdk%2Foran~master~I1c816846ebc2d459d0619550c6e127735652d076",
                "project": "ccsdk/oran",
                "branch": "master",
                "hashtags": [],
                "change_id": "I1c816846ebc2d459d0619550c6e127735652d076",
                "subject": "Add the Policy Management Service API",
                "status": "MERGED",
                "created": "2020-07-30 13:45:02.000000000",
                "updated": "2020-07-30 16:03:42.000000000",
                "submitted": "2020-07-30 15:58:50.000000000",
                "submitter": {"_account_id": 865},
                "insertions": 842,
                "deletions": 0,
                "total_comment_count": 0,
                "unresolved_comment_count": 0,
                "has_review_started": True,
                "current_revision": "39b0ae8275440fed45ea68bb8941e90a2a5f1d28",
                "submission_id": "110737-1596124730201-3ead5e5d",
                "meta_rev_id": "0a39fc46fb26cd68fd238aa4bdfa21e9f0560c7d",
                "_number": 110737,
                "owner": {
                    "_account_id": 3763,
                    "name": "Henrik Andersson",
                    "email": "henrik.b.andersson@est.tech",
                    "username": "elinuxhenrik",
                },
                "requirements": [],
                "submit_records": [
                    {
                        "status": "CLOSED",
                        "labels": [
                            {
                                "label": "Verified",
                                "status": "OK",
                                "applied_by": {"_account_id": 459},
                            },
                            {
                                "label": "Code-Review",
                                "status": "OK",
                                "applied_by": {"_account_id": 865},
                            },
                            {
                                "label": "Non-Author-Code-Review",
                                "status": "OK",
                                "applied_by": {"_account_id": 865},
                            },
                        ],
                    }
                ],
            },
            code.PullRequest(
                base_repo_id="gerrit:GerritProject:1:s",
                head_repo_id="gerrit:GerritProject:1:s",
                status="MERGED",
                original_status="MERGED",
                title="Add the Policy Management Service API",
                description="Add the Policy Management Service API",
                url="https://gerrit.onap.org/r/c/ccsdk/oran/+/110737",
                author_name="Henrik Andersson",
                author_id="henrik.b.andersson@est.tech",
                pull_request_key=110737,
                created_date=datetime(2020, 7, 30, 13, 45, 2),
                merged_date=datetime(2020, 7, 30, 16, 3, 42),
                merge_commit_sha="39b0ae8275440fed45ea68bb8941e90a2a5f1d28",
                head_ref=None,
                base_ref="master",
                head_commit_sha=None,
                base_commit_sha=None,
            ),
        ),
    ],
)
def test_changes_stream_convert(raw, expected, context):
    assert_stream_convert(GerritPlugin, "gerritchanges", raw, expected, context)


def test_change_commits_stream(context):
    state = {}
    stream = GerritPlugin().get_stream("gerritchanges")
    parent_dict = next(stream.collect(state, context))[0]
    parent = stream.extract(parent_dict)
    stream = GerritPlugin().get_stream("gerritchangecommits")
    for change_commit_data, state in stream.collect(state, context, parent):
        change_commit = stream.extract(change_commit_data)
        stream.convert(change_commit, context)
