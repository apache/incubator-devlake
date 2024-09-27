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

import { useState } from 'react';
import { Flex } from 'antd';

import { ConnectionName, ConnectionFormModal } from '@/plugins';

interface Props {
  plugin: string;
  connectionId: ID;
}

export const ConnectionCheck = ({ plugin, connectionId }: Props) => {
  const [open, setOpen] = useState(false);

  return (
    <Flex align="center" gap="small" style={{ paddingLeft: 16, cursor: 'pointer' }}>
      <span>-</span>
      <ConnectionName plugin={plugin} connectionId={connectionId} onClick={() => setOpen(true)} />
      <ConnectionFormModal plugin={plugin} connectionId={connectionId} open={open} onCancel={() => setOpen(false)} />
    </Flex>
  );
};
