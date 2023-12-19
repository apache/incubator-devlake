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

import { useEffect, useState } from 'react';
import { Modal, Radio, Input, Button, Collapse, message } from 'antd';

import API from '@/api';
import { Block } from '@/components';
import JiraIssueTipsImg from '@/images/jira-issue-tips.png';
import { operator } from '@/utils';

import * as S from './styled';

interface Props {
  connectionId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
  open: boolean;
  onCancel: () => void;
}

export const DevPanel = ({ connectionId, transformation, setTransformation, open, onCancel }: Props) => {
  const [step, setStep] = useState(1);
  const [issueKey, setIssueKey] = useState('');
  const [searching, setSearching] = useState(false);
  const [applicationTypes, setApplicationTypes] = useState<string[]>([]);
  const [applicationType, setApplicationType] = useState<string>();
  const [operating, setOperating] = useState(false);
  const [devPanelCommits, setDevPanelCommits] = useState<string[]>([]);
  const [pattern, setPattern] = useState('');
  const [regex, setRegex] = useState('');
  const [preview, setPreview] = useState([]);

  const getRegex = async () => {
    if (!pattern) return;

    const [success, res] = await operator(
      async () => {
        const { regex } = await API.plugin.jira.generateRegex(pattern);
        const preview = await API.plugin.jira.applyRegex(regex, devPanelCommits);
        return {
          regex,
          preview,
        };
      },
      {
        setOperating,
        hideToast: true,
      },
    );

    if (success) {
      setRegex(res.regex);
      setPreview(res.preview);
      setTransformation({
        ...transformation,
        applicationType,
        remotelinkRepoPattern: [{ pattern, regex: res.regex }],
      });
    }
  };

  useEffect(() => {
    const timer = setTimeout(getRegex, 500);
    return () => clearTimeout(timer);
  }, [pattern]);

  const handleSearch = async () => {
    const [success, res] = await operator(() => API.plugin.jira.applicationTypes(connectionId, { key: issueKey }), {
      setOperating: setSearching,
      hideToast: true,
    });

    if (success && res) {
      setApplicationTypes(res);
      setApplicationType(res[0]);
    } else {
      message.error('Cannot find the Jira issue, please input the right issue key.');
    }
  };

  const handleSubmit = async () => {
    if (step === 1 && applicationType) {
      const [success, res] = await operator(
        () => API.plugin.jira.devPanelCommits(connectionId, { key: issueKey, applicationType }),
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

    onCancel();
  };

  const handleCancel = () => {
    if (step === 1) {
      onCancel();
      return;
    }

    setStep(1);
  };

  return (
    <Modal
      open={open}
      width={820}
      centered
      title="Configure the `Application Type` and `Commit Pattern` by sending request(s) to Jira."
      okText={step === 1 ? 'Next' : 'Save'}
      okButtonProps={{
        disabled: (step === 1 && !applicationType) || (step === 2 && (!pattern || !regex)),
        loading: operating,
      }}
      cancelText={step === 2 ? 'Prev' : 'Cancel'}
      onOk={handleSubmit}
      onCancel={handleCancel}
    >
      <S.DialogBody>
        {step === 1 && (
          <>
            <Block
              title="Jira Issue Key"
              description="Input any issue key that has connected commit(s) in the development panel"
              required
            >
              <div className="search">
                <Input placeholder="Please enter..." value={issueKey} onChange={(e) => setIssueKey(e.target.value)} />
                <Button loading={searching} disabled={!issueKey} onClick={handleSearch}>
                  See Results
                </Button>
              </div>
            </Block>
            {applicationTypes.length > 0 && (
              <Block title="Application Type" description="Please choose an application type." required>
                <Radio.Group value={applicationType} onChange={(e) => setApplicationType(e.target.value)}>
                  {applicationTypes.map((at) => (
                    <Radio key={at} value={at}>
                      {at}
                    </Radio>
                  ))}
                </Radio.Group>
              </Block>
            )}
          </>
        )}
        {step === 2 && (
          <>
            <Block title="Jira Issue Key">{issueKey}</Block>
            <Block title="Application Type">{applicationType}</Block>
            <Block title="Commit Url Preview" description="The latest five commit(s) associated with the issue.">
              <ul>
                {devPanelCommits.map((commit) => (
                  <li key={commit}>{commit}</li>
                ))}
              </ul>
            </Block>
            <Block
              title="Commit Pattern"
              description={
                <Collapse
                  ghost
                  expandIconPosition="end"
                  items={[
                    {
                      key: '1',
                      label:
                        'Input pattern(s) to match and parse commits and repo identifiers from above commit URLs. See examples',
                      children: <img src={JiraIssueTipsImg} width="100%" alt="" />,
                    },
                  ]}
                />
              }
              required
            >
              <Input
                placeholder="eg. https://gitlab.com/{namespace}/{repo_name}/commit/{commit_sha}"
                value={pattern}
                onChange={(e) => setPattern(e.target.value)}
              />
            </Block>
            <Block title="Configuration Results Preview">
              <code>
                <pre>{JSON.stringify(preview, null, '  ')}</pre>
              </code>
            </Block>
          </>
        )}
      </S.DialogBody>
    </Modal>
  );
};
