import { Column, Entity } from 'typeorm';
import { BaseModel, UniqueID } from './base';

@Entity()
export class SourceTask extends BaseModel {
  @Column('char', {
    length: 36,
  })
  source_id: UniqueID;

  @Column('json')
  collector: string[];

  @Column('json')
  enricher: string[];

  @Column('json')
  options: Record<string, unknown>;
}
