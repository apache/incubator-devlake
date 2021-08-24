import { Inject, Injectable, Scope, Type } from '@nestjs/common';
import BaseEntity from 'plugins/core/src/base.entity';
import DependencyResolver, { DAG } from 'plugins/core/src/dependency.resolver';
import Plugin from 'plugins/core/src/plugin.interface';
import IssueLeadTimeEntity from './tasks/leadtime/leadtime.entity';

export * from './tasks';

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
    console.info('Jira Plugin Executed', options);
    const dag = await this.resolver.resolve(this.exports()[0]);
    return dag;
  }

  exports(): Type<BaseEntity>[] {
    return [IssueLeadTimeEntity];
  }
}

export default Jira;
