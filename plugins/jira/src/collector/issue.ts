import { IExecutable } from 'plugins/core/src/executable.interface';

class IssueCollector implements IExecutable<void> {
  async execute(): Promise<void> {
    console.info('execute Issue Collector')
    return;
  }
}

export default IssueCollector;
