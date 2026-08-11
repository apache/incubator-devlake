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

import { useState, useEffect, useMemo } from 'react';
import { Flex, Button, Alert } from 'antd';

import API from '@/api';
import type { ScopeDuplicateGroup } from '@/api/scope';
import { getPluginConfig } from '@/plugins';
import { operator } from '@/utils';

import { SearchLocal } from './search-local';
import { SearchRemote } from './search-remote';

interface Props {
  mode?: 'single' | 'multiple';
  plugin: string;
  connectionId: ID;
  selectedScope?: any[];
  disabledScope?: any[];
  onChangeSelectedScope?: (scope: any[]) => void;
  footer?: React.ReactNode;
  onCancel?: () => void;
  onSubmit?: (origin: any) => void;
}

const getGithubId = (scope: any): number | undefined => {
  const fromData = scope?.data?.githubId;
  if (typeof fromData === 'number' && fromData > 0) {
    return fromData;
  }
  const fromId = Number(scope?.id);
  if (!Number.isNaN(fromId) && fromId > 0) {
    return fromId;
  }
  return undefined;
};

const buildDuplicateWarning = (duplicates: ScopeDuplicateGroup[]): string => {
  const connectionNames = Array.from(
    new Set(
      duplicates.flatMap((d) => d.connections.map((c) => c.connectionName).filter(Boolean)),
    ),
  );
  const repoLabel =
    duplicates.length === 1
      ? duplicates[0].fullName || duplicates[0].htmlUrl || 'This repository'
      : 'One or more selected repositories';

  const via =
    connectionNames.length === 1
      ? `Connection "${connectionNames[0]}"`
      : connectionNames.length > 1
        ? `Connections ${connectionNames.map((n) => `"${n}"`).join(', ')}`
        : 'another connection';

  return `${repoLabel} is already connected via ${via}. Collecting it here will create duplicate pull requests and issue records, which will inflate all metrics for this repository.`;
};

export const DataScopeRemote = ({
  mode = 'multiple',
  plugin,
  connectionId,
  disabledScope,
  onChangeSelectedScope,
  footer,
  onCancel,
  onSubmit,
  ...props
}: Props) => {
  const [selectedScope, setSelectedScope] = useState<any[]>([]);
  const [operating, setOperating] = useState(false);
  const [duplicates, setDuplicates] = useState<ScopeDuplicateGroup[]>([]);
  const [warningDismissed, setWarningDismissed] = useState(false);

  useEffect(() => {
    setSelectedScope(props.selectedScope ?? []);
  }, [props.selectedScope]);

  const config = useMemo(() => getPluginConfig(plugin).dataScope, [plugin]);

  const githubIdsKey = useMemo(() => {
    if (plugin !== 'github') {
      return '';
    }
    return selectedScope
      .map(getGithubId)
      .filter((id): id is number => id !== undefined)
      .sort((a, b) => a - b)
      .join(',');
  }, [plugin, selectedScope]);

  useEffect(() => {
    if (plugin !== 'github' || !githubIdsKey) {
      setDuplicates([]);
      setWarningDismissed(false);
      return;
    }

    let cancelled = false;
    setWarningDismissed(false);

    API.scope
      .scopeDuplicates(plugin, {
        connectionId,
        githubIds: githubIdsKey,
      })
      .then((res) => {
        if (!cancelled) {
          setDuplicates(res.duplicates ?? []);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setDuplicates([]);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [plugin, connectionId, githubIdsKey]);

  const handleSubmit = async () => {
    const [success, res] = await operator(
      () => API.scope.batch(plugin, connectionId, { data: selectedScope.map((it) => it.data) }),
      {
        setOperating,
        formatMessage: () => 'Add data scope successful.',
      },
    );

    if (success) {
      onSubmit?.(res);
    }
  };

  const showWarning = plugin === 'github' && duplicates.length > 0 && !warningDismissed;

  return (
    <Flex vertical>
      {showWarning && (
        <Alert
          type="warning"
          showIcon
          closable
          onClose={() => setWarningDismissed(true)}
          style={{ marginBottom: 16 }}
          message={buildDuplicateWarning(duplicates)}
        />
      )}
      {config.render ? (
        config.render({
          connectionId,
          disabledItems: disabledScope?.map((it) => ({ id: it.id })),
          selectedItems: selectedScope,
          onChangeSelectedItems: onChangeSelectedScope ?? setSelectedScope,
        })
      ) : config.localSearch ? (
        <SearchLocal
          mode={mode}
          plugin={plugin}
          connectionId={connectionId}
          config={config}
          disabledScope={disabledScope ?? []}
          selectedScope={selectedScope}
          onChange={onChangeSelectedScope ?? setSelectedScope}
        />
      ) : (
        <SearchRemote
          mode={mode}
          plugin={plugin}
          connectionId={connectionId}
          config={config}
          disabledScope={disabledScope ?? []}
          selectedScope={selectedScope}
          onChange={onChangeSelectedScope ?? setSelectedScope}
        />
      )}
      {footer !== undefined ? (
        footer
      ) : (
        <Flex style={{ marginTop: 16 }} justify="flex-end" gap="small">
          <Button disabled={operating} onClick={onCancel}>
            Cancel
          </Button>
          <Button type="primary" loading={operating} disabled={!selectedScope.length} onClick={handleSubmit}>
            Save
          </Button>
        </Flex>
      )}
    </Flex>
  );
};
