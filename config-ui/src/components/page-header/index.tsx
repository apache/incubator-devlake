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

import { Link } from 'react-router-dom';
import { RightOutlined } from '@ant-design/icons';

import * as S from './styled';

interface Props {
  breadcrumbs: Array<{
    name: string;
    path: string;
  }>;
  extra?: React.ReactNode;
  description?: string;
  children: React.ReactNode;
}

export const PageHeader = ({ breadcrumbs, extra, description, children }: Props) => {
  return (
    <S.Container>
      <S.Title>
        <S.Breadcrumbs>
          {breadcrumbs.map(({ name, path }, i, arr) => (
            <S.Breadcrumb key={i}>
              <Link to={path}>
                <span>{name}</span>
                {arr.length - 1 !== i && <RightOutlined />}
              </Link>
            </S.Breadcrumb>
          ))}
        </S.Breadcrumbs>
        <S.Extra>{extra}</S.Extra>
      </S.Title>
      <S.Description>{description}</S.Description>
      <S.Content>{children}</S.Content>
    </S.Container>
  );
};
