import { Column, Entity } from 'typeorm';
import { BaseModel } from './base';

@Entity()
export default class Source extends BaseModel {
  @Column('varchar')
  type: string;

  @Column('varchar', {
    nullable: true,
  })
  name?: string;

  @Column('json')
  options: Record<string, unknown>;
}
