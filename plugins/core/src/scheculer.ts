import { IExecutable } from './executable.interface';

class Scheduler<T, P> implements IExecutable<T> {
  async execute(options: P): Promise<T> {
    return;
  }
}

export default Scheduler;
