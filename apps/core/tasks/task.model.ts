import { Redis } from 'ioredis';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { v4 } from 'uuid';

export type Job = {
  id: string;
  name: string;
  data?: any;
};

export default class Task {
  private dag: DAG;

  constructor(private taskId: string, private redis: Redis) {}

  private async _initiDag(): Promise<void> {
    const serilizedPip = await this.redis.get(this.taskId);
    const pipline = JSON.parse(serilizedPip);
    this.dag = new DAG(pipline);
  }

  private async _getJobAtIndex(index: number): Promise<Job> {
    if (!this.dag) {
      await this._initiDag();
    }
    const job = this.dag.get(index);
    if (!job) {
      return null;
    }
    if (!job.id) {
      const jobId = v4();
      job.id = jobId;
    }
    return job;
  }

  async init(dag: DAG): Promise<void> {
    this.dag = dag;
  }

  async next(jobId?: string): Promise<Job[]> {
    if (jobId) {
      const currentIndex = this.dag.findIndex({ id: jobId });
      if (currentIndex < this.dag.length - 1) {
        const job = await this._getJobAtIndex(currentIndex + 1);
        return [job];
      }
    } else {
      const job = await this._getJobAtIndex(0);
      if (job) {
        return [job];
      }
    }
    return [];
  }

  async save(): Promise<void> {
    if (!this.dag) {
      return;
    }
    const piplines = this.dag.getPipline();
    await this.redis.set(this.taskId, JSON.stringify(piplines));
  }
}
