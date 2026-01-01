# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import pytest

import pydevlake.domain_layer.code as code
import pydevlake.domain_layer.devops as devops
from azuredevops.main import AzureDevOpsPlugin
from pydevlake.testing import assert_stream_convert, ContextBuilder


@pytest.fixture
def context():
    return (
        ContextBuilder(AzureDevOpsPlugin())
        .with_connection(token='token')
        .with_scope_config(deployment_pattern='deploy',
                           production_pattern='prod')
        .with_scope('johndoe/test-repo', url='https://github.com/johndoe/test-repo')
        .build()
    )


def test_builds_stream(context):
    raw = {
        'properties': {},
        'tags': [],
        'validationResults': [],
        'plans': [{'planId': 'c672e778-a9e9-444a-b1e0-92f839c061e0'}],
        'triggerInfo': {
            "ci.sourceBranch": "refs/heads/main",
            "ci.sourceSha": "40e3d9cb9f208f431cf1fb0e33963f5a1405491b",
            "ci.message": "Add azure-pipelines.yml jobs to Azure Pipelines",
            "ci.triggerRepository": "eaf116f6-821f-42d7-920e-a867e564302e"
        },
        '_links':{                    
            "self": {
                "href": "https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_apis/build/Builds/15"
            },
            "web": {
                "href": "https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_build/results?buildId=15"
            },
            "sourceVersionDisplayUri": {
                "href": "https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_apis/build/builds/15/sources"
            },
            "timeline": {
                "href": "https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_apis/build/builds/15/Timeline"
            },
            "badge": {
                "href": "https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_apis/build/status/1"
            }
        },
        'id': 12,
        'buildNumber': 'azure-job',
        'status': 'completed',
        'result': 'succeeded',
        'queueTime': '2023-02-25T06:22:21.2237625Z',
        'start_time': '2023-02-25T06:22:32.8097789Z',
        'finish_time': '2023-02-25T06:23:04.0061884Z',
        'url': 'https://dev.azure.com/testorg/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/Builds/12',
        'definition': {
            'drafts': [],
            'id': 5,
            'name': 'deploy_to_prod',
            'url': 'https://dev.azure.com/testorg/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/Definitions/5?revision=1',
            'uri': 'vstfs:///Build/Definition/5',
            'path': '\\',
            'type': 'build',
            'queueStatus': 'enabled',
            'revision': 1,
            'project': {
                'id': '7a3fd40e-2aed-4fac-bac9-511bf1a70206',
                'name': 'test-project',
                'url': 'https://dev.azure.com/testorg/_apis/projects/7a3fd40e-2aed-4fac-bac9-511bf1a70206',
                'state': 'wellFormed',
                'revision': 11,
                'visibility': 'private',
                'lastUpdateTime': '2023-01-26T19:38:04.267Z'
            }
        },
        'project': {
            'id': '7a3fd40e-2aed-4fac-bac9-511bf1a70206',
            'name': 'Test project',
            'url': 'https://dev.azure.com/testorg/_apis/projects/7a3fd40e-2aed-4fac-bac9-511bf1a70206',
        },
        'uri': 'vstfs:///Build/Build/12',
        'sourceBranch': 'refs/heads/main',
        'sourceVersion': '40c59264e73fc5e1a6cab192f1622d26b7bd5c2a',
        'queue': {
            'id': 9,
            'name': 'Azure Pipelines',
            'pool': {'id': 9, 'name': 'Azure Pipelines', 'isHosted': True}
        },
        'priority': 'normal',
        'reason': 'manual',
        'requestedFor': {
            'displayName': 'John Doe',
            'url': 'https://spsprodcus5.vssps.visualstudio.com/A1def512a-251e-4668-9a5d-a4bc1f0da4aa/_apis/Identities/bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'id': 'bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'uniqueName': 'john.dow@merico.dev',
            'imageUrl': 'https://dev.azure.com/testorg/_apis/GraphProfile/MemberAvatars/aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5',
            'descriptor': 'aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5'
        },
        'requestedBy': {
            'displayName': 'John Doe',
            'url': 'https://spsprodcus5.vssps.visualstudio.com/A1def512a-251e-4668-9a5d-a4bc1f0da4aa/_apis/Identities/bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'id': 'bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'uniqueName': 'john.doe@merico.dev',
            'imageUrl': 'https://dev.azure.com/testorg/_apis/GraphProfile/MemberAvatars/aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5',
            'descriptor': 'aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5'
        },
        'lastChangedDate': '2023-02-25T06:23:04.343Z',
        'lastChangedBy': {
            'displayName': 'Microsoft.VisualStudio.Services.TFS',
            'url': 'https://spsprodcus5.vssps.visualstudio.com/A1def512a-251e-4668-9a5d-a4bc1f0da4aa/_apis/Identities/00000002-0000-8888-8000-000000000000',
            'id': '00000002-0000-8888-8000-000000000000',
            'uniqueName': '00000002-0000-8888-8000-000000000000@2c895908-04e0-4952-89fd-54b0046d6288',
            'imageUrl': 'https://dev.azure.com/testorg/_apis/GraphProfile/MemberAvatars/s2s.MDAwMDAwMDItMDAwMC04ODg4LTgwMDAtMDAwMDAwMDAwMDAwQDJjODk1OTA4LTA0ZTAtNDk1Mi04OWZkLTU0YjAwNDZkNjI4OA',
            'descriptor': 's2s.MDAwMDAwMDItMDAwMC04ODg4LTgwMDAtMDAwMDAwMDAwMDAwQDJjODk1OTA4LTA0ZTAtNDk1Mi04OWZkLTU0YjAwNDZkNjI4OA'
        },
        'orchestrationPlan': {'planId': 'c672e778-a9e9-444a-b1e0-92f839c061e0'},
        'logs': {
            'id': 0,
            'type': 'Container',
            'url': 'https://dev.azure.com/testorg/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/builds/12/logs'
        },
        'repository': {
            'id': 'johndoe/test-repo',
            'type': 'GitHub'
        },
        'retainedByRelease': False,
        'triggeredByBuild': None,
        'appendCommitMessageToRunName': True
    }

    expected = [
        devops.CICDPipeline(
            name='deploy_to_prod',
            status=devops.CICDStatus.DONE,
            created_date='2023-02-25T06:22:21.2237625Z',
            queued_date='2023-02-25T06:22:21.2237625Z',
            started_date='2023-02-25T06:22:32.8097789Z',
            finished_date='2023-02-25T06:23:04.0061884Z',
            result=devops.CICDResult.SUCCESS,
            original_status='Completed',
            original_result='Succeeded',
            duration_sec=31.196409940719604,
            environment=devops.CICDEnvironment.PRODUCTION,
            type=devops.CICDType.DEPLOYMENT,
            cicd_scope_id=context.scope.domain_id(),
            display_title='Add azure-pipelines.yml jobs to Azure Pipelines',
            url='https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_build/results?buildId=15'
        ),
        devops.CiCDPipelineCommit(
            pipeline_id='azuredevops:Build:1:12',
            commit_sha='40c59264e73fc5e1a6cab192f1622d26b7bd5c2a',
            branch='refs/heads/main',
            repo_id=context.scope.domain_id(),
            repo_url='https://github.com/johndoe/test-repo',
            display_title='Add azure-pipelines.yml jobs to Azure Pipelines',
            url='https://dev.azure.com/linweihoumerico-lake/e8af9e7b-d4bf-4afd-9d0a-c9f8dfac1d59/_build/results?buildId=15'
        )
    ]

    assert_stream_convert(AzureDevOpsPlugin, 'builds', raw, expected, context)


