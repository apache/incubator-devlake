import base64
import typing

from pydevlake.api import API, request_hook, Paginator, TokenPaginator, response_hook, Response, Request


class AzureDevOpsAPI(API):

    def __init__(self, base_url: str, pat: str):
        self.base_url = base_url
        self.pat = pat

    def base_url(self):
        return self.base_url

    @property
    def paginator(self) -> Paginator:
        TokenPaginator("", "", "")
        return None

    @response_hook  # how to use this?
    def on_response(self, response: Response):
        return response

    @response_hook  # how to use this?
    def rate_limiter(self, response: Response):
        # TODO
        return response

    @request_hook
    def authenticate(self, request: Request):
        pat_b64 = base64.b64encode((':' + self.pat).encode()).decode()
        if self.pat:
            request.headers['Authorization'] = 'Basic ' + pat_b64

    def accounts(self) -> Response:
        return self.get('https://app.vssps.visualstudio.com', '_apis/accounts')

    def orgs(self) -> list[str]:
        response = self.accounts()
        return [acct["AccountName"] for acct in response.json]

    def projects(self, org: str) -> Response:
        return self.get(f'{org}/_apis/projects')

    # Get a project
    def project(self, org: str, project: str) -> Response:
        return self.get(f'{org}/_apis/projects/{project}')

    # List repos under an org
    def git_repos(self, org: str, project: str) -> Response:
        return self.get(f'{org}/{project}/_apis/git/repositories')

    def git_repo_pull_requests(self, org: str, project: str, repo_id: str) -> Response:
        # see https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests?view=azure-devops-rest-7.1&tabs=HTTP
        return self.get(f'{org}/{project}/_apis/git/repositories/{repo_id}/pullrequests')

    def git_repo_pull_request_commits(self, org: str, project: str, repo_id: str, pull_request_id: int) -> Response:
        return self.get(f'{org}/{project}/_apis/git/repositories/{repo_id}/pullRequests/{pull_request_id}/commits')

    def git_repo_pull_request_comments(self, org: str, project: str, repo_id: str, pull_request_id: int) -> Response:
        return self.get(f'{org}/{project}/_apis/git/repositories/{repo_id}/pullRequests/{pull_request_id}/threads')

    # not needed
    def commits(self, org: str, project: str, repo_id: str) -> Response:
        return self.get(f'{org}/{project}/_apis/git/repositories/{repo_id}/commits')

    def builds(self, org: str, project: str) -> Response:
        return self.get(f'{org}/{project}/_apis/build/builds')

    def jobs(self, org: str, project: str, build_id: int) -> Response:
        return self.get(f'{org}/{project}/_apis/build/builds/{build_id}/timeline')

    # unused
    def deployments(self, org: str, project: str) -> Response:
        return self.get(f'{org}/{project}/_apis/release/deployments')

    # unused
    def releases(self, org: str, project: str) -> Response:
        return self.get(f'{org}/{project}/_apis/release/releases')

    def parse_response(self, res: Response) -> typing.Iterable[dict]:
        return res.json['value']
