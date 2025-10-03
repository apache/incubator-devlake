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

import { useEffect, useMemo } from 'react';
import { Button, Checkbox, Flex, Form, Input, InputNumber, Segmented, Table, Tooltip, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { DeleteOutlined } from '@ant-design/icons';

interface ScopeData {
  prefix?: string;
  year?: number;
  month?: number | null;
  basePath?: string;
}

interface ScopeItem {
  id: string;
  name: string;
  fullName: string;
  data?: ScopeData;
}

interface Props {
  connectionId: ID;
  disabledItems?: Array<{ id: ID }>;
  selectedItems: ScopeItem[];
  onChangeSelectedItems: (items: ScopeItem[]) => void;
}

const CURRENT_YEAR = new Date().getUTCFullYear();
const MONTHS = Array.from({ length: 12 }, (_, idx) => idx + 1);
const MONTH_LABELS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

const DEFAULT_BASE_PATH = 'user-report/AWSLogs';

const ensureLeadingZero = (value: number) => value.toString().padStart(2, '0');

const normalizeBasePath = (value: string) => value.trim().replace(/^\/+/, '').replace(/\/+$/, '');

const trimTrailingSlashes = (value: string) => value.replace(/\/+$/, '');

const extractScopeMeta = (item: ScopeItem) => {
  const data = item.data ?? {};
  const rawPrefix = data.prefix ?? item.fullName ?? item.id;
  const prefix = typeof rawPrefix === 'string' ? trimTrailingSlashes(rawPrefix) : '';
  const segments = prefix ? prefix.split('/').filter(Boolean) : [];

  let month = data.month ?? null;
  if (month === undefined || month === null) {
    const last = segments[segments.length - 1];
    if (last && /^(0[1-9]|1[0-2])$/.test(last)) {
      month = Number(last);
    } else {
      month = null;
    }
  }

  let year = data.year;
  if (year === undefined || year === null) {
    const idx = month ? segments.length - 2 : segments.length - 1;
    const candidate = idx >= 0 ? segments[idx] : undefined;
    if (candidate && /^\d{4}$/.test(candidate)) {
      year = Number(candidate);
    } else {
      year = undefined;
    }
  }

  let baseSegments: string[];
  if (segments.length === 0) {
    baseSegments = [];
  } else if (month) {
    baseSegments = segments.slice(0, Math.max(segments.length - 2, 0));
  } else {
    baseSegments = segments.slice(0, Math.max(segments.length - 1, 0));
  }

  const basePath = normalizeBasePath(data.basePath ?? (baseSegments.length ? baseSegments.join('/') : ''));

  return {
    basePath,
    year: typeof year === 'number' ? year : null,
    month,
    prefix,
  };
};

const deriveBasePathFromSelection = (items: ScopeItem[]) => {
  for (const item of items) {
    const meta = extractScopeMeta(item);
    if (meta.basePath !== undefined) {
      return meta.basePath;
    }
  }
  return undefined;
};

const buildPrefix = (basePath: string, year: number, month: number | null) => {
  const segments = [] as string[];
  const sanitizedBase = normalizeBasePath(basePath);
  if (sanitizedBase) {
    segments.push(sanitizedBase);
  }
  segments.push(String(year));
  if (month !== null && month !== undefined) {
    segments.push(ensureLeadingZero(month));
  }
  return segments.join('/');
};

const createScopeItem = (basePath: string, year: number, month: number | null): ScopeItem => {
  const sanitizedBase = normalizeBasePath(basePath);
  const prefix = buildPrefix(sanitizedBase, year, month);
  const isFullYear = month === null;
  const name = isFullYear
    ? `${year} (Full Year)`
    : `${year}-${ensureLeadingZero(month as number)} (${MONTH_LABELS[(month as number) - 1]})`;

  return {
    id: prefix,
    name,
    fullName: prefix,
    data: {
      basePath: sanitizedBase,
      prefix,
      year,
      month,
    },
  };
};

const formatScopeLabel = (item: ScopeItem) => {
  const meta = extractScopeMeta(item);
  if (!meta.year) {
    return item.name;
  }

  if (meta.month) {
    const monthLabel = MONTH_LABELS[meta.month - 1] ?? ensureLeadingZero(meta.month);
    return `${meta.year}-${ensureLeadingZero(meta.month)} (${monthLabel})`;
  }

  return `${meta.year} (Full Year)`;
};

const MONTH_OPTIONS = MONTHS.map((value) => ({
  label: `${MONTH_LABELS[value - 1]} (${ensureLeadingZero(value)})`,
  value,
}));

type FormValues = {
  basePath: string;
  year: number;
  mode: 'year' | 'months';
  months?: number[];
};

export const QDevDataScope = ({
  connectionId: _connectionId,
  disabledItems,
  selectedItems,
  onChangeSelectedItems,
}: Props) => {
  const [form] = Form.useForm<FormValues>();

  const disabledIds = useMemo(() => new Set(disabledItems?.map((it) => String(it.id)) ?? []), [disabledItems]);

  const derivedBasePath = useMemo(
    () => deriveBasePathFromSelection(selectedItems) ?? DEFAULT_BASE_PATH,
    [selectedItems],
  );

  useEffect(() => {
    if (!form.isFieldsTouched(['basePath'])) {
      form.setFieldsValue({ basePath: derivedBasePath });
    }
  }, [derivedBasePath, form]);

  useEffect(() => {
    form.setFieldsValue({ mode: 'year', year: form.getFieldValue('year') ?? CURRENT_YEAR });
  }, [form]);

  const handleAdd = async () => {
    const { basePath, year, mode, months = [] } = await form.validateFields();

    const normalizedBase = normalizeBasePath(basePath ?? '');
    const normalizedYear = Number(year);
    if (!normalizedYear || Number.isNaN(normalizedYear)) {
      return;
    }

    const currentIds = new Set(selectedItems.map((item) => item.id));
    const hasFullYear = selectedItems.some((item) => {
      const meta = extractScopeMeta(item);
      return (
        meta.basePath === normalizedBase &&
        meta.year === normalizedYear &&
        (meta.month === null || meta.month === undefined)
      );
    });

    const additions: ScopeItem[] = [];

    if (mode === 'year') {
      if (hasFullYear) {
        return;
      }

      const hasMonths = selectedItems.some((item) => {
        const meta = extractScopeMeta(item);
        return meta.basePath === normalizedBase && meta.year === normalizedYear && meta.month !== null;
      });

      if (hasMonths) {
        return;
      }

      const item = createScopeItem(normalizedBase, normalizedYear, null);
      if (!currentIds.has(item.id) && !disabledIds.has(item.id)) {
        additions.push(item);
      }
    } else {
      if (hasFullYear) {
        return;
      }

      const uniqueMonths = Array.from(new Set(months))
        .map((m) => Number(m))
        .filter((m) => !Number.isNaN(m));
      uniqueMonths.sort((a, b) => a - b);

      uniqueMonths.forEach((month) => {
        if (month < 1 || month > 12) {
          return;
        }

        const item = createScopeItem(normalizedBase, normalizedYear, month);
        if (currentIds.has(item.id) || disabledIds.has(item.id)) {
          return;
        }
        additions.push(item);
      });
    }

    if (!additions.length) {
      return;
    }

    const next = [...selectedItems, ...additions];
    next.sort((a, b) => a.id.localeCompare(b.id));
    onChangeSelectedItems(next);

    if (mode === 'months') {
      form.setFieldsValue({ months: [] });
    }
  };

  const handleRemove = (id: string) => {
    onChangeSelectedItems(selectedItems.filter((item) => item.id !== id));
  };

  const columns: ColumnsType<ScopeItem> = [
    {
      title: 'Time Range',
      dataIndex: 'id',
      key: 'name',
      render: (_: unknown, item) => formatScopeLabel(item),
    },
    {
      title: 'S3 Prefix',
      dataIndex: 'id',
      key: 'prefix',
      render: (_: unknown, item) => {
        const meta = extractScopeMeta(item);
        return <Typography.Text code>{meta.prefix}</Typography.Text>;
      },
    },
    {
      title: 'Base Path',
      dataIndex: 'id',
      key: 'basePath',
      render: (_: unknown, item) => {
        const meta = extractScopeMeta(item);
        return meta.basePath ? (
          <Typography.Text>{meta.basePath}</Typography.Text>
        ) : (
          <Typography.Text type="secondary">(bucket root)</Typography.Text>
        );
      },
    },
    {
      title: '',
      dataIndex: 'id',
      key: 'action',
      width: 80,
      align: 'center',
      render: (id: string) => (
        <Tooltip title={disabledIds.has(id) ? 'Scope is used by existing blueprint' : 'Remove'}>
          <Button
            type="text"
            danger
            icon={<DeleteOutlined />}
            disabled={disabledIds.has(id)}
            onClick={() => handleRemove(id)}
          />
        </Tooltip>
      ),
    },
  ];

  return (
    <Flex vertical gap="middle">
      <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>
        Pick which year and month prefixes DevLake should collect from your Q Developer S3 bucket. Leave empty to
        collect all available data.
      </Typography.Paragraph>

      <Form
        form={form}
        layout="inline"
        initialValues={{
          basePath: derivedBasePath,
          year: CURRENT_YEAR,
          mode: 'year',
          months: [],
        }}
        onFinish={handleAdd}
        style={{ rowGap: 16 }}
      >
        <Form.Item
          label="Base Path"
          name="basePath"
          style={{ flex: 1 }}
          tooltip="Common prefix in S3 between the bucket root and the year directory"
        >
          <Input placeholder="user-report/AWSLogs/.../us-east-1" />
        </Form.Item>

        <Form.Item label="Year" name="year" rules={[{ required: true, message: 'Enter year' }]} style={{ width: 160 }}>
          <InputNumber min={2000} max={2100} style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item name="mode" style={{ width: 180 }}>
          <Segmented
            options={[
              { label: 'Full Year', value: 'year' },
              { label: 'Specific Months', value: 'months' },
            ]}
          />
        </Form.Item>

        <Form.Item noStyle shouldUpdate>
          {({ getFieldValue }) =>
            getFieldValue('mode') === 'months' ? (
              <Form.Item
                name="months"
                rules={[{ required: true, message: 'Select at least one month' }]}
                style={{ minWidth: 260 }}
              >
                <Checkbox.Group
                  options={MONTH_OPTIONS}
                  style={{ display: 'grid', gridTemplateColumns: 'repeat(4, minmax(60px, 1fr))', gap: 8 }}
                />
              </Form.Item>
            ) : null
          }
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Add Scope
          </Button>
        </Form.Item>
      </Form>

      <Table
        size="middle"
        rowKey="id"
        columns={columns}
        dataSource={selectedItems}
        pagination={false}
        locale={{ emptyText: 'No scope selected yet.' }}
      />

      {selectedItems.length > 0 && (
        <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>
          These selections will be stored as S3 prefixes and used during data collection.
        </Typography.Paragraph>
      )}
    </Flex>
  );
};