def test_jobs_stream(context):
    raw = {
        'previousAttempts': [],
        'id': 'cfa20e98-6997-523c-4233-f0a7302c929f',
        'parentId': '9ecf18fe-987d-5811-7c63-300aecae35da',
        'type': 'Job',
        'name': 'deploy production',
        'build_id': 'azuredevops:Build:1:12',  # Added by collector,
        'start_time': '2023-02-25T06:22:36.8066667Z',
        'finish_time': '2023-02-25T06:22:43.2333333Z',
        'currentOperation': None,
        'percentComplete': None,
        'state': 'completed',
        'result': 'succeeded',
        'resultCode': None,
        'changeId': 18,
        'lastModified': '0001-01-01T00:00:00',
        'workerName': 'Hosted Agent',
        'queueId': 9,
        'order': 1,
        'details': None,
        'errorCount': 0,
        'warningCount': 0,
        'url': None,
        'log': {
            'id': 10,
            'type': 'Container',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/builds/12/logs/10'
        },
        'task': None,
        'attempt': 1,
        'identifier': 'job_2.__default'
    }

    expected = devops.CICDTask(
        id='cfa20e98-6997-523c-4233-f0a7302c929f',
        name='deploy production',
        pipeline_id='azuredevops:Build:1:12',
        status=devops.CICDStatus.DONE,
        original_status='Completed',
        original_result='Succeeded',
        created_date='2023-02-25T06:22:36.8066667Z',
        started_date='2023-02-25T06:22:36.8066667Z',
        finished_date='2023-02-25T06:22:43.2333333Z',
        result=devops.CICDResult.SUCCESS,
        type=devops.CICDType.DEPLOYMENT,
        duration_sec=6.426667213439941,
        environment=devops.CICDEnvironment.PRODUCTION,
        cicd_scope_id=context.scope.domain_id()
    )
    assert_stream_convert(AzureDevOpsPlugin, 'jobs', raw, expected, context)


