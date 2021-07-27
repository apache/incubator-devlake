# Contributing to Lake

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
