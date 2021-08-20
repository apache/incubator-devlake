import * as _ from 'lodash';
import { getManager } from 'typeorm';
import { BaseModel } from '../src/models/base';

export function ignoreEntityProps<T extends BaseModel>(
  input: T,
): Omit<T, 'id' | 'create_time' | 'update_time'> {
  return _.omit(input, 'id', 'create_time', 'update_time');
}

export async function truncateTableForTest(
  tableNames: string[],
): Promise<void> {
  await getManager().query(`TRUNCATE TABLE ${tableNames.join(',')}`);
}
