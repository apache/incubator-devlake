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

export type CSVTable = {
  // Kept in the order the server sent them: the org plugin's endpoints round-trip
  // whole files, so columns the UI does not edit still have to be written back.
  headers: string[];
  rows: Array<Record<string, string>>;
};

// parseCSV implements the subset of RFC 4180 that gocsv (used by the org plugin)
// emits: double quotes around fields containing a delimiter, quote or newline,
// with embedded quotes doubled.
export const parseCSV = (text: string): CSVTable => {
  const records: string[][] = [];
  let record: string[] = [];
  let field = '';
  let quoted = false;
  let i = 0;

  const endField = () => {
    record.push(field);
    field = '';
  };

  const endRecord = () => {
    endField();
    // A trailing newline must not produce a spurious empty record.
    if (record.length > 1 || record[0] !== '') {
      records.push(record);
    }
    record = [];
  };

  while (i < text.length) {
    const char = text[i];

    if (quoted) {
      if (char === '"') {
        if (text[i + 1] === '"') {
          field += '"';
          i += 2;
          continue;
        }
        quoted = false;
        i += 1;
        continue;
      }
      field += char;
      i += 1;
      continue;
    }

    if (char === '"' && field === '') {
      quoted = true;
      i += 1;
      continue;
    }

    if (char === ',') {
      endField();
      i += 1;
      continue;
    }

    if (char === '\r' && text[i + 1] === '\n') {
      endRecord();
      i += 2;
      continue;
    }

    if (char === '\n' || char === '\r') {
      endRecord();
      i += 1;
      continue;
    }

    field += char;
    i += 1;
  }

  if (field !== '' || record.length) {
    endRecord();
  }

  const [headers = [], ...body] = records;

  return {
    headers,
    rows: body.map((cells) =>
      headers.reduce<Record<string, string>>((acc, header, index) => {
        acc[header] = cells[index] ?? '';
        return acc;
      }, {}),
    ),
  };
};

const escapeField = (value: string) => (/[",\r\n]/.test(value) ? `"${value.replace(/"/g, '""')}"` : value);

export const stringifyCSV = ({ headers, rows }: CSVTable): string =>
  [headers, ...rows.map((row) => headers.map((header) => row[header] ?? ''))]
    .map((cells) => cells.map((cell) => escapeField(cell)).join(','))
    .join('\n') + '\n';
