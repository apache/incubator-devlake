import BaseEntity from 'plugins/core/src/base.entity';
import { Column, Entity } from 'typeorm';

@Entity()
export default class ChangelogEntity extends BaseEntity {
  @Column()
  key: string;

  @Column()
  fromString: string;

  @Column()
  toString: string;
}
