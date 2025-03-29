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

import { useEffect, useMemo } from 'react';
import { useParams, useNavigate, useLocation, Link, Outlet } from 'react-router-dom';
import { RollbackOutlined } from '@ant-design/icons';
import { Layout, Menu } from 'antd';

import { PageHeader, PageLoading } from '@/components';
import { request, selectProjectStatus, selectProject } from '@/features/project';
import { useAppDispatch, useAppSelector } from '@/hooks';

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

const breadcrumbs = (paths: string[]) => {
  const map: Record<
    string,
    {
      path: string;
      name: string;
    }
  > = {
    '/config': {
      path: '/',
      name: 'Configurations',
    },
    projects: {
      path: '/projects',
      name: 'Projects',
    },
  };

  return paths
    .filter((p) => p)
    .map(
      (p) =>
        map[p] ?? {
          path: `/projects/${p}`,
          name: p,
        },
    );
};

export const ProjectLayout = () => {
  const { pname } = useParams() as { pname: string };
  const navigate = useNavigate();
  const { pathname } = useLocation();

  const dispatch = useAppDispatch();
  const status = useAppSelector(selectProjectStatus);
  const project = useAppSelector(selectProject);

  useEffect(() => {
    dispatch(request(pname));
  }, [pname]);

  const { paths, selectedKeys, title } = useMemo(() => {
    const paths = pathname.split('/');
    const key = paths.pop();
    const item = items.find((i) => i.key === key);

    return {
      paths,
      selectedKeys: key ? [key] : [],
      title: [
        {
          name: item?.label ?? '',
          path: '',
        },
      ],
    };
  }, [pathname]);

  if (status === 'loading' || !project) {
    return <PageLoading />;
  }

  return (
    <Layout style={{ height: '100%', overflow: 'hidden' }}>
      <Sider width={240} style={{ padding: '36px 12px', backgroundColor: '#F9F9FA', borderRight: '1px solid #E7E9F3' }}>
        <S.Top onClick={() => navigate('/projects')}>
          <RollbackOutlined />
          <span className="back">Back to Projects</span>
        </S.Top>
        <ProjectSelector name={pname} />
        <Menu
          mode="inline"
          style={{ backgroundColor: '#F9F9FA', border: 'none' }}
          items={items}
          selectedKeys={selectedKeys}
          onClick={({ key }) => navigate(`/projects/${encodeURIComponent(pname)}/${key}`)}
        />
      </Sider>
      <Layout>
        <Content style={{ padding: '36px 48px', overflowY: 'auto' }}>
          <p>
            {breadcrumbs(paths).map((b, i) => (
              <span key={b.path}>
                {i !== paths.length - 2 ? <Link to={b.path}>{b.name}</Link> : <span>{decodeURIComponent(b.name)}</span>}
                <span> / </span>
              </span>
            ))}
          </p>
          <PageHeader breadcrumbs={title}>
            <Outlet />
          </PageHeader>
        </Content>
      </Layout>
    </Layout>
  );
};
