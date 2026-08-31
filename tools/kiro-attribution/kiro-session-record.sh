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
# Records the Kiro session id that last touched this repository.
#
# Registered as a Kiro Stop hook, which fires after the agent finishes a turn -
# by then the files are written and a commit is the likely next step. Recording
# at turn start instead would capture a session that went on to touch other
# files, or one the user abandoned.
#
# The payload arrives as JSON on stdin and carries session_id; see
# https://kiro.dev/docs/cli/hooks. Verified against kiro-cli 2.6.1.
set -eu

GIT_DIR=$(git rev-parse --git-dir 2>/dev/null) || exit 0
PAYLOAD=$(cat)

# Extract session_id without depending on jq being installed.
SESSION_ID=$(printf '%s' "$PAYLOAD" \
  | sed -n 's/.*"session_id"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
  | head -1)

[ -n "$SESSION_ID" ] || exit 0

# Append rather than overwrite: several sessions may contribute to one commit,
# and dropping the earlier ones would under-report AI involvement.
#
# Epoch seconds first, because that is what the commit hook compares against.
# An ISO timestamp would have to agree on a timezone with git's own output, and
# a mismatch there fails silently - the trailer simply never appears. The
# readable form is kept as a second column for humans debugging the file.
printf '%s\t%s\t%s\n' "$(date -u +%s)" "$SESSION_ID" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  >> "$GIT_DIR/kiro-sessions"

# Keep only recent entries so the file cannot grow without bound. 200 turns is
# far more than any single commit spans.
if [ "$(wc -l < "$GIT_DIR/kiro-sessions")" -gt 200 ]; then
  tail -100 "$GIT_DIR/kiro-sessions" > "$GIT_DIR/kiro-sessions.tmp"
  mv "$GIT_DIR/kiro-sessions.tmp" "$GIT_DIR/kiro-sessions"
fi

exit 0
