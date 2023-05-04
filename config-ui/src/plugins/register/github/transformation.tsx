/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import React, { useState, useEffect } from 'react';
import { FormGroup, InputGroup, TextArea, Tag, Switch, Icon, Intent, Colors } from '@blueprintjs/core';

import { ExternalLink, HelpTooltip, Divider } from '@/components';

import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const GitHubTransformation = ({ transformation, setTransformation }: Props) => {
  const [enableCICD, setEnableCICD] = useState(true);

  useEffect(() => {
    if (!transformation.deploymentPattern) {
      setEnableCICD(false);
    }
  }, [transformation]);

  const handleChangeEnableCICD = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (!checked) {
      setTransformation({
        ...transformation,
        deploymentPattern: undefined,
        productionPattern: undefined,
      });
    }

    setEnableCICD(checked);
  };

  return (
    <S.Transformation>
      {/* Issue Tracking */}
      <div className="issue-tracking">
        <h2>Issue Tracking</h2>
        <p>
          Tell DevLake what your issue labels mean to view metrics such as{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Metrics/BugAge">Bug Age</ExternalLink>,{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Metrics/MTTR">
            DORA - Median Time to Restore Service
          </ExternalLink>
          , etc.
        </p>
        <div className="issue-type">
          <div className="title">
            <span>Issue Type</span>
            <HelpTooltip content="DevLake defines three standard types of issues: FEATURE, BUG and INCIDENT. Set your issues to these three types with issue labels that match the RegEx." />
          </div>
          <div className="list">
            <FormGroup inline label="Requirement">
              <InputGroup
                placeholder="(feat|feature|proposal|requirement)"
                value={transformation.issueTypeRequirement ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    issueTypeRequirement: e.target.value,
                  })
                }
              />
            </FormGroup>
            <FormGroup inline label="Bug">
              <InputGroup
                placeholder="(bug|broken)"
                value={transformation.issueTypeBug ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    issueTypeBug: e.target.value,
                  })
                }
              />
            </FormGroup>
            <FormGroup
              inline
              label={
                <span>
                  Incident
                  <Tag minimal intent={Intent.PRIMARY} style={{ marginLeft: 4 }}>
                    DORA
                  </Tag>
                </span>
              }
            >
              <InputGroup
                placeholder="(incident|failure)"
                value={transformation.issueTypeIncident ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    issueTypeIncident: e.target.value,
                  })
                }
              />
            </FormGroup>
          </div>
        </div>
        <FormGroup
          inline
          label={
            <>
              <span>Issue Priority</span>
              <HelpTooltip content="Labels that match the RegEx will be set as the priority of an issue." />
            </>
          }
        >
          <InputGroup
            placeholder="(highest|high|medium|low|p0|p1|p2|p3)"
            value={transformation.issuePriority ?? ''}
            onChange={(e) =>
              setTransformation({
                ...transformation,
                issuePriority: e.target.value,
              })
            }
          />
        </FormGroup>
        <FormGroup
          inline
          label={
            <>
              <span>Issue Component</span>
              <HelpTooltip content="Labels that match the RegEx will be set as the component of an issue." />
            </>
          }
        >
          <InputGroup
            placeholder="component(.*)"
            value={transformation.issueComponent ?? ''}
            onChange={(e) =>
              setTransformation({
                ...transformation,
                issueComponent: e.target.value,
              })
            }
          />
        </FormGroup>
        <FormGroup
          inline
          label={
            <>
              <span>Issue Severity</span>
              <HelpTooltip content="Labels that match the RegEx will be set as the serverity of an issue." />
            </>
          }
        >
          <InputGroup
            placeholder="severity(.*)"
            value={transformation.issueSeverity ?? ''}
            onChange={(e) =>
              setTransformation({
                ...transformation,
                issueSeverity: e.target.value,
              })
            }
          />
        </FormGroup>
      </div>
      <Divider />
      {/* CI/CD */}
      <S.CICD>
        <h2>CI/CD</h2>
        <h3>
          <span>Deployment</span>
          <Tag minimal intent={Intent.PRIMARY} style={{ marginLeft: 8 }}>
            DORA
          </Tag>
          <div className="switch">
            <span>Enable</span>
            <Switch alignIndicator="right" inline checked={enableCICD} onChange={handleChangeEnableCICD} />
          </div>
        </h3>
        {enableCICD && (
          <>
            <p>
              Use Regular Expression to define Deployments in DevLake in order to measure DORA metrics.{' '}
              <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitHub#step-3---adding-transformation-rules-optional">
                Learn more
              </ExternalLink>
            </p>
            <div style={{ marginTop: 16 }}>Convert a GitHub Workflow run as a DevLake Deployment when: </div>
            <div className="text">
              <span>
                The name of the <strong>GitHub workflow run</strong> or <strong>one of its jobs</strong> matches
              </span>
              <InputGroup
                style={{ width: 200, margin: '0 8px' }}
                placeholder="(deploy|push-image)"
                value={transformation.deploymentPattern ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    deploymentPattern: e.target.value,
                    productionPattern: !e.target.value ? '' : transformation.productionPattern,
                  })
                }
              />
              <i style={{ color: '#E34040' }}>*</i>
              <HelpTooltip content="GitHub Workflow Runs: https://docs.github.com/en/actions/managing-workflow-runs/manually-running-a-workflow" />
            </div>
            <div className="text">
              <span>If the name also matches</span>
              <InputGroup
                style={{ width: 200, margin: '0 8px' }}
                disabled={!transformation.deploymentPattern}
                placeholder="prod(.*)"
                value={transformation.productionPattern ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    productionPattern: e.target.value,
                  })
                }
              />
              <span>, this Deployment is a ‘Production Deployment’</span>
              <HelpTooltip content="If you leave this field empty, all DevLake Deployments will be tagged as in the Production environment. " />
            </div>
          </>
        )}
      </S.CICD>
      <Divider />
      {/* Code Review */}
      <div>
        <h2>Code Review</h2>
        <p>
          If you use labels to identify types and components of pull requests, use the following RegExes to extract them
          into corresponding columns.{' '}
          <ExternalLink link="https://devlake.apache.org/docs/DataModels/DevLakeDomainLayerSchema#pull_requests">
            Learn More
          </ExternalLink>
        </p>
        <FormGroup
          inline
          label={
            <>
              <span>PR Type</span>
              <HelpTooltip content="Labels that match the RegEx will be set as the type of a pull request." />
            </>
          }
        >
          <InputGroup
            placeholder="type(.*)$"
            value={transformation.prType ?? ''}
            onChange={(e) => setTransformation({ ...transformation, prType: e.target.value })}
          />
        </FormGroup>
        <FormGroup
          inline
          label={
            <>
              <span>PR Component</span>
              <HelpTooltip content="Labels that match the RegEx will be set as the component of a pull request." />
            </>
          }
        >
          <InputGroup
            placeholder="component(.*)$"
            value={transformation.prComponent ?? ''}
            onChange={(e) =>
              setTransformation({
                ...transformation,
                prComponent: e.target.value,
              })
            }
          />
        </FormGroup>
      </div>
      <Divider />
      {/* Cross-domain */}
      <div>
        <h2>Cross-domain</h2>
        <p>
          Connect entities across domains to measure metrics such as{' '}
          <ExternalLink link="https://devlake.apache.org/docs/Metrics/BugCountPer1kLinesOfCode">
            Bug Count per 1k Lines of Code
          </ExternalLink>
          .
        </p>
        <FormGroup
          inline
          label={
            <div className="label">
              <span>Connect PRs and Issues</span>
              <HelpTooltip
                content={
                  <>
                    <div>
                      <Icon icon="tick-circle" size={12} color={Colors.GREEN4} style={{ marginRight: '4px' }} />
                      Example 1: PR #321 body contains "<strong>Closes #1234</strong>" (PR #321 and issue #1234 will be
                      mapped by the following RegEx)
                    </div>
                    <div>
                      <Icon icon="delete" size={12} color={Colors.RED4} style={{ marginRight: '4px' }} />
                      Example 2: PR #321 body contains "<strong>Related to #1234</strong>" (PR #321 and issue #1234 will
                      NOT be mapped by the following RegEx)
                    </div>
                  </>
                }
              />
            </div>
          }
        >
          <TextArea
            value={transformation.prBodyClosePattern ?? ''}
            placeholder="(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[s]*.*(((and )?(#|https://github.com/%s/%s/issues/)d+[ ]*)+)"
            onChange={(e) =>
              setTransformation({
                ...transformation,
                prBodyClosePattern: e.target.value,
              })
            }
            fill
            rows={2}
          />
        </FormGroup>
      </div>
    </S.Transformation>
  );
};
