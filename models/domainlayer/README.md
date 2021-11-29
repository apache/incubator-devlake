# Why

We have Jira/TAPD for **Project Management**, we have Github/Gitlab/BitBucket for **Code Hosting**, that is, we have
multiple _Platforms_ for a certain type of problem. So, how can we calculate metrics across different _Platforms_?

For example, some users may use Jira as their **Project Management** platform, the others might opt for TAPD, if we
were to implement a Requirement Count metrics for all users, should we implement 2 charts for Jira and TAPD
independently? It's too impractical to begin with.


# How

Domain Layer is designed to solve the problem by offering a set of Platform Independent Entities, Devlake divides all
platforms into three categories: Project Management / Code Hosting / Devops, by abstracting common properties from
different platforms, we can define a set of Domain Entities for each category.

The following rules make sure Domain Layer Entities serve its purpose

1. Every platform specific entity can be mapped (or split) to one (or more) Domain Layer Entity
2. Every Domain Layer Entity contains enough information for metrics calculation
3. Domain Layer Entity should contains some sort of pointer to its origin record, and all entities should use the same
   pointer schema

# What

## Domain Layer Entity

- Each entity has a incremental primary key named `ID` which is platform independent
- Each entity has a `OriginKey` describe its origin record in format `<Plugin>:<Entity>:<PK0>:<PK1>`
- All fields needed for metric calculation

PS: `many-to-many` Entities are not considered as Domain Layer Entity, but part of the system, so they are freed from
these rules.


## Conversion Plugin

- Read data from platform specific table, convert and store record into one or multiple domain table(s)
- Set `OriginKey`
- Fields conversion

Sample code:

```go

type Issue struct {
    ID              uint64
    OriginKey       string
    ...
}

issue := Issue {
    OriginKey: "jira:jira_issues:10"
    ...
}

```
