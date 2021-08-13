# Contributing to Lake

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

The following is a set of guidelines for contributing to Lake. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.*

## Collaboration Guidelines

1. If you'd like to start working on a new feature, bug, or change of any kind, please open a new Issue in GitHub with the "proposal" label before you begin work. Describe your proposal in great detail. You may use code samples. Invite the related colleagues to discuss and iterate on your proposal. Once it's approved, you can start implementing it and remove the proposal label, and add "in progress" label. This is our mechanism for "acquiring the lock".

2. The maintainer team aims to achieve an SLA of 24 hrs for replying to issues. A timely reply to Issues will encourage our colleagues to communicate via high-quality Issues.

## Proposing a Change

If you intend to change the public API, or make any non-trivial changes to the implementation, we recommend filing an issue. This lets us reach an agreement on your proposal before you put significant effort into it.

If you’re only fixing a bug, it’s fine to submit a pull request right away but we still recommend to file an issue detailing what you’re fixing. This is helpful in case we don’t accept that specific fix but want to keep track of the issue.**

## Approval of Issues
- Read, understand, and discuss Issues and get consensus from contributors
- Assign the Issue to @hezheng or mention him in a comment
- An Issue is approved when Hezheng comments on it as "approved"

## Closing Issues
- Initiate discussion by asking questions before closing issues
- Wait for 5 days of inactivity before closing any issues

## Style guides

### Git Commit message

We follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) guidelines.

#### Commit tool

It's recommended to use [commitizen](https://www.npmjs.com/package/commitizen) to generate "conventional commit message".

```shell
# npm run commit or npx cz
$ npm run commit

> lake@1.0.0 commit /home/code/merico-dev/lake
> cz

cz-cli@4.2.4, cz-conventional-changelog@3.3.0

? Select the type of change that you're committing: (Use arrow keys)
> feat:     A new feature
  fix:      A bug fix
  docs:     Documentation only changes
  style:    Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
  refactor: A code change that neither fixes a bug nor adds a feature
  perf:     A code change that improves performance
  test:     Adding missing tests or correcting existing tests
(Move up and down to reveal more choices)
? What is the scope of this change (e.g. component or file name): (press enter to skip)
? Write a short, imperative tense description of the change (max 93 chars):
 (23) add commit message tool
? Provide a longer description of the change: (press enter to skip)

? Are there any breaking changes? No
? Does this change affect any open issues? No
[chore/commit_message dc34f57] chore: add commit message tool
 5 files changed, 585 insertions(+), 4 deletions(-)
```

--Attribution: 

*quoted from Atom documentation
**quoted from React documentation