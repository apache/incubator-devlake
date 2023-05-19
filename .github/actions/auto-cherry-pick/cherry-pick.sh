#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#!/bin/bash

TARGET_BRANCH="release-${LABEL_NAME##*-}"
PR_BRANCH="auto-cherry-pick-$TARGET_BRANCH-$GITHUB_SHA"

echo "REPOSITORY: $REPOSITORY"
echo "PR Number: $PR_NUMBER"
echo "Label: $LABEL_NAME"
echo "GitHub SHA: $GITHUB_SHA"
echo "Author Email: $AUTHOR_EMAIL"
echo "Author Name: $AUTHOR_NAME"
echo "Assignees: $ASSIGNEES"
echo "Target Branch: $TARGET_BRANCH"
echo "PR Branch: $PR_BRANCH"

git config --global user.email "$AUTHOR_EMAIL"
git config --global user.name "$AUTHOR_NAME"

git remote update
git fetch --all
git restore .
git checkout -b $PR_BRANCH origin/$TARGET_BRANCH
git cherry-pick $GITHUB_SHA
git push origin $PR_BRANCH

gh pr create -B $TARGET_BRANCH -H $PR_BRANCH -t "cherry-pick: $PR_NUMBER" -b "this a auto create pr!<br />cherry pick from https://github.com/$REPOSITORY/pull/$PR_NUMBER" -a $ASSIGNEES
gh pr comment $PR_NUMBER --body "🤖 cherry pick finished successfully 🎉!"
