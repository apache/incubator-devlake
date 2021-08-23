import BaseEntity from 'plugins/core/src/base.entity';
import { Column, Entity } from 'typeorm';

@Entity()
export default class ChangelogEntity extends BaseEntity {
  @Column()
  key: string;

  @Column()
  from: string;

  @Column()
  to: string;
}
