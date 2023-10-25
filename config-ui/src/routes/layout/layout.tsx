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

import { useEffect, useRef } from 'react';
import { useLoaderData, Outlet, useNavigate, useLocation } from 'react-router-dom';
import { CSSTransition } from 'react-transition-group';
import { Menu, MenuItem, Navbar, Alignment } from '@blueprintjs/core';

import { useAppDispatch, useAppSelector } from '@/app/hook';
import { PageLoading, Logo, ExternalLink, IconButton } from '@/components';
import { init, selectStatus } from '@/features';
import { DOC_URL } from '@/release';
import { TipsContextProvider, TipsContextConsumer } from '@/store';

import DashboardIcon from '@/images/icons/dashboard.svg';
import FileIcon from '@/images/icons/file.svg';
import APIIcon from '@/images/icons/api.svg';
import GitHubIcon from '@/images/icons/github.svg';
import SlackIcon from '@/images/icons/slack.svg';

import { loader } from './loader';
import { useMenu, MenuItemType } from './use-menu';
import * as S from './styled';
import './tips-transition.css';

export const Layout = () => {
  const { version, plugins } = useLoaderData() as Awaited<ReturnType<typeof loader>>;

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const dispatch = useAppDispatch();
  const status = useAppSelector(selectStatus);

  const menu = useMenu();

  const tipsRef = useRef(null);

  useEffect(() => {
    dispatch(init(plugins));
  }, []);

  if (['idle', 'loading'].includes(status)) {
    return <PageLoading />;
  }

  const handlePushPath = (it: MenuItemType) => {
    if (!it.target) {
      navigate(it.path);
    } else {
      window.open(it.path, '_blank');
    }
  };

  const getGrafanaUrl = () => {
    const suffix = '/d/lCO8w-pVk/homepage?orgId=1';
    const { protocol, hostname } = window.location;

    return import.meta.env.DEV ? `${protocol}//${hostname}:3002${suffix}` : `/grafana${suffix}`;
  };

  return (
    <TipsContextProvider>
      <TipsContextConsumer>
        {({ tips, setTips }) => (
          <S.Wrapper>
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
                            </S.SiderMenuItem>
                          }
                          icon={cit.icon}
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
            <S.Main>
              <S.Header>
                <Navbar.Group align={Alignment.RIGHT}>
                  <S.DashboardIcon>
                    <ExternalLink link={getGrafanaUrl()}>
                      <img src={DashboardIcon} alt="dashboards" />
                      <span>Dashboards</span>
                    </ExternalLink>
                  </S.DashboardIcon>
                  <Navbar.Divider />
                  <ExternalLink link={DOC_URL.TUTORIAL}>
                    <img src={FileIcon} alt="documents" />
                    <span>Docs</span>
                  </ExternalLink>
                  <Navbar.Divider />
                  <ExternalLink link="/api/swagger/index.html">
                    <img src={APIIcon} alt="api" />
                    <span>API</span>
                  </ExternalLink>
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
              <S.Inner>
                <S.Content>
                  <Outlet />
                </S.Content>
              </S.Inner>
              <CSSTransition in={!!tips} unmountOnExit timeout={300} nodeRef={tipsRef} classNames="tips">
                <S.Tips ref={tipsRef}>
                  <div className="content">{tips}</div>
                  <IconButton style={{ color: '#fff' }} icon="cross" tooltip="Close" onClick={() => setTips('')} />
                </S.Tips>
              </CSSTransition>
            </S.Main>
          </S.Wrapper>
        )}
      </TipsContextConsumer>
    </TipsContextProvider>
  );
};
