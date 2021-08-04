## gitlab-pond

## Summary of API calls to cover metrics

https://merico.feishu.cn/docs/doccnRszNS8i7mhvtjRynIGGR4g#

## Metrics we want to cover for a GitLab Project:

1. Number of contributors
2. Number of commits
3. Number of merge reviewers
4. Merge review time (defined as from the first comment to merge, the MR should have at least one comment to be considered as reviewed)
5. Number of merge requests
6. Merge review pass rate (defined as the percentage of merged MRs vs all MRs)
7. Added lines of code
8. Removed lines of code
9. Accumulated lines (defined as the sum of added lines of code and removed lines of code during a time window)

## The Endpoints

1. Commits API
  - Link to Docs: https://docs.gitlab.com/ee/api/commits.html#list-repository-commits 
  - Endpoint: GET /projects/:id/repository/commits?with_stats=true
  - Metrics: #1-5
2. Merge Requests API
  - Link to Docs: https://docs.gitlab.com/ee/api/merge_requests.html#list-project-merge-requests 
  - Endpoint: GET /projects/:id/merge_requests
  - Metrics: #6-8
3. Notes API
  - Link to Docs: https://docs.gitlab.com/ee/api/notes.html#list-all-merge-request-notes
  - Endpoint: GET /projects/:id/merge_requests/:merge_request_iid/notes
  - Metrics: #9
4. Projects API
  - Link to Docs: https://docs.gitlab.com/ee/api/projects.html#get-single-project
  - Endpoint: GET /projects/:id


