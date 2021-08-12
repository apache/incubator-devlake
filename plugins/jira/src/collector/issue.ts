import { IExecutable } from 'plugins/core/src/executable.interface';

class IssueCollector implements IExecutable<void> {
  async execute(): Promise<void> {
    return;
  }
}

export default IssueCollector;
