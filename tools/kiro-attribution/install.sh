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
# Installs Kiro commit attribution into a repository.
#
# Usage:
#   ./install.sh                 install into the current repository
#   ./install.sh /path/to/repo   install into another repository
#   ./install.sh --check         report status without changing anything
#   ./install.sh --uninstall     remove what this script installed
#
# Where the hook goes depends on core.hooksPath, which git consults in place of
# .git/hooks whenever it is set. Two tools set it in practice and they behave
# differently, so the destination is decided by testing rather than by name:
#
#   husky (repository scope, e.g. config-ui/.husky) does NOT forward to
#   .git/hooks - its shim only re-executes itself. Verified by installing into
#   .git/hooks in such a repo and watching the trailer never appear. So the hook
#   must go into the husky directory itself.
#
#   git-defender (system scope, a root-owned directory) DOES invoke
#   .git/hooks/prepare-commit-msg - its own error text says "Your local
#   prepare-commit-msg hook failed". Its directory is not writable anyway, so
#   .git/hooks is both the only option and the correct one.
#
# The rule that covers both: install into the directory git reads, unless that
# directory is unwritable, in which case fall back to .git/hooks - which is
# exactly the case where the tooling forwards there.
#
# core.hooksPath itself is never modified. Redirecting it would be worse than
# useless under git-defender, which reports the event ("User had a non Code
# Defender hooks path value set"), tripping a security signal while solving
# nothing.
set -eu

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
MODE=install
TARGET_REPO=.
MARKER='# kiro-attribution'

for arg in "$@"; do
  case "$arg" in
    --check)     MODE=check ;;
    --uninstall) MODE=uninstall ;;
    --help|-h)   sed -n '2,8p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    -*)          echo "unknown option: $arg" >&2; exit 2 ;;
    *)           TARGET_REPO=$arg ;;
  esac
done

cd "$TARGET_REPO" 2>/dev/null || { echo "not a directory: $TARGET_REPO" >&2; exit 1; }
# --absolute-git-dir, not --git-dir: the latter returns a path relative to the
# current directory (".git" at the repo root), and the verification step below
# cd's elsewhere. Every path derived from a relative GIT_DIR would then resolve
# against the wrong directory - silently, since sed on a missing file just
# produces an empty hook that chmod happily marks executable.
GIT_DIR=$(git rev-parse --absolute-git-dir 2>/dev/null) || {
  echo "not a git repository: $TARGET_REPO" >&2; exit 1
}
REPO_ROOT=$(git rev-parse --show-toplevel)
SESSION_FILE="$GIT_DIR/kiro-sessions"

