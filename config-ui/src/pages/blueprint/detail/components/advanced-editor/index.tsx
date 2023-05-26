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

import React from 'react';
import { TextArea, ButtonGroup, Button, Menu, MenuItem, Position } from '@blueprintjs/core';
import { Popover2 } from '@blueprintjs/popover2';
import styled from 'styled-components';

import { EXAMPLE_CONFIG } from './example';

const Wrapper = styled.div`
  h2 {
    margin: 0;
    padding: 0;
    font-size: 16px;
    font-weight: 600;
  }

  h3 {
    margin: 0 0 8px;
    padding: 0;
    font-size: 14px;
    font-weight: 600;
  }

  p {
    margin: 0 0 8px;
  }

  textarea {
    margin-bottom: 8px;
    min-height: 240px;
    font-size: 12px;
    font-family: 'JetBrains Mono', monospace;
  }
`;

interface Props {
  value: string;
  onChange: (value: string) => void;
}

export const AdvancedEditor = ({ value, onChange }: Props) => {
  return (
    <Wrapper>
      <h3>Task Editor</h3>
      <p>
        <span>Enter JSON Configuration or preload from a template.</span>
        <a
          href="https://devlake.apache.org/docs/UserManuals/ConfigUI/AdvancedMode/#examples"
          rel="noreferrer"
          target="_blank"
        >
          See examples
        </a>
      </p>
      <TextArea fill value={value} onChange={(e) => onChange(e.target.value)} />
      <ButtonGroup minimal>
        <Button small text="Reset" icon="eraser" onClick={() => onChange(JSON.stringify([[]], null, '  '))} />
        <Popover2
          placement={Position.TOP}
          content={
            <Menu>
              {EXAMPLE_CONFIG.map((it) => (
                <MenuItem
                  key={it.id}
                  icon="code"
                  text={it.name}
                  onClick={() => onChange(JSON.stringify(it.config, null, '  '))}
                />
              ))}
            </Menu>
          }
          renderTarget={({ isOpen, ref, ...targetProps }) => (
            <Button {...targetProps} elementRef={ref} small text="Load Templates" rightIcon="caret-down" />
          )}
        />
      </ButtonGroup>
    </Wrapper>
  );
};
