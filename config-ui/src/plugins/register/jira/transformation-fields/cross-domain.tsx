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
import { Radio, Icon, Collapse, InputGroup, Button, Intent } from '@blueprintjs/core';

import { ExternalLink, IconButton } from '@/components';
import JiraIssueTipsImg from '@/images/jira-issue-tips.png';

import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const CrossDomain = ({ transformation, setTransformation }: Props) => {
  const [radio, setRadio] = useState<'repo' | 'commitSha'>('repo');
  const [repoTips, setRepoTips] = useState(false);
  const [repoLinks, setRepoLinks] = useState([]);

  // console.log(transformation);

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

  return (
    <S.CrossDomain>
      <h2>Cross-domain</h2>
      <p>
        Connect `commits` and `issues` to measure metrics such as{' '}
        <ExternalLink link="https://devlake.apache.org/docs/Metrics/BugCountPer1kLinesOfCode">
          Bug Count per 1k Lines of Code
        </ExternalLink>{' '}
        or man hour distribution on different work types. Connect `commits` and `issues` to measure metrics such as{' '}
      </p>
      <div className="radio">
        <div className="radio-item">
          <Radio checked={radio === 'repo'} onChange={() => handleChangeRadio('repo')} />
          <div className="content">
            <h5>Connect Jira issues and commits via Jira issues’ remote links that match the following pattern</h5>
            <p onClick={handleToggleRepoTips}>
              The default pattern shows how to match and parse GitLab(Cloud) commits from issue remote links. See More{' '}
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
                      placeholder="https://gitlab.com/{namespace}/{repo_name}/-/commit/{commit_sha}"
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
            <h5>Connect Jira issues and commits via Jira’s development panel</h5>
            <p>
              Choose this if you’ve enabled{' '}
              <ExternalLink link="https://support.atlassian.com/jira-software-cloud/docs/view-development-information-for-an-issue/">
                issue’ development panel
              </ExternalLink>
              . Usually, it happens when you’re using BitBucket for source code management.
            </p>
            {radio === 'commitSha' && (
              <InputGroup
                fill
                placeholder="/commit/([0-9a-f]{40})$"
                value={transformation.remotelinkCommitShaPattern ?? ''}
                onChange={(e) =>
                  setTransformation({
                    ...transformation,
                    remotelinkCommitShaPattern: e.target.value,
                  })
                }
              />
            )}
          </div>
        </div>
      </div>
    </S.CrossDomain>
  );
};
