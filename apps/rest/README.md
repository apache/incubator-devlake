# Restful Api

## Source
Source represent a set of plugin's configuration. For example, if we has a Jira plugin who is responsible for collecting jira data, a source may correspond with a jira host.

## SourceTask
SourceTask represent a plugin execution, we can deliver more parameters to plugin when we create a source task. For example, we will send data scope(eg: jira issue) we want to collect to plugin when creating the source task.

> Notice: options for Source and SourceTask is defined by plugin itself, lake only take responsible to deliver the parameters.
