import { Inject, Injectable, Scope } from '@nestjs/common';
import DependencyResolver, { DAG } from 'plugins/core/src/dependency.resolver';
import Plugin from 'plugins/core/src/plugin.interface';
import IssueLeadTimeEntity from './tasks/leadtime/leadtime.entity';

export type JiraOptions = {
  source: {
    host: string;
    username: string;
    token: string;
  };
};

@Injectable({ scope: Scope.TRANSIENT })
class Jira implements Plugin {
  constructor(private resolver: DependencyResolver) {}

  async execute(options: JiraOptions): Promise<DAG> {
    //TODO: Add jira collector and enrichment
    console.info('Excute Jira', options);
    const dag = await this.resolver.resolve(IssueLeadTimeEntity);
    return dag;
  }
}

export default Jira;
