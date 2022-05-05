## Why advanced mode?

Advanced mode allows users to create any pipeline by writing JSON. This is most useful for users who'd like to:

1. Collect multiple GitHub/GitLab repos or Jira projects within a single pipeline
2. Have fine-grained control over what entities to collect or what subtasks to run for each plugin
3. Orchestrate a complex pipeline that consists of multiple stages of plugins.

Advaned mode gives the most flexiblity to users by exposing the JSON API

## How to use advanced mode to create pipelines?

1. Visit the "Create Pipeline Run" page on `config-ui`

![image](https://user-images.githubusercontent.com/2908155/164569669-698da2f2-47c1-457b-b7da-39dfa7963e09.png)

2. Scroll to the bottom and toggle on the "Advanced Mode" button

![image](https://user-images.githubusercontent.com/2908155/164570039-befb86e2-c400-48fe-8867-da44654194bd.png)

3. The pipeline editor expects a 2D array of plugins. The first dimension represents different stages of the pipeline and the second dimension describes the plugins in each stage. Stages run in sequential order and plugins within the same stage runs in parallel. We provide some templates for users to get started. Please also see the next section for some examples.

![image](https://user-images.githubusercontent.com/2908155/164576122-fc015fea-ca4a-48f2-b2f5-6f1fae1ab73c.png)

## Examples

1. Collect multiple GitLab repos sequentially. 

>When there're multiple collection tasks against a single data source, we recommend running these tasks sequentially since the collection speed is mostly limited by the API rate limit of the data source. 
>Running multiple tasks against the same data source is unlikely to speed up the process and may overwhelm the data source.


Below is an example for collecting 2 GitLab repos sequentially. It has 2 stages, each contains a GitLab task. 


```
[
  [
    {
      "Plugin": "gitlab",
      "Options": {
        "projectId": 15238074
      }
    }
  ],
  [
    {
      "Plugin": "gitlab",
      "Options": {
        "projectId": 11624398
      }
    }
  ]
]
```


2. Collect a GitHub repo and a Jira board in parallel

Below is an example for collecting a GitHub repo and a Jira board in parallel. It has a single stage with a GitHub task and a Jira task. Since users can configure multiple Jira connection, it's required to pass in a `connectionId` for Jira task to specify which connection to use.

```
[
  [
    {
      "Plugin": "github",
      "Options": {
        "repo": "lake",
        "owner": "merico-dev"
      }
    },
    {
      "Plugin": "jira",
      "Options": {
        "connectionId": 1,
        "boardId": 76
      }
    }
  ]
]
```
