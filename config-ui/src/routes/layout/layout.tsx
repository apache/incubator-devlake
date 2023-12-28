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

import { useState, useEffect, useRef } from 'react';
import { useLoaderData, Outlet, useNavigate, useLocation } from 'react-router-dom';
import { CSSTransition } from 'react-transition-group';
import { CloseOutlined } from '@ant-design/icons';
import { Layout as AntdLayout, Menu, Divider, Button } from 'antd';

import { useAppDispatch, useAppSelector } from '@/app/hook';
import { PageLoading, Logo, ExternalLink } from '@/components';
import { init, selectError, selectStatus } from '@/features';
import { TipsContextProvider, TipsContextConsumer } from '@/store';

import { loader } from './loader';
import { menuItems, menuItemsMatch, headerItems } from './config';
import * as S from './styled';
import './tips-transition.css';

const { Sider, Header, Content, Footer } = AntdLayout;

export const Layout = () => {
  const [openKeys, setOpenKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);

  const { version, plugins } = useLoaderData() as Awaited<ReturnType<typeof loader>>;

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const dispatch = useAppDispatch();
  const status = useAppSelector(selectStatus);
  const error = useAppSelector(selectError);

  const tipsRef = useRef(null);

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
    setSelectedKeys([pathname]);
  }, [pathname]);

  if (['idle', 'loading'].includes(status)) {
    return <PageLoading />;
  }

  if (status === 'failed') {
    throw error.message;
  }

  return (
    <TipsContextProvider>
      <TipsContextConsumer>
        {({ tips, setTips }) => (
          <AntdLayout style={{ minHeight: '100vh' }}>
            <Sider
              style={{
                position: 'fixed',
                top: 0,
                bottom: 0,
                left: 0,
                height: '100vh',
                overflow: 'auto',
              }}
            >
              <Logo style={{ padding: 24 }} />
              <Menu
                mode="inline"
                theme="dark"
                items={menuItems}
                openKeys={openKeys}
                selectedKeys={selectedKeys}
                onSelect={({ key }) => navigate(key)}
                onOpenChange={(keys) => setOpenKeys(keys)}
              />
              <div style={{ position: 'absolute', right: 0, bottom: 20, left: 0, color: '#fff', textAlign: 'center' }}>
                {version}
              </div>
            </Sider>
            <AntdLayout style={{ marginLeft: 200 }}>
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
                {headerItems.map((item, i) => (
                  <ExternalLink key={item.label} link={item.link} style={{ display: 'flex', alignItems: 'center' }}>
                    {item.icon}
                    <span style={{ marginLeft: 4 }}>{item.label}</span>
                    {i !== headerItems.length - 1 && <Divider type="vertical" />}
                  </ExternalLink>
                ))}
              </Header>
              <Content style={{ margin: '0 auto', width: 1188 }}>
                <Outlet />
              </Content>
              <Footer style={{ color: '#a1a1a1', textAlign: 'center' }}>
                {import.meta.env.DEVLAKE_COPYRIGHT ?? 'Apache 2.0 License'}
              </Footer>
              <CSSTransition in={!!tips} unmountOnExit timeout={300} nodeRef={tipsRef} classNames="tips">
                <S.Tips ref={tipsRef}>
                  <div className="content">{tips}</div>
                  <Button type="primary" icon={<CloseOutlined />} onClick={() => setTips('')} />
                </S.Tips>
              </CSSTransition>
            </AntdLayout>
          </AntdLayout>
        )}
      </TipsContextConsumer>
    </TipsContextProvider>
  );
};
