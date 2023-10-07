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
import { useDebounce } from 'ahooks';

import { IconButton } from '@/components';
import { operator } from '@/utils';

import * as API from '../api';
import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const RemoteLink = ({ transformation, setTransformation }: Props) => {
  const [index, setInedx] = useState<number>();
  const [pattern, setPattern] = useState('');
  const [error, setError] = useState('');
  const [links, setLinks] = useState<Array<{ pattern: string; regex: string }>>(
    transformation.remotelinkRepoPattern ?? [],
  );
  const [generating, setGenerating] = useState(false);

  useEffect(() => {
    setTransformation({
      ...transformation,
      remotelinkRepoPattern: links.filter((link) => link.regex),
    });
  }, [links]);

  const debouncedPattern = useDebounce(pattern, { wait: 500 });

  const getRegex = async () => {
    const [success, res] = await operator(() => API.generateRegex(pattern), {
      hideToast: true,
      setOperating: setGenerating,
    });

    if (success) {
      setLinks(links.map((link, i) => (i === index ? { ...res, pattern } : link)));
    } else {
      setError(res?.response?.data?.message ?? '');
    }
  };

  useEffect(() => {
    if (debouncedPattern) {
      getRegex();
    }
  }, [debouncedPattern]);

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
          <div className="inner">
            <InputGroup
              key={i}
              placeholder="E.g. https://gitlab.com/{namespace}/{repo_name}/-/commit/{commit_sha}"
              value={index === i ? pattern : link.pattern}
              onChange={(e) => {
                setPattern(e.target.value);
                setError('');
              }}
              onFocus={() => {
                setInedx(i);
                setPattern(link.pattern);
                setError('');
              }}
            />
            {links.length > 1 && (
              <IconButton loading={generating} icon="cross" tooltip="Delete" onClick={() => handleDeleteLink(i)} />
            )}
          </div>
          {index === i && error && <div className="error">{error}</div>}
        </div>
      ))}
      <Button
        outlined
        loading={generating}
        intent={Intent.PRIMARY}
        icon="add"
        text="Add a Pattern"
        onClick={() => handleAddLink()}
      />
    </S.RemoteLinkWrapper>
  );
};
