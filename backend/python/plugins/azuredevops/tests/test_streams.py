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

from pydevlake.testing import assert_convert
import pydevlake.domain_layer.code as code
import pydevlake.domain_layer.devops as devops

from azuredevops.main import AzureDevOpsPlugin


def test_builds_stream():
    raw = {
        'properties': {},
        'tags': [],
        'validationResults': [],
        'plans': [{'planId': 'c672e778-a9e9-444a-b1e0-92f839c061e0'}],
        'triggerInfo': {},
        'id': 12,
        'buildNumber': 'azure-job',
        'status': 'completed',
        'result': 'succeeded',
        'queueTime': '2023-02-25T06:22:21.2237625Z',
        'startTime': '2023-02-25T06:22:32.8097789Z',
        'finishTime': '2023-02-25T06:23:04.0061884Z',
        'url': 'https://dev.azure.com/testorg/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/Builds/12',
        'definition': {
            'drafts': [],
            'id': 5,
            'name': 'johndoe.test-repo',
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
            'imageUrl': 'https://dev.azure.com/testorg/_apis/GraphProfile/MemberAvatars/aad.YmM1MzhmZWItOWZkZC03Y2Y4LTgwZTEtN2M1Njk1MGQwMjg5',
            'descriptor': 'aad.YmM1MzhmZWItOWZkZC03Y2Y4LTgwZTEtN2M1Njk1MGQwMjg5'
        },
        'requestedBy': {
            'displayName': 'John Doe',
            'url': 'https://spsprodcus5.vssps.visualstudio.com/A1def512a-251e-4668-9a5d-a4bc1f0da4aa/_apis/Identities/bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'id': 'bc538feb-9fdd-6cf8-80e1-7c56950d0289',
            'uniqueName': 'john.doe@merico.dev',
            'imageUrl': 'https://dev.azure.com/testorg/_apis/GraphProfile/MemberAvatars/aad.YmM1MzhmZWItOWZkZC03Y2Y4LTgwZTEtN2M1Njk1MGQwMjg5',
            'descriptor': 'aad.YmM1MzhmZWItOWZkZC03Y2Y4LTgwZTEtN2M1Njk1MGQwMjg5'
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
        'logs': {'id': 0,
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
            name=12,
            status=devops.CICDStatus.DONE,
            created_date='2023-02-25T06:22:32.8097789Z',
            finished_date='2023-02-25T06:23:04.0061884Z',
            result=devops.CICDResult.SUCCESS,
            duration_sec=28,
            environment=devops.CICDEnvironment.PRODUCTION,
            type=devops.CICDType.DEPLOYMENT,
            cicd_scope_id='johndoe/test-repo'
        ),
        devops.CiCDPipelineCommit(
            pipeline_id=12,
            commit_sha='40c59264e73fc5e1a6cab192f1622d26b7bd5c2a',
            branch='refs/heads/main',
            repo_id='johndoe/test-repo',
            repo='https://github.com/johndoe/test-repo'
        )
    ]

    assert_convert(AzureDevOpsPlugin, 'builds', raw, expected)


def test_jobs_stream():
    raw = {
        'previousAttempts': [],
        'id': 'cfa20e98-6997-523c-4233-f0a7302c929f',
        'parentId': '9ecf18fe-987d-5811-7c63-300aecae35da',
        'type': 'Job',
        'name': 'job_2',
        'build_id': 12, # Added by collector,
        'repo_id': 'johndoe/test-repo', # Added by collector,
        'startTime': '2023-02-25T06:22:36.8066667Z',
        'finishTime': '2023-02-25T06:22:43.2333333Z',
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
        'log': {'id': 10,
        'type': 'Container',
        'url': 'https://dev.azure.com/johndoe/7a3fd40e-2aed-4fac-bac9-511bf1a70206/_apis/build/builds/12/logs/10'},
        'task': None,
        'attempt': 1,
        'identifier': 'job_2.__default'
    }

    expected = devops.CICDTask(
        id='cfa20e98-6997-523c-4233-f0a7302c929f',
        name='job_2',
        pipeline_id=12,
        status=devops.CICDStatus.DONE,
        created_date='2023-02-25T06:22:36.8066667Z',
        finished_date='2023-02-25T06:22:43.2333333Z',
        result=devops.CICDResult.SUCCESS,
        type=devops.CICDType.BUILD,
        duration_sec=7,
        environment=devops.CICDEnvironment.PRODUCTION,
        cicd_scope_id='johndoe/test-repo'
    )

    assert_convert(AzureDevOpsPlugin, 'jobs', raw, expected)