def test_pull_requests_stream(context):
    raw = {
        'repository': {
            'id': '0d50ba13-f9ad-49b0-9b21-d29eda50ca33',
            'name': 'test-repo2',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33',
            'project': {
                'id': '7a3fd40e-2aed-4fac-bac9-511bf1a70206',
                'name': 'test-project',
                'state': 'unchanged',
                'visibility': 'unchanged',
                'lastUpdateTime': '0001-01-01T00:00:00'
            }
        },
        'pullRequestId': 1,
        'codeReviewId': 1,
        'status': 'active',
        'createdBy': {
            'displayName': 'John Doe',
            'url': 'https://spsprodcus5.vssps.visualstudio.com/A1def512a-251e-4668-9a5d-a4bc1f0da4aa/_apis/Identities/bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            '_links': {
                'avatar': {
                    'href': 'https://dev.azure.com/johndoe/_apis/GraphProfile/MemberAvatars/aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5'
                }
            },
            'id': 'bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'uniqueName': 'john.doe@merico.dev',
            'imageUrl': 'https://dev.azure.com/johndoe/_api/_common/identityImage?id=bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'descriptor': 'aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5'
        },
        'creationDate': '2023-02-07T04:41:26.6424314Z',
        'title': 'ticket-2 PR',
        'description': 'Updated main.java by ticket-2',
        'sourceRefName': 'refs/heads/ticket-2',
        'targetRefName': 'refs/heads/main',
        'mergeStatus': 'succeeded',
        'isDraft': False,
        'mergeId': '99da29c2-4d27-4620-989f-5b59908917cd',
        'lastMergeSourceCommit': {
            'commitId': '85ede91717145a1e6e2bdab4cab689ac8f2fa3a2',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33/commits/85ede91717145a1e6e2bdab4cab689ac8f2fa3a2'
        },
        'lastMergeTargetCommit': {
            'commitId': '4bc26d92b5dbee7837a4d221035a4e2f8df120b2',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33/commits/4bc26d92b5dbee7837a4d221035a4e2f8df120b2'
        },
        'lastMergeCommit': {
            'commitId': 'ebc6c7a2a5e3c155510d0ba44fd4385bf7ae6e22',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33/commits/ebc6c7a2a5e3c155510d0ba44fd4385bf7ae6e22'
        },
        'reviewers': [
            {
                'reviewerUrl': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33/pullRequests/1/reviewers/bc538feb-9fdd-6cf8-80e1-7c56950d0289',
                'vote': 0,
                'hasDeclined': False,
                'isFlagged': False,
                'displayName': 'John Doe',
                'url': 'https://spsprodcus5.vssps.visualstudio.com/A1def512a-251e-4668-9a5d-a4bc1f0da4aa/_apis/Identities/bc538feb-9fdd-6cf8-80e1-7c56950d0289',
                '_links': {'avatar': {
                    'href': 'https://dev.azure.com/johndoe/_apis/GraphProfile/MemberAvatars/aad.YmM1MzhmZWItOWZkZC03Y2Y4LT3wXTXtN2M1Njk1MGQwMjg5'}},
                'id': 'bc538feb-9fdd-6cf8-80e1-7c56950d0289',
                'uniqueName': 'john.doe@merico.dev',
                'imageUrl': 'https://dev.azure.com/johndoe/_api/_common/identityImage?id=bc538feb-9fdd-6cf8-80e1-7c56950d0289'
            }
        ],
        'labels': [
            {
                'id': '98db191b-f0a5-421b-8433-e982ad05fe06',
                'name': 'feature',
                'active': True
            }
        ],
        'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33/pullRequests/1',
        'supportsIterations': True
    }

    expected = code.PullRequest(
        base_repo_id='azuredevops:GitRepository:1:johndoe/test-repo',
        head_repo_id='azuredevops:GitRepository:1:johndoe/test-repo',
        status='OPEN',
        original_status='active',
        title='ticket-2 PR',
        description='Updated main.java by ticket-2',
        url='https://github.com/johndoe/test-repo/pullrequest/1',
        author_name='John Doe',
        author_id='bc538feb-9fdd-6cf8-80e1-7c56950d0289',
        pull_request_key=1,
        created_date='2023-02-07T04:41:26.6424314Z',
        merged_date=None,
        closed_date=None,
        type='feature',
        component="",
        merge_commit_sha='ebc6c7a2a5e3c155510d0ba44fd4385bf7ae6e22',
        head_ref='refs/heads/ticket-2',
        base_ref='refs/heads/main',
        head_commit_sha='85ede91717145a1e6e2bdab4cab689ac8f2fa3a2',
        base_commit_sha='4bc26d92b5dbee7837a4d221035a4e2f8df120b2'
    )

    assert_stream_convert(AzureDevOpsPlugin, 'gitpullrequests', raw, expected, context)


