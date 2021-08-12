import { Injectable, Scope } from '@nestjs/common';
import { IExecutable } from 'plugins/core/src/executable.interface';

@Injectable({
  scope: Scope.TRANSIENT,
})
class IssueCollector implements IExecutable<void> {
  async execute(): Promise<void> {
    return;
  }
}

export default IssueCollector;
