import { Injectable, Scope } from '@nestjs/common';
import { IExecutable } from 'plugins/core/src/executable.interface';

@Injectable({
  scope: Scope.TRANSIENT,
})
class IssueCollector implements IExecutable<void> {
  constructor() {
    console.info('initi issue collector')
  }
  async execute(): Promise<void> {
    console.info('execute Issue Collector')
    return;
  }
}

export default IssueCollector;