def test_pull_request_commits_stream():
    raw = {
        'commitId': '85ede91717145a1e6e2bdab4cab689ac8f2fa3a2',
        'author': {
            'name': 'John Doe',
            'email': 'john.doe@merico.dev',
            'date': '2023-02-07T04:49:28Z'
        },
        'committer': {
            'name': 'John Doe',
            'email': 'john.doe@merico.dev',
            'date': '2023-02-07T04:49:28Z'
        },
        'comment': 'Fixed main.java',
        'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/git/repositories/0d50ba13-f9ad-49b0-9b21-d29eda50ca33/commits/85ede91717145a1e6e2bdab4cab689ac8f2fa3a2',
        'pull_request_id': "azuredevops:gitpullrequest:1:12345"
        # This is not part of the API response, but is added in collect method
    }

    expected = code.PullRequestCommit(
        commit_sha='85ede91717145a1e6e2bdab4cab689ac8f2fa3a2',
        pull_request_id="azuredevops:gitpullrequest:1:12345",
        commit_author_name='John Doe',
        commit_authored_date='2023-02-07T04:49:28Z',
        commit_author_email='john.doe@merico.dev'
    )

    assert_stream_convert(AzureDevOpsPlugin, 'gitpullrequestcommits', raw, expected)


@pytest.fixture
def context_with_environment_pattern():
    """Context with environment_pattern configured to extract environment names from job names."""
    return (
        ContextBuilder(AzureDevOpsPlugin())
        .with_connection(token='token')
        .with_scope_config(
            deployment_pattern='deploy',
            production_pattern='prod',
            # Pattern to extract environment name from job names like 'deploy_xxxx-prod_helm'
            environment_pattern=r'(?:deploy|predeploy)[_-](.+?)(?:[_-](?:helm|terraform))?$'
        )
        .with_scope('johndoe/test-repo', url='https://github.com/johndoe/test-repo')
        .build()
    )


def test_jobs_stream_with_environment_pattern(context_with_environment_pattern):
    """Test that environment_pattern extracts environment name and uses it for production matching."""
    raw = {
        'previousAttempts': [],
        'id': 'cfa20e98-6997-523c-4233-f0a7302c929f',
        'parentId': '9ecf18fe-987d-5811-7c63-300aecae35da',
        'type': 'Job',
        'name': 'deploy_xxxx-prod_helm',  # environment name 'xxxx-prod' should be extracted
        'build_id': 'azuredevops:Build:1:12',
        'start_time': '2023-02-25T06:22:36.8066667Z',
        'finish_time': '2023-02-25T06:22:43.2333333Z',
        'currentOperation': None,
        'percentComplete': None,
        'state': 'completed',
        'result': 'succeeded',
        'resultCode': None,
        'changeId': 18,
        'lastModified': '0001-01-01T00:00:00',
        'workerName': 'Hosted Agent',
        'queueId': 9,
        'order': 1,
        'details': None,
        'errorCount': 0,
        'warningCount': 0,
        'url': None,
        'log': {
            'id': 10,
            'type': 'Container',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/builds/12/logs/10'
        },
        'task': None,
        'attempt': 1,
        'identifier': 'deploy_xxxx-prod_helm.__default'
    }

    expected = devops.CICDTask(
        id='cfa20e98-6997-523c-4233-f0a7302c929f',
        name='deploy_xxxx-prod_helm',
        pipeline_id='azuredevops:Build:1:12',
        status=devops.CICDStatus.DONE,
        original_status='Completed',
        original_result='Succeeded',
        created_date='2023-02-25T06:22:36.8066667Z',
        started_date='2023-02-25T06:22:36.8066667Z',
        finished_date='2023-02-25T06:22:43.2333333Z',
        result=devops.CICDResult.SUCCESS,
        type=devops.CICDType.DEPLOYMENT,
        duration_sec=6.426667213439941,
        environment=devops.CICDEnvironment.PRODUCTION,  # Should match because 'xxxx-prod' contains 'prod'
        cicd_scope_id=context_with_environment_pattern.scope.domain_id()
    )
    assert_stream_convert(AzureDevOpsPlugin, 'jobs', raw, expected, context_with_environment_pattern)