# resolve_hook_dir prints the directory to install into, plus a short reason.
resolve_hook_dir() {
  _hp=$(git config --get core.hooksPath 2>/dev/null || true)
  if [ -z "$_hp" ]; then
    echo "$GIT_DIR/hooks|git default"
    return
  fi
  case "$_hp" in
    /*) _abs=$_hp ;;
    *)  _abs="$REPO_ROOT/$_hp" ;;
  esac
  # Writability is the deciding test. An unwritable hooksPath belongs to managed
  # tooling, and that tooling is the kind that forwards to .git/hooks.
  if [ -d "$_abs" ] && [ -w "$_abs" ]; then
    echo "$_abs|core.hooksPath ($_hp)"
  elif [ ! -d "$_abs" ]; then
    echo "$_abs|core.hooksPath ($_hp), creating it"
  else
    echo "$GIT_DIR/hooks|core.hooksPath ($_hp) is not writable, falling back"
  fi
}

RESOLVED_DIR=$(resolve_hook_dir)
HOOK_DIR=${RESOLVED_DIR%|*}
HOOK_WHY=${RESOLVED_DIR#*|}
HOOK_FILE="$HOOK_DIR/prepare-commit-msg"

# ---------------------------------------------------------------- check --------

if [ "$MODE" = check ]; then
  echo "repository:  $REPO_ROOT"

  echo "hook dir:    $HOOK_DIR"
  echo "             chosen by: $HOOK_WHY"

  if [ -f "$HOOK_FILE" ] && grep -q "$MARKER" "$HOOK_FILE" 2>/dev/null; then
    if [ -x "$HOOK_FILE" ]; then
      echo "commit hook: installed"
    else
      echo "commit hook: installed but NOT EXECUTABLE - git will skip it"
    fi
  elif [ -f "$HOOK_FILE" ]; then
    echo "commit hook: a different prepare-commit-msg is present (not ours)"
  else
    echo "commit hook: not installed"
  fi

  KIRO_HOOK="$REPO_ROOT/.kiro/hooks/kiro-attribution.json"
  if [ -f "$KIRO_HOOK" ]; then
    echo "kiro hook:   registered"
  else
    echo "kiro hook:   not registered"
  fi

  if [ -f "$SESSION_FILE" ]; then
    echo "session log: $(wc -l < "$SESSION_FILE" | tr -d ' ') entries"
  else
    echo "session log: none yet (the Kiro hook has not fired in this repo)"
  fi

  # Coverage decides whether attribution can support a conclusion at all.
  # Unknown coverage means unusable data, not partial data: a commit with no
  # trailer is indistinguishable between "no AI was used" and "the hook was
  # never installed on that machine".
  TOTAL=$(git rev-list --count HEAD 2>/dev/null || echo 0)
  TAGGED=$(git log --all --grep='^Kiro-Session-Id:' --oneline 2>/dev/null | wc -l | tr -d ' ')
  echo "commits:     $TOTAL total, $TAGGED with a Kiro-Session-Id trailer"
  exit 0
fi

# ------------------------------------------------------------ uninstall --------

if [ "$MODE" = uninstall ]; then
  if [ -f "$HOOK_FILE" ] && grep -q "$MARKER" "$HOOK_FILE" 2>/dev/null; then
    if [ -f "$HOOK_FILE.kiro-backup" ]; then
      mv "$HOOK_FILE.kiro-backup" "$HOOK_FILE"
      echo "restored the pre-existing prepare-commit-msg"
    else
      rm -f "$HOOK_FILE"
      echo "removed $HOOK_FILE"
    fi
  else
    echo "no hook installed by this script"
  fi
  # Also clear a stale copy in .git/hooks. A repo can acquire core.hooksPath
  # after the hook was installed (adding husky does exactly that), leaving an
  # orphan behind that --check would no longer look at.
  STALE="$GIT_DIR/hooks/prepare-commit-msg"
  if [ "$STALE" != "$HOOK_FILE" ] && [ -f "$STALE" ] && grep -q "$MARKER" "$STALE" 2>/dev/null; then
    rm -f "$STALE"
    echo "removed a stale copy at $STALE"
  fi
  rm -f "$REPO_ROOT/.kiro/hooks/kiro-attribution.json" 2>/dev/null || true
  echo "note: $SESSION_FILE is left in place; delete it manually if unwanted"
  exit 0
fi

# -------------------------------------------------------------- install --------

mkdir -p "$HOOK_DIR"

# Chain to an unrelated hook rather than replacing it: silently dropping another
# team's hook would break their workflow with nothing to point at.
CHAIN_TO=""
if [ -f "$HOOK_FILE" ] && ! grep -q "$MARKER" "$HOOK_FILE" 2>/dev/null; then
  mv "$HOOK_FILE" "$HOOK_FILE.kiro-backup"
  CHAIN_TO="$HOOK_FILE.kiro-backup"
  echo "kept the existing prepare-commit-msg as prepare-commit-msg.kiro-backup"
fi

{
  printf '#!/bin/sh\n'
  printf '%s - appends a Kiro-Session-Id trailer.\n' "$MARKER"
  printf '# Installed by tools/kiro-attribution/install.sh; safe to remove.\n'
  printf 'set -eu\n\n'
  if [ -n "$CHAIN_TO" ]; then
    printf '# Run the pre-existing hook first. Its failure must still abort the\n'
    printf '# commit, hence propagating the exit code.\n'
    printf 'if [ -x "%s" ]; then "%s" "$@" || exit $?; fi\n\n' "$CHAIN_TO" "$CHAIN_TO"
  fi
  # Inlined rather than sourced: a developer's checkout of this repository can
  # move or be deleted, and a hook pointing at a missing file would break every
  # commit in that repo.
  sed '1d' "$SCRIPT_DIR/prepare-commit-msg.sh"
} > "$HOOK_FILE"

chmod +x "$HOOK_FILE"
echo "installed $HOOK_FILE"
echo "  location chosen by: $HOOK_WHY"
[ -n "$CHAIN_TO" ] && echo "  chaining to $CHAIN_TO"

# Register the Kiro-side recorder for this workspace.
#
# The schema matters: Kiro loads only .kiro/hooks/*.json files shaped as
# {"version":"v1","hooks":[...]}. The older one-object-per-file format with a
# .kiro.hook extension is silently ignored - the log says "loaded 0 standalone
# hooks" and nothing else indicates why. Field names changed too: "when.type":
# "agentStop" became "trigger": "Stop", and "then" became "action".
#
# Hooks load at session start, so a running Kiro session must be restarted before
# a newly installed hook fires.
if [ -d "$REPO_ROOT/.kiro" ]; then
  mkdir -p "$REPO_ROOT/.kiro/hooks"
  # Invoked via bash, and with a repo-relative path.
  #
  # The bash prefix is required by Kiro's permission model: shell capability
  # rules match on command prefix, and a bare script path matches nothing, so the
  # hook is skipped with no error logged. "bash *" is a standard allowed pattern.
  #
  # Relative rather than absolute so the hook survives the repo being moved or
  # re-cloned; Kiro runs hooks with the workspace root as cwd.
  REL_RECORDER=$(printf '%s' "$SCRIPT_DIR/kiro-session-record.sh" | sed "s|^$REPO_ROOT/||")
  cat > "$REPO_ROOT/.kiro/hooks/kiro-attribution.json" <<EOF
{
  "version": "v1",
  "hooks": [
    {
      "name": "kiro-session-record",
      "description": "Records the Kiro session id so the next commit carries a Kiro-Session-Id trailer. Installed by tools/kiro-attribution/install.sh.",
      "trigger": "Stop",
      "action": {
        "type": "command",
        "command": "bash $REL_RECORDER"
      },
      "timeout": 15,
      "enabled": true
    }
  ]
}
EOF
  echo "registered the Kiro Stop hook in $REPO_ROOT/.kiro/hooks/kiro-attribution.json"
  echo "  RESTART Kiro for it to take effect - hooks load at session start"
else
  echo "no .kiro/ directory here, so the Kiro-side hook was not registered."
  echo "  Create .kiro/hooks/kiro-attribution.json with trigger \"Stop\" running:"
  echo "    $SCRIPT_DIR/kiro-session-record.sh"
fi

# ---------------------------------------------------------------- verify -------
# Proved rather than assumed. Every failure mode in this setup is silent: a
# missing executable bit, a timezone mismatch between the recorder and the hook,
# an off-by-one at the commit boundary. All three produce "no trailer, no
# error", which is indistinguishable from "no AI was used".
echo
echo "verifying in a scratch repository..."
VERIFY_DIR=$(mktemp -d)
trap 'rm -rf "$VERIFY_DIR"' EXIT
(
  cd "$VERIFY_DIR"
  git init -q
  git config user.email verify@example.com
  git config user.name verify
  mkdir -p .git/hooks
  # Point core.hooksPath at this scratch repo's own directory. Repository scope
  # beats system scope, so any org-wide tooling (git-defender) stays out of the
  # way while our hook still runs.
  #
  # --no-verify would be the obvious alternative and is wrong twice over: it
  # skips every hook including the one under test, and on a managed laptop it
  # records a security-scan bypass.
  git config core.hooksPath .git/hooks
  # Drop the chain line: the pre-existing hook belongs to the real repository and
  # may well refuse to run here.
  sed '/kiro-backup/d' "$HOOK_FILE" > .git/hooks/prepare-commit-msg
  chmod +x .git/hooks/prepare-commit-msg

  echo x > x && git add . && git commit -q -m baseline
  sleep 1
  printf '%s\tverify-session-id\t%s\n' "$(date -u +%s)" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    >> .git/kiro-sessions
  echo y > y && git add . && git commit -q -m "with session"

  if ! git log -1 --format=%B | grep -q 'Kiro-Session-Id: verify-session-id'; then
    echo "  FAILED: no trailer written. Do not rely on this install." >&2
    exit 1
  fi
  echo "  a trailer is written when a session was recorded"

  echo z > z && git add . && git commit -q -m "no new session"
  if git log -1 --format=%B | grep -q 'Kiro-Session-Id'; then
    echo "  FAILED: a commit with no new session inherited a trailer." >&2
    exit 1
  fi
  echo "  no trailer when no new session was recorded"
)
echo
echo "done. Check status any time with:"
echo "  $SCRIPT_DIR/install.sh --check"
