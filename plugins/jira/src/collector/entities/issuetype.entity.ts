import { Column } from 'typeorm';
import JiraType from './jiratype.entity';

class IssueType extends JiraType {
  @Column()
  description: string;

  @Column()
  subtask: boolean;
}

export default IssueType;
