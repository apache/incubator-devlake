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
import { Button, Icon, Intent } from '@blueprintjs/core';

import { Loading, DeleteButton } from '@/components';
import { Plugins } from '@/plugins';

import type { UseDataScopeList } from './use-data-scope-list';
import { useDataScopeList } from './use-data-scope-list';
import * as S from './styled';

interface Props extends UseDataScopeList {
  groupByTs: boolean;
  onDelete?: (plugin: Plugins, connectionId: ID, scopeId: ID) => void;
}

export const DataScopeList = ({ groupByTs, onDelete, ...props }: Props) => {
  const { loading, scope, scopeTsMap } = useDataScopeList({ ...props });

  if (!scope.length) {
    return <span>No Data Scope Selected</span>;
  }

  if (loading) {
    return <Loading />;
  }

  return (
    <S.ScopeList>
      {!groupByTs &&
        scope.map((sc) => (
          <S.ScopeItem key={sc.id}>
            <span>{sc.name}</span>
            <DeleteButton onDelete={() => onDelete?.(props.plugin, props.connectionId, sc.id)}>
              <Button small minimal intent={Intent.PRIMARY} icon="cross" />
            </DeleteButton>
          </S.ScopeItem>
        ))}

      {groupByTs &&
        Object.keys(scopeTsMap).map((name) => (
          <S.ScopeItemMap key={name}>
            <div className="name">
              <Icon icon="function" />
              <span>{name}</span>
            </div>
            <ul>
              {scopeTsMap[name].map((sc) => (
                <li key={sc.id}>
                  <span>{sc.name}</span>
                  {onDelete && (
                    <DeleteButton onDelete={() => onDelete(props.plugin, props.connectionId, sc.id)}>
                      <Button small minimal intent={Intent.PRIMARY} icon="cross" />
                    </DeleteButton>
                  )}
                </li>
              ))}
            </ul>
          </S.ScopeItemMap>
        ))}
    </S.ScopeList>
  );
};
