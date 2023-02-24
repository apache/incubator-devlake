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

import Img from '@/images/no-data.svg';
import { Card } from '@/components';

import styled from 'styled-components';

const Wrapper = styled(Card)`
  text-align: center;

  img {
    display: inline-block;
    width: 120px;
    height: 120px;
  }

  .action {
    margin-top: 24px;
  }
`;

interface Props {
  text: React.ReactNode;
  action?: React.ReactNode;
}

export const NoData = ({ text, action }: Props) => {
  return (
    <Wrapper>
      <img src={Img} alt="" />
      <p>{text}</p>
      <div className="action">{action}</div>
    </Wrapper>
  );
};
