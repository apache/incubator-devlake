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

import { useState, useEffect, useReducer, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { CaretDownOutlined } from '@ant-design/icons';
import { Input, Button, Flex } from 'antd';
import { useDebounce } from 'ahooks';

import { PageLoading } from '@/components';
import { PATHS } from '@/config';
import { useOutsideClick } from '@/hooks';
import { operator } from '@/utils';

import API from '@/api';

import * as S from './styled';

type StateType = { name: string }[];

const reducer = (state: StateType, action: { type: string; payload: StateType }) => {
  switch (action.type) {
    case 'RESET':
      return [...action.payload];
    case 'APPEND':
      return [...state, ...action.payload];
    default:
      return state;
  }
};

interface Props {
  name: string;
}

export const ProjectSelector = ({ name }: Props) => {
  const [open, setOpen] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const [operating, setOperating] = useState(false);

  const [state, dispatch] = useReducer(reducer, []);

  const ref = useRef(null);

  useOutsideClick(ref, () => setOpen(false));

  const navigate = useNavigate();

  const keywordDebounce = useDebounce(keyword, { wait: 500 });

  useEffect(() => {
    setOperating(true);
    (async () => {
      const res = await API.project.list({ page: 1, pageSize, keyword: keywordDebounce });
      dispatch({ type: 'RESET', payload: res.projects });
      setTotal(res.count);
      setOperating(false);
    })();
  }, [keywordDebounce]);

  const handleAppend = async () => {
    const [success, res] = await operator(() => API.project.list({ page: page + 1, pageSize, keyword }), {
      hideToast: true,
      setOperating,
    });

    if (success) {
      setPage(page + 1);
      dispatch({ type: 'APPEND', payload: res.projects });
    }
  };

  return (
    <S.ProjectSelector ref={ref}>
      <h1 onClick={() => setOpen(true)}>
        <span>{name}</span>
        <CaretDownOutlined style={{ fontSize: 12 }} />
      </h1>
      {open && (
        <S.Selector>
          <Input.Search placeholder="Search" value={keyword} onChange={(e) => setKeyword(e.target.value)} />
          {operating ? (
            <PageLoading />
          ) : state.length ? (
            <ul>
              {state.map((it) => (
                <li key={it.name} onClick={() => navigate(PATHS.PROJECT(it.name))}>
                  {it.name}
                </li>
              ))}
            </ul>
          ) : (
            <p style={{ textAlign: 'center' }}>No Results</p>
          )}
          {total > state.length && (
            <Flex style={{ marginTop: 16 }} align="center" justify="center" onClick={handleAppend}>
              <Button loading={operating}>Load More</Button>
            </Flex>
          )}
        </S.Selector>
      )}
    </S.ProjectSelector>
  );
};
