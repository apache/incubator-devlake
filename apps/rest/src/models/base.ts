import { BaseEntity, BeforeInsert, Column, PrimaryColumn } from 'typeorm';
import * as uuid from 'uuid';

export type UniqueID = string;

export class BaseModel extends BaseEntity {
  @PrimaryColumn('char', {
    length: 36,
  })
  id: UniqueID;

  @Column('datetime', {
    default: () => 'CURRENT_TIMESTAMP',
  })
  create_time: Date;

  @Column('datetime', {
    default: () => 'CURRENT_TIMESTAMP',
  })
  update_time: Date;

  // FIXME: using database trigger instead of before insert hook
  @BeforeInsert()
  generateUUID() {
    this.id = uuid.v4();
  }
}
