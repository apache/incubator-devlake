import { Column, Entity } from 'typeorm';
import { BaseModel } from './base';

export const SupportedSourceType = ['jira', 'gitlab'] as const;

export type SourceType = typeof SupportedSourceType[number];

@Entity()
export default class Source extends BaseModel {
  @Column('varchar')
  type: SourceType;

  @Column('json')
  options: Record<string, unknown>;
}
