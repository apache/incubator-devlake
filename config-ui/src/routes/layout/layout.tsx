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

import { useState, useEffect, useMemo } from 'react';
import { useLoaderData, Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import { Layout as AntdLayout, Menu, Divider } from 'antd';

import { PageLoading, Logo, ExternalLink } from '@/components';
import { init, selectError, selectStatus } from '@/features';
import { OnboardCard } from '@/routes/onboard/components';
import { useAppDispatch, useAppSelector } from '@/hooks';

import { menuItems, menuItemsMatch, headerItems } from './config';

const { Sider, Header, Content, Footer } = AntdLayout;

const brandName = import.meta.env.DEVLAKE_BRAND_NAME ?? 'DevLake';

export const Layout = () => {
  const [openKeys, setOpenKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);

  const { version, plugins } = useLoaderData() as { version: string; plugins: string[] };

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const dispatch = useAppDispatch();
  const status = useAppSelector(selectStatus);
  const error = useAppSelector(selectError);

  useEffect(() => {
    dispatch(init(plugins));
  }, []);

  useEffect(() => {
    const curMenuItem = menuItemsMatch[pathname];
    const parentKey = curMenuItem?.parentKey;
    if (parentKey) {
      setOpenKeys([parentKey]);
    }
  }, []);

  useEffect(() => {
    const selectedKeys = pathname.split('/').reduce((acc, cur, i, arr) => {
      if (i === 0) {
        acc.push('/');
        return acc;
      } else {
        acc.push(`${arr.slice(0, i + 1).join('/')}`);
        return acc;
      }
    }, [] as string[]);
    setSelectedKeys(selectedKeys);
  }, [pathname]);

  const title = useMemo(() => {
    const curMenuItem = menuItemsMatch[pathname];
    return curMenuItem?.label ?? '';
  }, [pathname]);

  if (['idle', 'loading'].includes(status)) {
    return <PageLoading />;
  }

  if (status === 'failed') {
    throw error.message;
  }

  return (
    <AntdLayout style={{ height: '100vh' }}>
      <Helmet>
        <title>
          {title ? `${title} - ` : ''}
          {brandName}
        </title>
      </Helmet>
      <Sider>
        {import.meta.env.DEVLAKE_TITLE_CUSTOM ? (
          <h2 style={{ margin: '36px 0', textAlign: 'center', color: '#fff' }}>
            {import.meta.env.DEVLAKE_TITLE_CUSTOM}
          </h2>
        ) : (
          <Logo style={{ padding: 24 }} />
        )}
        <Menu
          mode="inline"
          theme="dark"
          items={menuItems}
          openKeys={openKeys}
          selectedKeys={selectedKeys}
          onClick={({ key }) => navigate(key)}
          onOpenChange={(keys) => setOpenKeys(keys)}
        />
        <div style={{ position: 'absolute', right: 0, bottom: 20, left: 0, color: '#fff', textAlign: 'center' }}>
          {version}
        </div>
      </Sider>
      <AntdLayout>
        <Header
          style={{
            display: 'flex',
            justifyContent: 'flex-end',
            alignItems: 'center',
            padding: '0 24px',
            height: 50,
            background: 'transparent',
          }}
        >
          {headerItems
            .filter((item) =>
              import.meta.env.DEVLAKE_COPYRIGHT_HIDE ? !['GitHub', 'Slack'].includes(item.label) : true,
            )
            .map((item, i, arr) => (
              <ExternalLink key={item.label} link={item.link} style={{ display: 'flex', alignItems: 'center' }}>
                {item.icon}
                <span style={{ marginLeft: 4 }}>{item.label}</span>
                {i !== arr.length - 1 && <Divider type="vertical" />}
              </ExternalLink>
            ))}
        </Header>
        <Content style={{ overflowY: 'auto' }}>
          <div style={{ padding: 24, margin: '0 auto', maxWidth: 1280 }}>
            <OnboardCard style={{ marginBottom: 32 }} />
            <Outlet />
          </div>
          {!import.meta.env.DEVLAKE_COPYRIGHT_HIDE && (
            <Footer>
              <p style={{ textAlign: 'center' }}>Apache 2.0 License</p>
            </Footer>
          )}
        </Content>
      </AntdLayout>
    </AntdLayout>
  );
};
