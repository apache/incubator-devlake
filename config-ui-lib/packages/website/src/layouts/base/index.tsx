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
import React, { useMemo, useState } from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import type { MenuProps } from 'antd';
import { Layout, Menu, Popover } from 'antd';
import {
  DatabaseOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  LinkOutlined,
  MailOutlined,
  SlackOutlined,
} from '@ant-design/icons';

import * as S from './styled';

import Logo from '@/images/logo.svg';
import Textmark from '@/images/textmark.svg';
import SlackLogo from '@/images/slack-logo.svg';

const { Header, Sider, Content, Footer } = Layout;

const items = [
  {
    label: <Link to="/connections">Connections</Link>,
    key: 'connections',
    icon: <DatabaseOutlined />,
  },
];

export const Base = () => {
  const { pathname } = useLocation();
  const [collapsed, setCollapsed] = useState(false);

  const handleChangeMenu: MenuProps['onClick'] = (e) => {
    console.log(e);
  };

  const selectedKeys = useMemo(() => {
    return pathname.split('/')[1];
  }, [pathname]);

  return (
    <S.Container>
      <Sider width={220} trigger={null} collapsible collapsed={collapsed}>
        <div className="logo">
          <img src={Logo} alt="" />
          {!collapsed && <img src={Textmark} alt="" />}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          items={items}
          selectedKeys={[selectedKeys]}
          onClick={handleChangeMenu}
        />
      </Sider>
      <Layout>
        <Header>
          {React.createElement(
            collapsed ? MenuUnfoldOutlined : MenuFoldOutlined,
            {
              className: 'trigger',
              onClick: () => setCollapsed(!collapsed),
            },
          )}

          <ul className="other-info">
            <li>
              <a
                href="https://github.com/apache/incubator-devlake"
                rel="noreferrer"
                target="_blank"
              >
                <LinkOutlined style={{ fontSize: 16 }} />
              </a>
            </li>
            <li>
              <a
                href="mailto:hello@merico.dev"
                rel="noreferrer"
                target="_blank"
              >
                <MailOutlined style={{ fontSize: 16 }} />
              </a>
            </li>
            <li>
              <Popover
                placement="bottomLeft"
                content={
                  <div style={{ width: 240, textAlign: 'center' }}>
                    <img src={SlackLogo} width={130} alt="" />
                    <p>
                      Want to interact with the{' '}
                      <strong>Merico Community</strong>? Join us on our Slack
                      Channel.
                      <br />
                      <a
                        href="https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ"
                        rel="noreferrer"
                        target="_blank"
                        className="bp3-button bp3-intent-warning bp3-elevation-1 bp3-small"
                        style={{ marginTop: '10px' }}
                      >
                        Message us on&nbsp;<strong>Slack</strong>
                      </a>
                    </p>
                  </div>
                }
              >
                <SlackOutlined style={{ fontSize: 16 }} />
              </Popover>
            </li>
          </ul>
        </Header>
        <Content>
          <Outlet />
        </Content>
        <Footer>
          <div className="copyright">Apache 2.0 License</div>
        </Footer>
      </Layout>
    </S.Container>
  );
};
