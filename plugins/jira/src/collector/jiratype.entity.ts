import { Column } from 'typeorm';

class JiraType {
  @Column()
  id: number;

  @Column()
  name: string;

  @Column()
  self: string;
}

export default JiraType;
