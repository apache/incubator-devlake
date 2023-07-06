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
import { InputGroup, Button, RadioGroup, Radio, Icon, Collapse } from '@blueprintjs/core';

import { Dialog, FormItem, toast } from '@/components';
import JiraIssueTipsImg from '@/images/jira-issue-tips.png';
import { operator } from '@/utils';

import * as API from '../api';
import * as S from './styled';

interface Props {
  connectionId: ID;
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
  isOpen: boolean;
  onCancel: () => void;
}

export const DevPanel = ({ connectionId, transformation, setTransformation, isOpen, onCancel }: Props) => {
  const [step, setStep] = useState(1);
  const [issueKey, setIssueKey] = useState('');
  const [searching, setSearching] = useState(false);
  const [applicationTypes, setApplicationTypes] = useState<string[]>([]);
  const [applicationType, setApplicationType] = useState<string>();
  const [showTip, setShowTip] = useState(false);
  const [operating, setOperating] = useState(false);
  const [devPanelCommits, setDevPanelCommits] = useState<string[]>([]);
  const [pattern, setPattern] = useState('');
  const [regex, setRegex] = useState('');
  const [preview, setPreview] = useState([]);

  const getRegex = async () => {
    if (!pattern) return;

    const [success, res] = await operator(
      async () => {
        const { regex } = await API.generateRegex(pattern);
        const preview = await API.applyRegex(regex, devPanelCommits);
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
    const [success, res] = await operator(() => API.getApplicationTypes(connectionId, { key: issueKey }), {
      setOperating: setSearching,
      hideToast: true,
    });

    if (success && res) {
      setApplicationTypes(res);
      setApplicationType(res[0]);
    } else {
      toast.error('Cannot find the Jira issue, please input the right issue key.');
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
    <Dialog
      style={{ width: 820 }}
      isOpen={isOpen}
      title="Configure the `Application Type` and `Commit Pattern` by sending request(s) to Jira."
      okText={step === 1 ? 'Next' : 'Save'}
      okDisabled={(step === 1 && !applicationType) || (step === 2 && (!pattern || !regex))}
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
                  <p style={{ display: 'flex', alignItems: 'center' }} onClick={() => setShowTip(!showTip)}>
                    Input pattern(s) to match and parse commits and repo identifiers from above commit URLs. See
                    examples <Icon icon={showTip ? 'chevron-up' : 'chevron-down'} style={{ cursor: 'pointer' }} />
                  </p>
                  <Collapse isOpen={showTip}>
                    <img src={JiraIssueTipsImg} width="100%" alt="" />
                  </Collapse>
                </>
              }
              required
            >
              <InputGroup
                placeholder="eg. https://gitlab.com/{namespace}/{repo_name}/commit/{commit_sha}"
                value={pattern}
                onChange={(e) => setPattern(e.target.value)}
              />
            </FormItem>
            <FormItem label="Configuration Results Preview">
              <code>
                <pre>{JSON.stringify(preview, null, '  ')}</pre>
              </code>
            </FormItem>
          </>
        )}
      </S.DialogBody>
    </Dialog>
  );
};
