import { Injectable, Scope } from '@nestjs/common';
import { ContextId } from '@nestjs/core';
import Collector from 'plugins/core/src/collector.decorate';
import CollectorRef from 'plugins/core/src/collectorref';
import Scheduler from 'plugins/core/src/scheculer';
import IssueCollector from './collector/issue';

export enum JiraCollector {
  ISSUE = 'ISSUE',
  CHANGELOG = 'CHANGELOG',
  COMMENTS = 'COMMENTS',
  REMOTELINK = 'REMOTELINK',
  BOARD = 'BOARD',
}

export type JiraOptions = {
  collectors: JiraCollector[];
};

@Injectable({
  scope: Scope.TRANSIENT,
})
@Collector({
  [JiraCollector.ISSUE]: IssueCollector,
})
class Jira extends Scheduler<void, JiraOptions> {
  constructor(private collectorRef: CollectorRef) {
    super();
  }

  async execute(options: JiraOptions, contextId?: ContextId): Promise<void> {
    //TODO: Add jira collector and enrichment
    const collector = await this.collectorRef.resolve(
      JiraCollector.ISSUE,
      'Jira',
      contextId,
    );
    await collector.execute(options);
    return;
  }
}

export default Jira;
