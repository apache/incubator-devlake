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


from typing import Optional

from pydantic import BaseModel

from pydevlake.message import PipelineTask


def gitextractor(url: str, repo_id: str, proxy: Optional[str]):
    return PipelineTask(
        plugin="gitextractor",
        options={
            "url": url,
            "repoId": repo_id,
            "proxy": proxy
        },
    )


class RefDiffOptions(BaseModel):
    tags_limit: Optional[int] = 10
    tags_order: Optional[str] = "reverse semver"
    tags_pattern: Optional[str] = r"/v\d+\.\d+(\.\d+(-rc)*\d*)*$/"


def refdiff(repo_id: str, options: RefDiffOptions=None):
    if options is None:
        options = RefDiffOptions()
    return PipelineTask(
        plugin="refdiff",
        options={
            "repoId":repo_id,
            "tagsLimit": options.tags_limit,
            "tagsOrder": options.tags_order,
            "tagsPattern": options.tags_pattern
        },
    )