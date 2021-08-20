# Plugin Execution Pipeline POC

This is a POC demo to show how plugin pipeline could be done.

# Assumption

## Plugins

1. `JiraPlugin` responsible for jira data collection and enrichment
2. `GitlabPlugin` responsible for gitlab data collection and enrichment
3. `QualityPlugin` is to calculate compound figures base on `JiraPlugin` and `GitlabPlugin`
4. All plugin forms a DAG

## Scenario

1. Want to see some `Quality` figures, so he issue a api request to our api
2. Server generate a pipeline base on Plugin Dependency DAG, which is:
    Pipeline ( all jobs of a step executed in parallel)

    | step   | parallel job 1 | parallel job 2 |
    | ------ | -------------- | -------------- |
    | step 1 | jira           | gitlab         |
    | step 2 | quality        |                |

## Problem

Say if we can generate a pipeline like this:

```js
[
    [
         { plugin: 'jira', args: { boardId: 8, force: false } },
         { plugin: 'gitlab', args: { projectId: 800012, force: false } }
    ],
    [
         { plugin: 'quality', args: { boardId: 8, projectId: 800012, force: false } }
    ]
]
```

Can we leverage `bulljs` to run those tasks in multiple process?

# Solution

1. `worker.ts` is the daemon process listening task queue and execute them, it will spawn 5 `worker-process.ts` in this demo.
2. `worker-process.ts` takes bulljs Job and call specific plugin
3. `main.ts` simulates API Server pipeline execution

# How to run

1. checkout branch of this PR
2. run `npm run compose` for `redis` service
3. run `npx ts-node proposal1/worker.ts` and keep it running
4. run `npx ts-node proposal1/main.ts` to simulate a pipeline execution
5. expected output of `worker.ts`:

```sh
[15772|15][2021-08-20T11:00:22.453Z] <jira> start pipeline job { boardId: 8, force: false }
[15766|16][2021-08-20T11:00:22.453Z] <gitlab> start pipeline job { projectId: 800012, force: false }
[15766|16][2021-08-20T11:00:22.453Z] <gitlab> INFO >>> gitlab plugin start
[15772|15][2021-08-20T11:00:22.453Z] <jira> INFO >>> jira plugin start
[15766|16][2021-08-20T11:00:22.453Z] <gitlab> INFO >>> gitlab start collect repo data
[15772|15][2021-08-20T11:00:22.453Z] <jira> INFO >>> jira start collect board data
[15772|15][2021-08-20T11:00:23.454Z] <jira> progress: 10
[15772|15][2021-08-20T11:00:23.454Z] <jira> INFO >>> jira end collect board data
[15772|15][2021-08-20T11:00:23.454Z] <jira> INFO >>> jira start collect issues data
[15766|16][2021-08-20T11:00:23.454Z] <gitlab> progress: 10
[15766|16][2021-08-20T11:00:23.454Z] <gitlab> INFO >>> gitlab end collect repo data
[15766|16][2021-08-20T11:00:23.454Z] <gitlab> INFO >>> gitlab start collect commits data
[15772|15][2021-08-20T11:00:24.455Z] <jira> progress: 50
[15772|15][2021-08-20T11:00:24.455Z] <jira> INFO >>> jira end collect issues data
[15772|15][2021-08-20T11:00:24.455Z] <jira> INFO >>> jira start enricher issues data
[15766|16][2021-08-20T11:00:24.455Z] <gitlab> progress: 50
[15766|16][2021-08-20T11:00:24.456Z] <gitlab> INFO >>> gitlab end collect commits data
[15766|16][2021-08-20T11:00:24.456Z] <gitlab> INFO >>> gitlab start enricher commits data
[15766|16][2021-08-20T11:00:25.456Z] <gitlab> progress: 100
[15766|16][2021-08-20T11:00:25.456Z] <gitlab> INFO >>> gitlab end enricher commits data
[15766|16][2021-08-20T11:00:25.456Z] <gitlab> INFO >>> gitlab plugin end
[15766|16][2021-08-20T11:00:25.456Z] <gitlab> end pipeline job { projectId: 800012, force: false }
[15772|15][2021-08-20T11:00:25.456Z] <jira> progress: 100
[15772|15][2021-08-20T11:00:25.457Z] <jira> INFO >>> jira end enricher issues data
[15772|15][2021-08-20T11:00:25.457Z] <jira> INFO >>> jira plugin end
[15772|15][2021-08-20T11:00:25.457Z] <jira> end pipeline job { boardId: 8, force: false }
[15772|17][2021-08-20T11:00:25.460Z] <quality> start pipeline job { boardId: 8, projectId: 800012, force: false }
[15772|17][2021-08-20T11:00:25.460Z] <quality> INFO >>> quality plugin start
[15772|17][2021-08-20T11:00:25.460Z] <quality> INFO >>> we can use entities from parent plugins now, like:  function function
[15772|17][2021-08-20T11:00:25.460Z] <quality> INFO >>> quality start calculating BUGS COUNT PER 1K LOC
[15772|17][2021-08-20T11:00:26.461Z] <quality> progress: 50
[15772|17][2021-08-20T11:00:26.461Z] <quality> INFO >>> quality end calculating BUGS COUNT PER 1K LOC
[15772|17][2021-08-20T11:00:26.461Z] <quality> INFO >>> quality start calculating INCIDENTS COUNT PER 1K LOC
[15772|17][2021-08-20T11:00:27.463Z] <quality> progress: 100
[15772|17][2021-08-20T11:00:27.463Z] <quality> INFO >>> quality end calculating INCIDENTS COUNT PER 1K LOC
[15772|17][2021-08-20T11:00:27.463Z] <quality> INFO >>> quality plugin end
[15772|17][2021-08-20T11:00:27.463Z] <quality> end pipeline job { boardId: 8, projectId: 800012, force: false }
```



