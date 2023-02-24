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
import { useLocation, useHistory } from 'react-router-dom';
import { Menu, MenuItem, Tag, Navbar, Intent, Alignment } from '@blueprintjs/core';

import { Logo } from '@/components';
import { useVersion } from '@/store';

import FileIcon from '@/images/icons/file.svg';
import GitHubIcon from '@/images/icons/github.svg';
import SlackIcon from '@/images/icons/slack.svg';

import { useMenu, MenuItemType } from './use-menu';
import * as S from './styled';

interface Props {
  children: React.ReactNode;
}

export const BaseLayout = ({ children }: Props) => {
  const menu = useMenu();
  const { pathname } = useLocation();
  const history = useHistory();
  const { version } = useVersion();

  const handlePushPath = (it: MenuItemType) => {
    if (!it.target) {
      history.push(it.path);
    } else {
      window.open(it.path, '_blank');
    }
  };

  return (
    <S.Container>
      <S.Sider>
        <Logo />
        <Menu className="menu">
          {menu.map((it) => {
            const paths = [it.path, ...(it.children ?? []).map((cit) => cit.path)];
            const active = !!paths.find((path) => pathname.includes(path));
            return (
              <MenuItem
                key={it.key}
                className="menu-item"
                text={it.title}
                icon={it.icon}
                active={active}
                onClick={() => handlePushPath(it)}
              >
                {it.children?.map((cit) => (
                  <MenuItem
                    key={cit.key}
                    className="sub-menu-item"
                    text={
                      <S.SiderMenuItem>
                        <span>{cit.title}</span>
                        {cit.isBeta && <Tag intent={Intent.WARNING}>beta</Tag>}
                      </S.SiderMenuItem>
                    }
                    icon={cit.icon ?? <img src={cit.iconUrl} width={16} alt="" />}
                    active={pathname.includes(cit.path)}
                    disabled={cit.disabled}
                    onClick={() => handlePushPath(cit)}
                  />
                ))}
              </MenuItem>
            );
          })}
        </Menu>
        <div className="copyright">
          <div>Apache 2.0 License</div>
          <div className="version">{version}</div>
        </div>
      </S.Sider>
      <S.Inner>
        <S.Header>
          <Navbar.Group align={Alignment.RIGHT}>
            <a href="https://devlake.apache.org/docs/UserManuals/ConfigUI/Tutorial/" rel="noreferrer" target="_blank">
              <img src={FileIcon} alt="documents" />
              <span>Docs</span>
            </a>
            <Navbar.Divider />
            <a
              href="https://github.com/apache/incubator-devlake"
              rel="noreferrer"
              target="_blank"
              className="navIconLink"
            >
              <img src={GitHubIcon} alt="github" />
              <span>GitHub</span>
            </a>
            <Navbar.Divider />
            <a
              href="https://join.slack.com/t/devlake-io/shared_invite/zt-17b6vuvps-x98pqseoUagM7EAmKC82xQ"
              rel="noreferrer"
              target="_blank"
            >
              <img src={SlackIcon} alt="slack" />
              <span>Slack</span>
            </a>
          </Navbar.Group>
        </S.Header>
        <S.Content>{children}</S.Content>
      </S.Inner>
    </S.Container>
  );
};
