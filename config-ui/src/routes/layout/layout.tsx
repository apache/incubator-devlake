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
import { Layout as AntdLayout, Menu, Flex, Divider, Tooltip, Button } from 'antd';

import API from '@/api';
import { PageLoading, Logo, ExternalLink, Message } from '@/components';
import { PATHS } from '@/config';
import {
  init,
  selectError,
  selectStatus,
  selectTipsShow,
  selectTipsType,
  selectTipsPayload,
  hideTips,
} from '@/features';
import { useAppDispatch, useAppSelector } from '@/hooks';
import { operator } from '@/utils';

import { layoutLoader } from './loader';
import { menuItems, menuItemsMatch, headerItems } from './config';
import * as S from './styled';
import './tips-transition.css';

const { Sider, Header, Content, Footer } = AntdLayout;

export const Layout = () => {
  const [openKeys, setOpenKeys] = useState<string[]>([]);
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [operating, setOperating] = useState(false);

  const { version, plugins } = useLoaderData() as Awaited<ReturnType<typeof layoutLoader>>;

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const dispatch = useAppDispatch();
  const status = useAppSelector(selectStatus);
  const error = useAppSelector(selectError);
  const tipsShow = useAppSelector(selectTipsShow);
  const tipsType = useAppSelector(selectTipsType);
  const tipsPayload = useAppSelector(selectTipsPayload);

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

  if (['idle', 'loading'].includes(status)) {
    return <PageLoading />;
  }

  if (status === 'failed') {
    throw error.message;
  }

  const handleRunBP = async () => {
    if (!tipsPayload) {
      return;
    }

    const { blueprintId, pname } = tipsPayload;

    const [success] = await operator(
      () => API.blueprint.trigger(tipsPayload.blueprintId, { skipCollectors: false, fullSync: false }),
      {
        setOperating,
        formatMessage: () => 'Trigger blueprint successful.',
      },
    );

    if (success) {
      navigate(pname ? PATHS.PROJECT(pname) : PATHS.BLUEPRINT(blueprintId));
    }
  };

  return (
    <AntdLayout style={{ height: '100vh' }}>
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
          <div style={{ padding: 24 }}>
            <Outlet />
          </div>
          {!import.meta.env.DEVLAKE_COPYRIGHT_HIDE && (
            <Footer>
              <p style={{ textAlign: 'center' }}>Apache 2.0 License</p>
            </Footer>
          )}
        </Content>
        <CSSTransition in={!!tipsShow} unmountOnExit timeout={300} nodeRef={tipsRef} classNames="tips">
          <S.Tips ref={tipsRef}>
            <div className="content">
              {tipsType === 'data-scope-changed' && (
                <Flex gap="middle">
                  <Message content="The change of Data Scope(s) will affect the metrics of this project. Would you like to recollect the data to get them updated?" />
                  <Button type="primary" loading={operating} onClick={handleRunBP}>
                    Recollect Data
                  </Button>
                </Flex>
              )}
              {tipsType === 'scope-config-changed' && (
                <Message
                  content="Scope Config(s) have been updated. If you would like to re-transform or re-collect the data in the related
              project(s), please go to the Project page and do so."
                />
              )}
            </div>
            <Tooltip title="Close">
              <Button shape="circle" ghost icon={<CloseOutlined />} onClick={() => dispatch(hideTips())} />
            </Tooltip>
          </S.Tips>
        </CSSTransition>
      </AntdLayout>
    </AntdLayout>
  );
};
