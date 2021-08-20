import BaseEntity from 'plugins/core/src/base.entity';
import { Column, Entity } from 'typeorm';

@Entity()
export default class IssueEntity extends BaseEntity {
  @Column()
  key: string;

  @Column()
  self: string;
}
