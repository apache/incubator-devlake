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

import { useState, useEffect } from 'react';
import { Radio, Icon, Collapse, InputGroup, Button, Intent, RadioGroup } from '@blueprintjs/core';

import { ExternalLink, IconButton, Dialog, FormItem } from '@/components';
import JiraIssueTipsImg from '@/images/jira-issue-tips.png';
import { operator } from '@/utils';

import * as API from '../api';

import * as S from './styled';

interface Props {
  connectionId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const CrossDomain = ({ connectionId, transformation, setTransformation }: Props) => {
  const [radio, setRadio] = useState<'repo' | 'commitSha'>('repo');
  const [repoTips, setRepoTips] = useState(false);
  const [repoLinks, setRepoLinks] = useState([]);
  const [isOpen, setIsOpen] = useState(false);
  const [step, setStep] = useState(1);
  const [issueKey, setIssueKey] = useState('');
  const [searching, setSearching] = useState(false);
  const [applicationTypes, setApplicationTypes] = useState<string[]>([]);
  const [applicationType, setApplicationType] = useState<string>();
  const [operating, setOperating] = useState(false);
  const [devPanelCommits, setDevPanelCommits] = useState<string[]>([]);

  useEffect(() => {
    if (transformation.remotelinkCommitShaPattern) {
      setRadio('commitSha');
    }

    if (transformation.remotelinkRepoPattern) {
      setRepoLinks(transformation.remotelinkRepoPattern);
    }
  }, [transformation]);

  const handleChangeRadio = (r: 'repo' | 'commitSha') => {
    if (r === 'repo') {
      setTransformation({
        ...transformation,
        remotelinkCommitShaPattern: '',
      });
    }

    if (r === 'commitSha') {
      setTransformation({
        ...transformation,
        remotelinkRepoPattern: [],
      });
    }

    setRadio(r);
  };

  const handleToggleRepoTips = () => setRepoTips(!repoTips);

  const handleChangeRepoLinks = (index: number, value: string) => {
    const newValue = repoLinks.map((link, i) => (index === i ? value : link));
    setTransformation({
      ...transformation,
      remotelinkRepoPattern: newValue,
    });
  };

  const handleAddRepoLinks = () => {
    const newValue = [...repoLinks, ''];
    setTransformation({
      ...transformation,
      remotelinkRepoPattern: newValue,
    });
  };

  const handleDeleteRepoLinks = (index: number) => {
    const newValue = repoLinks.filter((_, i) => i !== index);
    setTransformation({
      ...transformation,
      remotelinkRepoPattern: newValue,
    });
  };

  const handleShowDevPanel = () => setIsOpen(true);
  const handleHideDevPanel = () => setIsOpen(false);

  const handleSearch = async () => {
    const [success, res] = await operator(() => API.getApplicationTypes(connectionId, { key: issueKey }), {
      setOperating: setSearching,
      hideToast: true,
    });

    if (success) {
      setApplicationTypes(res);
    }
  };

  const handleSubmit = async () => {
    if (step === 1 && applicationType) {
      const [success, res] = await operator(
        () => API.getDevPanelCommits(connectionId, { key: issueKey, applicationType }),
        {
          setOperating,
          hideToast: true,
        },
      );

      if (success) {
        setStep(2);
        setDevPanelCommits(res);
        return;
      }
    }

    handleHideDevPanel();
  };

  const handleCancel = () => {
    if (step === 1) {
      handleHideDevPanel();
      return;
    }

    setStep(1);
  };

  return (
    <S.CrossDomain>
      <h2>Cross-domain</h2>
      <p>
        Connect `commits` and `issues` to measure metrics such as{' '}
        <ExternalLink link="https://devlake.apache.org/docs/Metrics/BugCountPer1kLinesOfCode">
          Bug Count per 1k Lines of Code
        </ExternalLink>{' '}
        or man hour distribution on different work types.
      </p>
      <div className="radio">
        <div className="radio-item">
          <Radio checked={radio === 'repo'} onChange={() => handleChangeRadio('repo')} />
          <div className="content">
            <h5>Connect Jira issues and commits via Jira issues’ remote links that match the following pattern</h5>
            <p onClick={handleToggleRepoTips}>
              Input pattern(s) to match and parse commits and repo identifiers from issue remote links. See examples{' '}
              <Icon icon={!repoTips ? 'chevron-down' : 'chevron-up'} style={{ marginLeft: 8, cursor: 'pointer' }} />
            </p>
            <Collapse isOpen={repoTips}>
              <img src={JiraIssueTipsImg} width="100%" alt="" />
            </Collapse>
            {radio === 'repo' && (
              <>
                {repoLinks.map((link, i) => (
                  <div className="input">
                    <InputGroup
                      key={i}
                      placeholder="E.g. https://gitlab.com/{namespace}/{repo_name}/-/commit/{commit_sha}"
                      value={link}
                      onChange={(e) => handleChangeRepoLinks(i, e.target.value)}
                    />
                    {repoLinks.length > 1 && (
                      <IconButton icon="cross" tooltip="Delete" onClick={() => handleDeleteRepoLinks(i)} />
                    )}
                  </div>
                ))}
                <Button
                  outlined
                  intent={Intent.PRIMARY}
                  icon="add"
                  text="Add a Pattern"
                  onClick={() => handleAddRepoLinks()}
                />
              </>
            )}
          </div>
        </div>
        <div className="radio-item">
          <Radio checked={radio === 'commitSha'} onChange={() => handleChangeRadio('commitSha')} />
          <div className="content">
            <h5>
              Connect Jira issues and commits via{' '}
              <ExternalLink link="Links to: https://support.atlassian.com/jira-software-cloud/docs/view-development-information-for-an-issue/">
                development panel
              </ExternalLink>
            </h5>
            <p>Finish the configuration so DevLake can get your Git data from your Jira development panel.</p>
            {radio === 'commitSha' && <Button text="Configure" onClick={handleShowDevPanel} />}
          </div>
        </div>
      </div>
      <Dialog
        style={{ width: 820 }}
        isOpen={isOpen}
        title="Configure the `Application Type` and `Commit Pattern` by sending request(s) to Jira."
        okText={step === 1 ? 'Next' : 'Save'}
        okDisabled={(step === 1 && !applicationType) || (step === 2 && !transformation.remotelinkCommitShaPattern)}
        okLoading={operating}
        cancelText={step === 2 ? 'Prev' : 'Cancel'}
        onOk={handleSubmit}
        onCancel={handleCancel}
      >
        <S.DialogBody>
          {step === 1 && (
            <>
              <FormItem
                label="Jira Issue Key"
                subLabel="Input any issue key that has connected commit(s) in the development panel"
                required
              >
                <div className="search">
                  <InputGroup
                    placeholder="Please enter..."
                    value={issueKey}
                    onChange={(e) => setIssueKey(e.target.value)}
                  />
                  <Button loading={searching} disabled={!issueKey} text="See Results" onClick={handleSearch} />
                </div>
              </FormItem>
              {applicationTypes.length > 0 && (
                <FormItem label="Application Type" subLabel="Please choose an application type." required>
                  <RadioGroup
                    selectedValue={applicationType}
                    onChange={(e) => setApplicationType((e.target as HTMLInputElement).value)}
                  >
                    {applicationTypes.map((at) => (
                      <Radio key={at} value={at} label={at} />
                    ))}
                  </RadioGroup>
                </FormItem>
              )}
            </>
          )}
          {step === 2 && (
            <>
              <FormItem label="Jira Issue Key">{issueKey}</FormItem>
              <FormItem label="Application Type">{applicationType}</FormItem>
              <FormItem label="Commit Url Preview" subLabel="The latest five commit(s) associated with the issue.">
                <ul>
                  {devPanelCommits.map((commit) => (
                    <li key={commit}>{commit}</li>
                  ))}
                </ul>
              </FormItem>
              <FormItem
                label="Commit Pattern"
                subLabel={
                  <>
                    Input pattern(s) to match and parse commits and repo identifiers from above commit URLs. See
                    examples
                  </>
                }
                required
              >
                <InputGroup
                  placeholder="eg. https://gitlab.com/{namespace}/{repo_name}/commit/{commit_sha}"
                  value={transformation.remotelinkCommitShaPattern}
                  onChange={(e) => setTransformation({ ...transformation, remotelinkCommitShaPattern: e.target.value })}
                />
              </FormItem>
            </>
          )}
        </S.DialogBody>
      </Dialog>
    </S.CrossDomain>
  );
};
