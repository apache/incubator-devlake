import { Entity, PrimaryGeneratedColumn, Column } from 'typeorm';

@Entity()
export default class Issue {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  project: string;
}
