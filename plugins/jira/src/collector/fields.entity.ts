import { Column } from 'typeorm';
import IssueType from './issuetype.entity';

class Fields {
  @Column(() => IssueType)
  issuetype: IssueType;
}

export default Fields;
