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

import styled from 'styled-components';

const Wrapper = styled.div<{ position: 'top' | 'bottom'; align: 'left' | 'right' | 'center' }>`
  display: flex;
  align-items: center;

  ${({ position }) => {
    if (position === 'top') {
      return 'margin-bottom: 24px;';
    }

    if (position === 'bottom') {
      return 'margin-top: 24px;';
    }
  }}

  ${({ align }) => {
    if (align === 'left') {
      return 'justify-content: flex-start;';
    }

    if (align === 'right') {
      return 'justify-content: flex-end;';
    }

    if (align === 'center') {
      return 'justify-content: space-around;';
    }
  }}

  .bp4-button + .bp4-button {
    margin-left: 8px;
  }
`;

interface Props {
  position?: 'top' | 'bottom';
  align?: 'left' | 'right' | 'center';
  children: React.ReactNode;
}

export const Buttons = ({ position = 'top', align = 'left', children }: Props) => {
  return (
    <Wrapper position={position} align={align}>
      {children}
    </Wrapper>
  );
};
