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

import { useNavigate, Navigate } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import { CloseOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import { theme, Layout, Modal } from 'antd';

import { selectOnboard } from '@/features/onboard';
import { useAppSelector } from '@/hooks';

import { Step0 } from './step-0';
import { Step1 } from './step-1';
import { Step2 } from './step-2';
import { Step3 } from './step-3';
import { Step4 } from './step-4';
import * as S from './styled';

const steps = [
  {
    step: 1,
    title: 'Create Project',
  },
  {
    step: 2,
    title: 'Configure Connection',
  },
  {
    step: 3,
    title: 'Add data scope',
  },
];

const brandName = import.meta.env.DEVLAKE_BRAND_NAME ?? 'DevLake';

interface Props {
  logo?: React.ReactNode;
  title?: React.ReactNode;
}

export const Onboard = ({ logo, title }: Props) => {
  const navigate = useNavigate();

  const { step, done } = useAppSelector(selectOnboard);

  const {
    token: { colorPrimary },
  } = theme.useToken();

  const [modal, contextHolder] = Modal.useModal();

  const handleClose = () => {
    modal.confirm({
      width: 820,
      title: 'Are you sure to exit the onboarding session?',
      content: 'You can get back to this session via the card on top of the Projects page.',
      icon: <ExclamationCircleOutlined />,
      okText: 'Confirm',
      onOk: () => navigate('/'),
    });
  };

  if (done) {
    return <Navigate to="/projects" />;
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Helmet>
        <title>Onboard - {brandName}</title>
      </Helmet>
      <S.Inner>
        {step === 0 ? (
          <Step0 logo={logo} title={title} />
        ) : (
          <>
            <S.Header>
              <h1>Connect to your first repository</h1>
              <CloseOutlined style={{ fontSize: 18, color: '#70727F', cursor: 'pointer' }} onClick={handleClose} />
            </S.Header>
            <S.Content>
              {[1, 2, 3].includes(step) && (
                <S.Step>
                  {steps.map((it) => (
                    <S.StepItem key={it.step} $actived={it.step === step} $activedColor={colorPrimary}>
                      <span>{it.step}</span>
                      <span>{it.title}</span>
                    </S.StepItem>
                  ))}
                </S.Step>
              )}
              {step === 1 && <Step1 />}
              {step === 2 && <Step2 />}
              {step === 3 && <Step3 />}
              {step === 4 && <Step4 />}
            </S.Content>
          </>
        )}
      </S.Inner>
      {contextHolder}
    </Layout>
  );
};
