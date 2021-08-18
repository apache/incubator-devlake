import { Column } from 'typeorm';
import IssueEntity from '../issue/issue.entity';

export default class IssueLeadTimeEntity extends IssueEntity {
  @Column('int')
  leadtime: number;
}
