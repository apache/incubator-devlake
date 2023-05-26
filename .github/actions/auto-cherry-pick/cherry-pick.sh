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

set -e

TARGET_BRANCH="release-${LABEL_NAME##*-}"
PR_BRANCH="auto-cherry-pick-$TARGET_BRANCH-$GITHUB_SHA"
AUTO_CHERRY_PICK_LABEL="bot/auto-cherry-pick"
AUTO_CHERRY_PICK_VERSION_LABEL="bot/auto-cherry-pick-for-$TARGET_BRANCH"
AUTO_CHERRY_PICK_FAILED_LABEL="bot/auto-cherry-pick-failed"
AUTO_CHERRY_PICK_COMPLETED_LABEL="bot/auto-cherry-pick-completed"

echo "==================== Basic Info ===================="
echo "PR Number: $PR_NUMBER"
echo "PR Title: $PR_TITLE"
echo "PR Body: $PR_BODY"
echo "Label: $LABEL_NAME"
echo "GitHub SHA: $GITHUB_SHA"
echo "Author Email: $AUTHOR_EMAIL"
echo "Author Name: $AUTHOR_NAME"
echo "Assignees: $ASSIGNEES"
echo "Target Branch: $TARGET_BRANCH"
echo "PR Branch: $PR_BRANCH"

echo "==================== Git Cherry Pick ===================="
git config --global user.email "$AUTHOR_EMAIL"
git config --global user.name "$AUTHOR_NAME"

git remote update
git fetch --all
git restore .
git checkout -b $PR_BRANCH origin/$TARGET_BRANCH
git cherry-pick -m 1 --strategy=recursive --strategy-option=theirs $GITHUB_SHA || (
	gh pr comment $PR_NUMBER --body "🤖 The current file has a conflict, and the pr cannot be automatically created."
	gh pr edit $PR_NUMBER --add-label $AUTO_CHERRY_PICK_FAILED_LABEL || (
		gh label create $AUTO_CHERRY_PICK_FAILED_LABEL -c "#D93F0B" -d "auto cherry pick failed"
		gh pr edit $PR_NUMBEr --add-label $AUTO_CHERRY_PICK_FAILED_LABEL
	)
	exit 1
)
git push origin $PR_BRANCH

echo "==================== GitHub Auto Create PR ===================="
AUTO_CREATED_PR_LINK=$(gh pr create \
	-B $TARGET_BRANCH \
	-H $PR_BRANCH \
	-t "$PR_TITLE (cherry-picked-from #$PR_NUMBER)" \
	-b "$PR_BODY" \
	-a $ASSIGNEES)

gh pr comment $PR_NUMBER --body "🤖 cherry pick finished successfully 🎉!"
gh pr edit $PR_NUMBER --add-label $AUTO_CHERRY_PICK_COMPLETED_LABEL || (
	gh label create $AUTO_CHERRY_PICK_COMPLETED_LABEL -c "#0E8A16" -d "auto cherry pick completed"
	gh pr edit $PR_NUMBER --add-label $AUTO_CHERRY_PICK_COMPLETED_LABEL
)

gh pr comment $AUTO_CREATED_PR_LINK --body "🤖 this a auto create pr!cherry picked from #$PR_NUMBER."
gh pr edit $AUTO_CREATED_PR_LINK --add-label "$AUTO_CHERRY_PICK_LABEL,$AUTO_CHERRY_PICK_VERSION_LABEL" || (
	gh label create $AUTO_CHERRY_PICK_LABEL -c "#5319E7" -d "auto cherry pick pr"
	gh label create $AUTO_CHERRY_PICK_VERSION_LABEL -c "#5319E7" -d "auto cherry pick pr for $TARGET_BRANCH"
	gh pr edit $AUTO_CREATED_PR_LINK --add-label "$AUTO_CHERRY_PICK_LABEL,$AUTO_CHERRY_PICK_VERSION_LABEL"
)