def test_jobs_stream_with_environment_pattern_non_prod(context_with_environment_pattern):
    """Test that non-prod environments are correctly identified."""
    raw = {
        'previousAttempts': [],
        'id': 'cfa20e98-6997-523c-4233-f0a7302c929f',
        'parentId': '9ecf18fe-987d-5811-7c63-300aecae35da',
        'type': 'Job',
        'name': 'deploy_xxxx-dev_helm',  # environment name 'xxxx-dev' should be extracted, not prod
        'build_id': 'azuredevops:Build:1:12',
        'start_time': '2023-02-25T06:22:36.8066667Z',
        'finish_time': '2023-02-25T06:22:43.2333333Z',
        'currentOperation': None,
        'percentComplete': None,
        'state': 'completed',
        'result': 'succeeded',
        'resultCode': None,
        'changeId': 18,
        'lastModified': '0001-01-01T00:00:00',
        'workerName': 'Hosted Agent',
        'queueId': 9,
        'order': 1,
        'details': None,
        'errorCount': 0,
        'warningCount': 0,
        'url': None,
        'log': {
            'id': 10,
            'type': 'Container',
            'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/builds/12/logs/10'
        },
        'task': None,
        'attempt': 1,
        'identifier': 'deploy_xxxx-dev_helm.__default'
    }

    expected = devops.CICDTask(
        id='cfa20e98-6997-523c-4233-f0a7302c929f',
        name='deploy_xxxx-dev_helm',
        pipeline_id='azuredevops:Build:1:12',
        status=devops.CICDStatus.DONE,
        original_status='Completed',
        original_result='Succeeded',
        created_date='2023-02-25T06:22:36.8066667Z',
        started_date='2023-02-25T06:22:36.8066667Z',
        finished_date='2023-02-25T06:22:43.2333333Z',
        result=devops.CICDResult.SUCCESS,
        type=devops.CICDType.DEPLOYMENT,
        duration_sec=6.426667213439941,
        environment=None,  # Should be None because 'xxxx-dev' does not contain 'prod'
        cicd_scope_id=context_with_environment_pattern.scope.domain_id()
    )
    assert_stream_convert(AzureDevOpsPlugin, 'jobs', raw, expected, context_with_environment_pattern)


def test_stage_record_collected():
    """Test that Stage records are also collected (not just Job records)."""
    context = (
        ContextBuilder(AzureDevOpsPlugin())
        .with_connection(token='token')
        .with_scope_config(
            deployment_pattern='deploy',
            production_pattern='prod'
        )
        .with_scope('johndoe/test-repo', url='https://github.com/johndoe/test-repo')
        .build()
    )

    raw = {
        'previousAttempts': [],
        'id': 'stage-id-123',
        'parentId': None,
        'type': 'Stage',  # This is a Stage record
        'name': 'deploy_prod_stage',
        'build_id': 'azuredevops:Build:1:12',
        'start_time': '2023-02-25T06:22:36.8066667Z',
        'finish_time': '2023-02-25T06:22:43.2333333Z',
        'currentOperation': None,
        'percentComplete': None,
        'state': 'completed',
        'result': 'succeeded',
        'resultCode': None,
        'changeId': 18,
        'lastModified': '0001-01-01T00:00:00',
        'workerName': None,
        'queueId': None,
        'order': 1,
        'details': None,
        'errorCount': 0,
        'warningCount': 0,
        'url': None,
        'log': None,
        'task': None,
        'attempt': 1,
        'identifier': 'deploy_prod_stage'
    }

    expected = devops.CICDTask(
        id='stage-id-123',
        name='deploy_prod_stage',
        pipeline_id='azuredevops:Build:1:12',
        status=devops.CICDStatus.DONE,
        original_status='Completed',
        original_result='Succeeded',
        created_date='2023-02-25T06:22:36.8066667Z',
        started_date='2023-02-25T06:22:36.8066667Z',
        finished_date='2023-02-25T06:22:43.2333333Z',
        result=devops.CICDResult.SUCCESS,
        type=devops.CICDType.DEPLOYMENT,
        duration_sec=6.426667213439941,
        environment=devops.CICDEnvironment.PRODUCTION,
        cicd_scope_id=context.scope.domain_id()
    )
    assert_stream_convert(AzureDevOpsPlugin, 'jobs', raw, expected, context)
