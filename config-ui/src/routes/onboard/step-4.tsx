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

import { SmileFilled } from '@ant-design/icons';
import { theme } from 'antd';
import styled from 'styled-components';

import * as S from './styled';

const Top = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 100px;
  margin-bottom: 24px;
  height: 70px;

  span.text {
    margin-left: 8px;
    font-size: 20px;
  }
`;

export const Step4 = () => {
  const {
    token: { green5 },
  } = theme.useToken();

  return (
    <>
      <Top>
        <SmileFilled style={{ fontSize: 36, color: green5 }} />
        <span className="text">Congratulations！You have successfully connected to your first repository!</span>
      </Top>
      <S.StepContent style={{ padding: 24 }}></S.StepContent>
    </>
  );
};
