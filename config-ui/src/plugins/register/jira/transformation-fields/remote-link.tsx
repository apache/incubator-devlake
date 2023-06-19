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
import { InputGroup, Button, Intent } from '@blueprintjs/core';

import { IconButton } from '@/components';
import { operator } from '@/utils';

import * as API from '../api';
import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const RemoteLink = ({ transformation, setTransformation }: Props) => {
  const [links, setLinks] = useState<Array<{ pattern: string; regex: string }>>(
    transformation.remotelinkRepoPattern ?? [],
  );
  const [generating, setGenerating] = useState(false);

  const generateRegex = async (pattern: string) => {
    try {
      const res = await API.generateRegex(pattern);
      return { pattern, regex: res.regex };
    } catch {
      return { pattern, regex: '' };
    }
  };

  const getRegex = async () => {
    const [success, res] = await operator(
      () =>
        Promise.all(
          links.map((link) => {
            if (!link.pattern || link.regex) {
              return link;
            }
            return generateRegex(link.pattern);
          }),
        ),
      {
        setOperating: setGenerating,
        hideToast: true,
      },
    );

    if (success) {
      setTransformation({
        ...transformation,
        remotelinkRepoPattern: res.filter((it: any) => it.regex),
      });
    }
  };

  useEffect(() => {
    const timer = setTimeout(getRegex, 1000);
    return () => clearTimeout(timer);
  }, [links]);

  const handleChangeLinks = (index: number, value: string) => {
    const newValue = links.map((link, i) => (index === i ? { pattern: value, regex: link.regex } : link));
    setLinks(newValue);
  };

  const handleAddLink = () => {
    const newValue = [...links, { pattern: '', regex: '' }];
    setLinks(newValue);
  };

  const handleDeleteLink = (index: number) => {
    const newValue = links.filter((_, i) => i !== index);
    setLinks(newValue);
  };

  return (
    <S.RemoteLinkWrapper>
      {links.map((link, i) => (
        <div key={i} className="input">
          <InputGroup
            key={i}
            placeholder="E.g. https://gitlab.com/{namespace}/{repo_name}/-/commit/{commit_sha}"
            value={link.pattern}
            onChange={(e) => handleChangeLinks(i, e.target.value)}
            // onBlur={() => handleGenerateRegex(i, link.pattern)}
          />
          {links.length > 1 && (
            <IconButton disabled={generating} icon="cross" tooltip="Delete" onClick={() => handleDeleteLink(i)} />
          )}
        </div>
      ))}
      <Button
        outlined
        disabled={generating}
        intent={Intent.PRIMARY}
        icon="add"
        text="Add a Pattern"
        onClick={() => handleAddLink()}
      />
    </S.RemoteLinkWrapper>
  );
};
