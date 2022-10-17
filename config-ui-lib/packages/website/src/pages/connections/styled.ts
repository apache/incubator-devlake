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
import styled from '@emotion/styled';

export const Container = styled.div`
  h4 {
    margin-top: 12px;
  }

  .item {
    margin-top: 48px;
  }

  ul.list {
    margin-top: 24px;

    li {
      display: inline-block;
      width: 130px;
      text-align: center;

      & > a {
        display: inline-flex;
        flex-direction: column;
        align-items: center;
        padding: 16px 24px;
        transition: all 0.3s ease;

        &:hover {
          box-shadow: 1px 1px 6px rgb(0 0 0 / 10%);
        }

        span {
          margin-top: 6px;
          color: #777;
        }
      }
    }
  }
`;
