import { Column, Entity } from 'typeorm';
import { BaseModel } from './base';

@Entity()
export default class Source extends BaseModel {
  /**
   * type correspond with plugin
   */
  @Column('varchar')
  type: string;

  /**
   * name source name
   */
  @Column('varchar', {
    nullable: true,
  })
  name?: string;

  /**
   * options plugin options
   */
  @Column('json')
  options: Record<string, unknown>;
}
