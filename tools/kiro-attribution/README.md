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
# Binding Kiro sessions to git commits

## What this solves

Usage data alone cannot answer whether AI helped a given person: high credit
consumption is compatible with both "the agent did a lot of work" and "the agent
went in circles". Answering it needs to know which commits the AI actually
contributed to.

Kiro has no git-lifecycle hook (the feature request, kirodotdev/Kiro#6436, is
still open), and its S3 interaction logs carry no `conversationId`, so neither
source can bind a session to a commit on its own.

## Why a trailer rather than a time window

The obvious alternative - "mark a commit as AI-assisted if Kiro was active
within 30 minutes" - has a base-rate trap. For anyone who keeps Kiro open all
day, every commit qualifies. The attribution rate approaches 100%, and the label
stops distinguishing between people, which is exactly what it was for. Heavy
users are the population the measurement most needs to separate.

A trailer records what happened instead of inferring it. That turns attribution
from probabilistic into deterministic, which in turn unlocks the strongest
comparison available from observational data: **the same person, in the same
week, on AI-assisted versus unassisted work**. That within-person comparison
holds individual skill, team environment and time period constant, so the
remaining difference can only come from AI involvement.

Trailers are also a standard git mechanism - `Co-Authored-By:` and
`Signed-off-by:` use the same convention - so nothing here is bespoke, and the
value stays readable to anyone inspecting the log.

## How it works

```
Kiro Stop hook  ──►  .git/kiro-sessions        (epoch, session id, ISO time)
                            │
git prepare-commit-msg  ────┘
                            ▼
                     commit message:
                       fix: payment retry

                       Kiro-Session-Id: sess-abc-123
                            │
                            ▼
                     DevLake already collects commits.message,
                     so extraction needs no new collector.
```

## Installation

Run the installer once per repository:

```sh
tools/kiro-attribution/install.sh
```

It installs the git `prepare-commit-msg` hook in the directory Git actually
uses (respecting a writable `core.hooksPath`, or falling back to `.git/hooks`
for managed tooling that forwards there). It does not change `core.hooksPath`.
If the repository has a `.kiro/` directory, it also writes the current Kiro
hook schema to `.kiro/hooks/kiro-attribution.json`:

```json
{
  "version": "v1",
  "hooks": [
    {
      "name": "kiro-session-record",
      "description": "Records the Kiro session id so the next commit carries a Kiro-Session-Id trailer.",
      "trigger": "Stop",
      "action": {
        "type": "command",
        "command": "bash tools/kiro-attribution/kiro-session-record.sh"
      },
      "timeout": 15,
      "enabled": true
    }
  ]
}
```

The `bash` prefix is required by Kiro's command permission rules. Restart Kiro
after installation because hooks are loaded when a session starts.

Check or remove the installation with:

```sh
tools/kiro-attribution/install.sh --check
tools/kiro-attribution/install.sh --uninstall
```

## Verified behaviour

| Scenario | Result |
|---|---|
| No AI session before the commit | No trailer |
| One session | One trailer line |
| Several sessions since the last commit | One line each, deduplicated |
| A commit with no new session | No trailer, and does **not** inherit the previous commit's |
| `--amend` | Trailer not duplicated |

Two defects were found and fixed while testing these, both of which failed
silently rather than erroring:

- **Timezone mismatch.** The recorder wrote UTC while the commit hook compared
  against git's local-time output, an 8-hour gap here. The trailer simply never
  appeared. Both sides now use epoch seconds, which has no timezone to disagree
  about.
- **Boundary off by one.** With `>=`, a session recorded in the same second as a
  commit was counted again by the *next* commit, so a commit containing no AI
  work inherited the previous trailer. Now strictly `>`.

## Known limitations

**Client-side, therefore optional.** The hooks live on the developer's machine
and can be removed or disabled. That reintroduces self-selection: "no AI usage"
and "no hook installed" become indistinguishable. **Hook coverage must itself be
tracked as a metric**, and attribution figures from an unknown-coverage
population must not enter a conclusion.

The S3 exports are the complement here: organization-controlled, tamper-proof,
and complete. Use S3 for the denominator (who uses Kiro, how much) and trailers
for attribution (which commits it touched).

**No backfill.** Only commits made after installation carry a trailer. Team-level
analysis over history has to rest on the S3 usage data.

**IDE session id unverified.** The CLI hook payload carries `session_id` (Kiro's
own docs, and confirmed by the existing adapter in `.kiro/hooks/`). The IDE
documentation lists only `USER_PROMPT`, so the IDE path needs a live probe
before it can be relied on.
