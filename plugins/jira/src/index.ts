import { Injectable, Module } from '@nestjs/common';
import Plugin from 'plugins/core/src';

export type JiraCollector =
  | 'ISSUE'
  | 'CHANGELOG'
  | 'COMMENTS'
  | 'REMOTELINK'
  | 'BOARD';

export type JiraOptions = {
  collectors: JiraCollector[];
};

@Injectable()
class Jira implements Plugin {
  async execute(options: JiraOptions): Promise<void> {
    //TODO: Add jira collector and enrichment
    console.info('Excute Jira');
  }
}

export default Jira;
