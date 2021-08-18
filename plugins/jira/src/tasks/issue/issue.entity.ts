import BaseEntity from 'plugins/core/src/base.entity';
import { Column } from 'typeorm';

export default class IssueEntity extends BaseEntity {
  @Column()
  key: string;

  @Column()
  self: string;
}
