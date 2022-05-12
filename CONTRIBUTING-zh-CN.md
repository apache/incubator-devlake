# 为 Dev Lake 做贡献

👍🎉 首先，感谢你付出的时间！ 🎉👍

以下是一套为 Dev Lake 做贡献的准则。这些主要是指导方针，而非规则。又任何修改建议，请根据个人判断提 PR 即可。


## 我如何贡献？

1. 通过填写所需的 issue 模板并将新问题标记为 "bug" 来报告bug

2. 建议增强功能

如果你打算更改公共 API，或者对实现进行任何非微不足道的改变，我们建议提交一个 issue。这可以确保在你投入大量精力之前，我们已经就方案达成一致。

如果你只是修复一个bug，马上提交一个 pull request 也是可以的，但我们仍然建议提交一个 issue，详细说明你要修复的内容。这对于我们“不接受一个特定的修复，但又想跟踪这个问题”的情况下是很有帮助的。


## 维护者团队 @ Merico

Dev Lake 由 Merico 的一群工程师维护，由 [@hezyin](https://github.com/hezyin) 领导。我们的目标是实现 24 小时内回复问题的 SLA

## 风格指南

### Git Commit Message

我们遵循此规范：[conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)

#### Commit 工具

我们使用 https://github.com/lintingzhen/commitizen-go 来提交Commit

```sh
make commit
```

```
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
