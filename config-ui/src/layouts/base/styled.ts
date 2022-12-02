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

import styled from '@emotion/styled'
import { Navbar } from '@blueprintjs/core'

export const Container = styled.div`
  display: flex;
  height: 100vh;
  background-color: #f9f9fa;
  overflow: hidden;
`

export const Sider = styled.div`
  flex: 0 0 200px;
  position: relative;
  width: 200px;
  background-color: #292b3f;

  .logo {
    display: flex;
    text-align: center;
    flex-direction: column;
    align-items: center;
    margin: 20px 0;

    img:first-child {
      width: 48px;
      height: 48px;
    }

    img:last-child {
      margin-top: 10px;
      width: 100px;
      height: 13px;
    }
  }

  .menu {
    margin: 0 4px;
    color: #fff;
    background-color: transparent;

    .menu-item,
    .sub-menu-item {
      display: flex;
      align-items: center;
      margin: 2px 0;
      line-height: 26px;
      transition: all 0.3s ease;
      border-radius: 8px;
      outline: none;
      cursor: pointer;

      &:hover {
        background-color: rgba(167, 182, 194, 0.3);
      }

      &.bp3-active:hover {
        background-color: #7497f7;
      }

      .bp3-icon {
        svg {
          width: 12px;
          height: 12px;
        }
      }
    }

    .sub-menu-item {
      border-radius: 3px;
    }
  }

  .copyright {
    position: absolute;
    right: 0;
    bottom: 30px;
    left: 0;
    text-align: center;
    color: rgba(124, 124, 124, 0.7);
  }
`

export const Inner = styled.div`
  display: flex;
  flex-direction: column;
  flex: auto;
  height: 100vh;
  overflow: auto;
`

export const Header = styled(Navbar)`
  flex: 0 0 50px;
  background-color: #f9f9fa;
  box-shadow: none;
`

export const Content = styled.div`
  flex: auto;
  margin: 24px auto 0;
  max-width: 820px;
  width: 100%;
`

export const SlackContainer = styled.div`
  padding: 10px;
  width: 200px;
  font-size: 12px;

  & > img {
    display: block;
    margin: 0 auto;
    width: 130px;
    height: 50px;
  }

  & > p {
    text-align: center;
  }
`
