import { IExecutable } from './executable.interface';

class Scheduler<T> implements IExecutable<T> {
  async execute(...args: any[]): Promise<T> {
    return;
  }
}

export default Scheduler;
