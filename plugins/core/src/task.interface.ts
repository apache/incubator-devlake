import IExecutable from './executable.interface';

/**
 * Interface for Collector/Enricher
 */
interface Task extends IExecutable<any> {
  /**
   * name of Task
   */
  name(): string;

  /**
   * start Task process
   */
  execute(...args: any[]): Promise<any>;
}

export default Task;
