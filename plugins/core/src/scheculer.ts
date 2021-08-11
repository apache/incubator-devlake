import { IExecutable } from './executable.interface';

class Scheduler implements IExecutable<any> {
  execute(): Promise<any> {
    return;
  }
}

export default Scheduler;
