#!/bin/sh
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
# Appends a Kiro-Session-Id trailer for sessions that touched this repo since
# the last commit.
#
# A git trailer is used rather than a time-window heuristic because the two
# differ in kind, not degree. Matching "was Kiro active within 30 minutes of
# this commit" hits a base-rate trap: for anyone who keeps Kiro open all day,
# every commit qualifies, the attribution rate approaches 100%, and the label
# can no longer distinguish between people - which is precisely what heavy
# users need it to do. A trailer records what actually happened.
#
# Trailers are a standard git mechanism (Co-Authored-By, Signed-off-by use the
# same convention), so nothing here is bespoke, and the value stays readable to
# a human inspecting the log.
set -eu

MSG_FILE=$1
COMMIT_SOURCE=${2:-}

# Leave amends, merges and squashes alone: their message already carries
# whatever trailers belong to it, and appending again would duplicate them.
case "$COMMIT_SOURCE" in
  merge|squash|commit) exit 0 ;;
esac

GIT_DIR=$(git rev-parse --git-dir 2>/dev/null) || exit 0
SESSION_FILE="$GIT_DIR/kiro-sessions"
[ -f "$SESSION_FILE" ] || exit 0

# Only sessions newer than the previous commit belong to this one. Without this
# cut, every future commit would inherit the whole history of session ids.
#
# Both sides are epoch seconds (%ct on git's side, date +%s on the recorder's).
# Comparing formatted timestamps instead requires the two to agree on a
# timezone, and when they disagree the hook fails silently - no trailer, no
# error. Epoch has no timezone to get wrong.
LAST_COMMIT_EPOCH=$(git log -1 --format=%ct 2>/dev/null || echo 0)
[ -n "$LAST_COMMIT_EPOCH" ] || LAST_COMMIT_EPOCH=0

# Strictly greater than, not >=: a session recorded in the same second as the
# previous commit already belongs to that commit. With >=, it would be counted
# again here, so a commit containing no AI work at all would inherit the
# previous one's trailer.
SESSIONS=$(awk -v since="$LAST_COMMIT_EPOCH" '$1 + 0 > since + 0 { print $2 }' \
  "$SESSION_FILE" 2>/dev/null | sort -u)

[ -n "$SESSIONS" ] || exit 0

# Skip if a trailer is already present, so re-running the hook is harmless.
if grep -q '^Kiro-Session-Id:' "$MSG_FILE" 2>/dev/null; then
  exit 0
fi

# A blank line before trailers is required for git to parse them as such.
printf '\n' >> "$MSG_FILE"
for s in $SESSIONS; do
  printf 'Kiro-Session-Id: %s\n' "$s" >> "$MSG_FILE"
done

exit 0
