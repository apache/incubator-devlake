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

import { useState } from 'react';
import { Switch, InputGroup } from '@blueprintjs/core';

import { ExternalLink, HelpTooltip } from '@/components';

import * as S from './styled';

interface Props {
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const AdditionalSettings = ({ transformation, setTransformation }: Props) => {
  const [enable, setEnable] = useState(true);

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (!checked) {
      setTransformation({
        ...transformation,
        refdiff: null,
      });
    }

    setEnable(checked);
  };

  return (
    <S.AdditionalSettings>
      <h2>
        <span>Additional Settings</span>
        <Switch alignIndicator="right" inline checked={enable} onChange={handleChange} />
      </h2>
      {enable && (
        <>
          <p>
            Enable the <ExternalLink link="https://devlake.apache.org/docs/Plugins/refdiff">RefDiff</ExternalLink>{' '}
            plugin to pre-calculate version-based metrics
            <HelpTooltip content="Calculate the commits diff between two consecutive tags that match the following RegEx. Issues closed by PRs which contain these commits will also be calculated. The result will be shown in table.refs_commits_diffs and table.refs_issues_diffs." />
          </p>
          <div className="refdiff">
            Compare the last
            <InputGroup
              style={{ width: 60 }}
              placeholder="10"
              value={transformation.refdiff?.tagsLimit ?? ''}
              onChange={(e) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsLimit: e.target.value,
                  },
                })
              }
            />
            tags that match the
            <InputGroup
              style={{ width: 200 }}
              placeholder="(regex)$"
              value={transformation.refdiff?.tagsPattern ?? ''}
              onChange={(e) =>
                setTransformation({
                  ...transformation,
                  refdiff: {
                    ...transformation?.refdiff,
                    tagsPattern: e.target.value,
                  },
                })
              }
            />
            for calculation
          </div>
        </>
      )}
    </S.AdditionalSettings>
  );
};
