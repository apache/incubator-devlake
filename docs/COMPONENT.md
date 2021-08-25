## Rest api
Rest api take responsibility to interact with lake core. We added simplest apis to drive lake work in the first stage. 

## Lake core


## Plugin
Plugin represent an executable components which could be triggered by lake core. Plugins may collect data, calculate data and write data to storage.  
- For example, Jira plugin may collect jira issue data, calculate lead-time and finally write the result to mysql.  
Plugin should be registered before lake starting. Plugin name is now the identity field for plugin.

## Source
Source represent a group of context for specific plugin. We will implement source in the **next** version.  
- For example, we may add two sources for Jira plugin if we have two jira instances.
  - one of the jira instances is deployed privately for company projects.
  - another is jira cloud version for open-source projects.

## Task
Task represent an execution of specific plugin. In this version, task will be triggered by plugin name.
