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

import { DeleteOutlined } from '@ant-design/icons';
import { Table, Input, Select, Button, Tooltip } from 'antd';

export type EditableRow = Record<string, string> & { _key: string };

export type EditableField = {
  key: string;
  title: string;
  type?: 'text' | 'select' | 'multiple';
  width?: number;
  placeholder?: string;
  options?: Array<{ label: string; value: string }>;
  // Multi-value fields are stored in the CSV as a single delimited cell.
  separator?: string;
  editable?: boolean;
  invalid?: (row: EditableRow, rows: EditableRow[]) => string | undefined;
};

interface Props {
  fields: EditableField[];
  rows: EditableRow[];
  onChange: (rows: EditableRow[]) => void;
  removable?: boolean;
  emptyText?: string;
}

export const EditableTable = ({ fields, rows, onChange, removable = true, emptyText }: Props) => {
  const update = (key: string, field: string, value: string) =>
    onChange(rows.map((row) => (row._key === key ? { ...row, [field]: value } : row)));

  const columns = fields.map((field) => ({
    title: field.title,
    dataIndex: field.key,
    width: field.width,
    render: (value: string, row: EditableRow) => {
      if (field.editable === false) {
        return value || <span style={{ color: '#9ca3af' }}>—</span>;
      }

      if (field.type === 'select' || field.type === 'multiple') {
        const separator = field.separator ?? ';';
        const multiple = field.type === 'multiple';
        return (
          <Select
            style={{ width: '100%' }}
            mode={multiple ? 'multiple' : undefined}
            allowClear
            showSearch
            optionFilterProp="label"
            placeholder={field.placeholder}
            options={field.options}
            value={multiple ? (value ? value.split(separator).filter(Boolean) : []) : value || undefined}
            onChange={(next) =>
              update(
                row._key,
                field.key,
                multiple ? ((next as string[]) ?? []).join(separator) : ((next as string) ?? ''),
              )
            }
          />
        );
      }

      const error = field.invalid?.(row, rows);
      return (
        <Tooltip title={error} open={error ? undefined : false}>
          <Input
            status={error ? 'error' : undefined}
            placeholder={field.placeholder}
            value={value}
            onChange={(e) => update(row._key, field.key, e.target.value)}
          />
        </Tooltip>
      );
    },
  }));

  if (removable) {
    columns.push({
      title: '',
      dataIndex: '_actions',
      width: 50,
      render: (_: string, row: EditableRow) => (
        <Button
          type="text"
          icon={<DeleteOutlined />}
          onClick={() => onChange(rows.filter((it) => it._key !== row._key))}
        />
      ),
    });
  }

  return (
    <Table
      rowKey="_key"
      size="small"
      columns={columns}
      dataSource={rows}
      pagination={false}
      locale={{ emptyText }}
      scroll={{ x: true }}
    />
  );
};
