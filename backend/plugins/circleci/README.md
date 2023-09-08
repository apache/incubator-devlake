
# Concepts

## Pipeline

A `pipeline` is the historical execution for a specific commit, each time you push a new commit to trigger an execution produces a new `pipeline` accordingly.

## Workflow

A `pipeline` could have multiple `workflows`, A `workflow` orchestrates a set of `jobs`, it is a historical execution record.

## Job

A `workflow` could have multiple `jobs`, A `job` could have multiple `steps`, and it is a historical execution record as well.

## Step

Not important at the moment


# Domain Layer Conversion

Based on the above concept, we need to convert the `workflow` and `job` to the Domain Layer as `cicd_pipeline` and `cicd_task`.
As for the CircleCI `pipeline`, it should be mapped to `cicd_pipeline_commit` because it contains the `commit sha`.

It may look weird at first glance, but it is correct since Domain Layer presumes a Pipeline could have multiple repos while CircleCI has only one repo.