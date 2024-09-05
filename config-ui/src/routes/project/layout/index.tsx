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

import { useMemo } from 'react';
import { useParams, useNavigate, useLocation, Outlet } from 'react-router-dom';
import { RollbackOutlined } from '@ant-design/icons';
import { Layout, Menu } from 'antd';

import { PageHeader } from '@/components';
import { PATHS } from '@/config';

import { ProjectSelector } from './project-selector';
import * as S from './styled';

const { Sider, Content } = Layout;

const items = [
  {
    key: 'general-settings',
    label: 'General Settings',
    style: {
      paddingLeft: 8,
    },
  },
  {
    key: 'webhooks',
    label: 'Webhooks',
    style: {
      paddingLeft: 8,
    },
  },
  {
    key: 'additional-settings',
    label: 'Additional Settings',
    style: {
      paddingLeft: 8,
    },
  },
];

export const ProjectLayout = () => {
  const { pname } = useParams() as { pname: string };
  const navigate = useNavigate();
  const { pathname } = useLocation();

  const { selectedKeys, breadcrumbs } = useMemo(() => {
    const key = pathname.split('/').pop();
    const item = items.find((i) => i.key === key);

    return {
      selectedKeys: key ? [key] : [],
      breadcrumbs: [
        {
          name: item?.label ?? '',
          path: '',
        },
      ],
    };
  }, [pathname]);

  return (
    <Layout style={{ height: '100%', overflow: 'hidden' }}>
      <Sider width={240} style={{ padding: '36px 12px', backgroundColor: '#F9F9FA', borderRight: '1px solid #E7E9F3' }}>
        <S.Top onClick={() => navigate(PATHS.PROJECTS())}>
          <RollbackOutlined />
          <span className="back">Back to Projects</span>
        </S.Top>
        <ProjectSelector name={pname} />
        <Menu
          mode="inline"
          style={{ backgroundColor: '#F9F9FA', border: 'none' }}
          items={items}
          selectedKeys={selectedKeys}
          onClick={({ key }) => navigate(`${PATHS.PROJECT(pname)}/${key}`)}
        />
      </Sider>
      <Layout>
        <Content style={{ padding: '36px 48px', overflowY: 'auto' }}>
          <p>Configurations / Projects / {pname} /</p>
          <PageHeader breadcrumbs={breadcrumbs}>
            <Outlet />
          </PageHeader>
        </Content>
      </Layout>
    </Layout>
  );
};
