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

import { Icon } from '@blueprintjs/core';

import * as S from './styled';

export interface PaginationProps {
  page: number;
  pageSize?: number;
  total: number;
  onChange: (page: number) => void;
}

export const Pagination = ({ page, pageSize = 20, total, onChange }: PaginationProps) => {
  const lastPage = Math.ceil(total / pageSize);

  const handlePrevPage = () => {
    if (page === 1) return;
    onChange(page - 1);
  };

  const handleNextPage = () => {
    if (page === lastPage) return;
    onChange(page + 1);
  };

  return (
    <S.List>
      <S.Item disabled={page === 1} onClick={handlePrevPage}>
        <Icon icon="chevron-left" />
      </S.Item>
      {Array.from({ length: lastPage }).map((_, i) => (
        <S.Item key={i + 1} active={page === i + 1} onClick={() => onChange(i + 1)}>
          <span>{i + 1}</span>
        </S.Item>
      ))}
      <S.Item disabled={page === lastPage} onClick={handleNextPage}>
        <Icon icon="chevron-right" />
      </S.Item>
    </S.List>
  );
};
