import { ContextId } from '@nestjs/core';
import { IExecutable } from './executable.interface';

class Scheduler<T, P> implements IExecutable<T> {
  async execute(options: P, contextId?: ContextId): Promise<T> {
    return;
  }
}

export default Scheduler;
