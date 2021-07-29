## gitlab-pond

## Summary of API calls to cover metrics

https://merico.feishu.cn/docs/doccnRszNS8i7mhvtjRynIGGR4g#

## Metrics we want to cover for a GitLab Project:

1. Number of contributors
2. Number of commits
3. Removed lines of code
4. Added lines of code
5. Accumulated lines (defined as the sum of added lines of code and removed lines of code during a time window)
6. Number of reviewers
7. Number of merge requests
8. MR review pass ratio (defined as the percentage of merged MRs vs all MRs)
9. Merge review time (defined as from the first comment to merge, the MR should have at least one comment to be considered as reviewed)
10. Absolute lines (defined as the diff between two snapshot of the codebase)

## The Endpoints

1. Commits API
  - Link to Docs: https://docs.gitlab.com/ee/api/commits.html#list-repository-commits 
  - Endpoint: GET /projects/:id/repository/commits?all=true&with_stats=true
  - Metrics: 1-5
2. Merge Requests API
  - Link to Docs: https://docs.gitlab.com/ee/api/merge_requests.html#list-project-merge-requests 
  - Endpoint: GET /projects/:id/merge_requests
  - Metrics: 6-8
3. Notes API
  - Link to Docs: 
  - Endpoint 
  - Metrics: 9

