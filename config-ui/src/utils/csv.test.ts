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

import { describe, it, expect } from 'vitest';

import { parseCSV, stringifyCSV } from './csv';

describe('parseCSV', () => {
  it('parses the org plugin teams.csv shape', () => {
    const { headers, rows } = parseCSV('Id,Name,Alias,ParentId,SortingIndex\n1,Platform,PLT,,0\n2,Payments,PAY,1,1\n');

    expect(headers).toEqual(['Id', 'Name', 'Alias', 'ParentId', 'SortingIndex']);
    expect(rows).toEqual([
      { Id: '1', Name: 'Platform', Alias: 'PLT', ParentId: '', SortingIndex: '0' },
      { Id: '2', Name: 'Payments', Alias: 'PAY', ParentId: '1', SortingIndex: '1' },
    ]);
  });

  it('keeps quoted delimiters, escaped quotes and embedded newlines in one field', () => {
    const { rows } = parseCSV('Id,Name\n1,"Platform, Core"\n2,"He said ""hi"""\n3,"two\nlines"\n');

    expect(rows.map((row) => row.Name)).toEqual(['Platform, Core', 'He said "hi"', 'two\nlines']);
  });

  it('handles CRLF line endings and a header-only file', () => {
    expect(parseCSV('Id,Name\r\n1,Platform\r\n').rows).toEqual([{ Id: '1', Name: 'Platform' }]);
    expect(parseCSV('Id,Name\n').rows).toEqual([]);
  });

  it('pads rows that are missing trailing fields', () => {
    expect(parseCSV('Id,Name,Email\n1,Platform\n').rows).toEqual([{ Id: '1', Name: 'Platform', Email: '' }]);
  });
});

describe('stringifyCSV', () => {
  it('round-trips a table unchanged', () => {
    const csv = 'Id,Name,Email,TeamIds\nu1,Itai,itai@example.com,1;2\nu2,Dana,dana@example.com,2\n';

    expect(stringifyCSV(parseCSV(csv))).toBe(csv);
  });

  it('quotes only the fields that need it', () => {
    const csv = stringifyCSV({
      headers: ['Id', 'Name'],
      rows: [
        { Id: '1', Name: 'Platform, Core' },
        { Id: '2', Name: 'He said "hi"' },
        { Id: '3', Name: 'plain' },
      ],
    });

    expect(csv).toBe('Id,Name\n1,"Platform, Core"\n2,"He said ""hi"""\n3,plain\n');
  });

  it('writes an empty cell for a column the row does not carry', () => {
    expect(stringifyCSV({ headers: ['Id', 'UserId'], rows: [{ Id: 'acc-1' }] })).toBe('Id,UserId\nacc-1,\n');
  });
});
