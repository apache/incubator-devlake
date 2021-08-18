import { Injectable } from '@nestjs/common';
import Plugin from 'plugins/core/src/plugin.interface';

export type JiraOptions = {
  source: {
    host: string;
    username: string;
    token: string;
  };
  exports: [];
};

@Injectable()
class Jira implements Plugin {
  name(): string {
    return 'Jira';
  }

  async execute(options: JiraOptions): Promise<void> {
    //TODO: Add jira collector and enrichment
    console.info('Excute Jira', options);
    return;
  }
}

export default Jira;
