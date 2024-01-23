<!--
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
-->
# lake-builder

Golang builder image for DevLake. Used by GitHub workflows.

https://hub.docker.com/r/mericodev/lake-builder

## Manual Release

```shell
export VERSION=0.0.11
docker build -t mericodev/lake-builder:$VERSION .
docker push mericodev/lake-builder:$VERSION
```

## Tagged Release
1. Create a tag matching the pattern `builder-v#.#.#`, e.g. `builder-v0.0.11`. Determine the previous tag first so you version
it correctly. Example command: `git tag builder-v0.0.11`
2. Push the tag to origin. Example command: `git push origin --tag builder-v0.0.11`
3. Done! This will trigger a GitHub workflow that will push this image to the Docker registry.