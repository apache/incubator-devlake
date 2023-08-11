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

ORIGIN_PR_LABELS=($(echo "$ORIGIN_PR_LABELS_JSON" | jq -r '.[]'))

echo "::group::Origin Info"
echo "Origin PR Number: $ORIGIN_PR_NUMBER"
echo "Origin PR Title: $ORIGIN_PR_TITLE"
echo "Origin PR Labels: $ORIGIN_PR_LABELS"
echo "GitHub SHA: $GITHUB_SHA"
echo "Author Email: $AUTHOR_EMAIL"
echo "Author Name: $AUTHOR_NAME"
echo "Assignees: $ASSIGNEES"
echo "::endgroup::"

TARGET_LABEL_PREFIX="needs-cherrypick-"
TARGET_LABEL=""

for label in "${ORIGIN_PR_LABELS[@]}"; do
	if [[ "$label" == "$TARGET_LABEL_PREFIX"* ]]; then
		TARGET_LABEL="$label"
		break
	fi
done

if [ -z "$TARGET_LABEL" ]; then
	echo "No need cherry-pick."
	exit 0
fi

TARGET_BRANCH="release-${TARGET_LABEL##*-}"
AUTO_CREATE_PR_BRANCH="$TARGET_BRANCH-auto-cherry-pick-$ORIGIN_PR_NUMBER"
AUTO_CHERRY_PICK_LABEL="bot/auto-cherry-pick"
AUTO_CHERRY_PICK_VERSION_LABEL="bot/auto-cherry-pick-for-$TARGET_BRANCH"
AUTO_CHERRY_PICK_FAILED_LABEL="bot/auto-cherry-pick-failed"
AUTO_CHERRY_PICK_COMPLETED_LABEL="bot/auto-cherry-pick-completed"

echo "::group::Generate Info"
echo "Target Branch: $TARGET_BRANCH"
echo "Auto Create PR Branch: $AUTO_CREATE_PR_BRANCH"
echo "Auto Cherry Pick Label: $AUTO_CHERRY_PICK_LABEL"
echo "Auto Cherry Pick Version Label: $AUTO_CHERRY_PICK_VERSION_LABEL"
echo "Auto Cherry Pick Failed Label: $AUTO_CHERRY_PICK_FAILED_LABEL"
echo "Auto Cherry Pick Completed Label: $AUTO_CHERRY_PICK_COMPLETED_LABEL"
echo "::endgroup::"

echo "::group::Git Cherry Pick"
git config --global user.email "$AUTHOR_EMAIL"
git config --global user.name "$AUTHOR_NAME"

git remote update
git fetch --all
git checkout -b $AUTO_CREATE_PR_BRANCH origin/$TARGET_BRANCH
git cherry-pick -m 1 $GITHUB_SHA || (
	gh pr comment $ORIGIN_PR_NUMBER --body "🤖 The current file has a conflict, and the pr cannot be automatically created."
	gh pr edit $ORIGIN_PR_NUMBER --add-label $AUTO_CHERRY_PICK_FAILED_LABEL
	exit 1
)
git push origin $AUTO_CREATE_PR_BRANCH
echo "::endgroup::"

echo "::group::GitHub Auto Create PR"
AUTO_CREATED_PR_LINK=$(gh pr create \
	-B $TARGET_BRANCH \
	-H $AUTO_CREATE_PR_BRANCH \
	-t "cherry-pick #$ORIGIN_PR_NUMBER $ORIGIN_PR_TITLE" \
	-b "cherry-pick #$ORIGIN_PR_NUMBER $ORIGIN_PR_TITLE" \
	-a $ASSIGNEES)

gh pr comment $ORIGIN_PR_NUMBER --body "🤖 cherry pick finished successfully 🎉!"
gh pr edit $ORIGIN_PR_NUMBER --add-label $AUTO_CHERRY_PICK_COMPLETED_LABEL || (
	gh label create $AUTO_CHERRY_PICK_COMPLETED_LABEL -c "#0E8A16" -d "auto cherry pick completed"
	gh pr edit $ORIGIN_PR_NUMBER --add-label $AUTO_CHERRY_PICK_COMPLETED_LABEL
)

gh pr comment $AUTO_CREATED_PR_LINK --body "🤖 this a auto create pr!cherry picked from #$ORIGIN_PR_NUMBER."
gh pr edit $AUTO_CREATED_PR_LINK --add-label "$AUTO_CHERRY_PICK_LABEL,$AUTO_CHERRY_PICK_VERSION_LABEL" || (
	gh label create $AUTO_CHERRY_PICK_VERSION_LABEL -c "#5319E7" -d "auto cherry pick pr for $TARGET_BRANCH"
	gh pr edit $AUTO_CREATED_PR_LINK --add-label "$AUTO_CHERRY_PICK_LABEL,$AUTO_CHERRY_PICK_VERSION_LABEL"
)
echo "::endgroup::"
