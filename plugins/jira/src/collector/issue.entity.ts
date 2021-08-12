import { Column, Entity, PrimaryGeneratedColumn } from "typeorm";
import Fields from "./fields.entity";

@Entity()
class JiraIssue {
  @PrimaryGeneratedColumn('uuid')
  uuid: string;

  @Column('varchar')
  key: string;

  @Column(() => Fields)
  fields: Fields;
}

export default JiraIssue;
