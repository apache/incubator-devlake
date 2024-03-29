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

import { useState, useContext } from 'react';
import { Button } from 'antd';
import styled from 'styled-components';

import API from '@/api';
import { operator } from '@/utils';

import { Context } from './context';

const Wrapper = styled.div`
  margin-top: 100px;
  text-align: center;

  .welcome {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 30px 0;
    font-size: 40px;

    span.line {
      display: inline-block;
      width: 242px;
      height: 1px;
      background-color: #dbdcdf;
    }

    span.content {
      margin: 0 24px;
    }
  }

  .title {
    margin: 15px 0;
    padding: 16px 0;
    font-size: 60px;
    font-weight: 600;
  }

  .subTitle {
    font-size: 16px;
  }

  .action {
    margin: 0 auto;
    width: 200px;
    margin-top: 64px;
  }
`;

export const Step0 = () => {
  const [operating, setOperating] = useState(false);

  const { step, records, done, projectName, plugin, setStep } = useContext(Context);

  const handleSubmit = async () => {
    const [success] = await operator(
      async () => API.store.set('onboard', { step: 1, records, done, projectName, plugin }),
      {
        setOperating,
        hideToast: true,
      },
    );

    if (success) {
      setStep(step + 1);
    }
  };

  return (
    <Wrapper>
      <div className="welcome">
        <span className="line" />
        <span className="content">Welcome</span>
        <span className="line" />
      </div>
      <div className="title">Connect to your first repository</div>
      <div className="subTitle">
        Integrate your first Git tool and observe engineering metrics with just a few clicks.
      </div>
      <div className="action">
        <Button block size="large" type="primary" loading={operating} onClick={handleSubmit}>
          Start
        </Button>
      </div>
    </Wrapper>
  );
};
