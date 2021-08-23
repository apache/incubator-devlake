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

  async next(jobId?: string, datas?: any): Promise<Job[]> {
    if (!this.dag) {
      await this._initiDag();
    }
    let jobindx = 0;
    if (jobId) {
      const id = jobId.split(':')[0];
      jobindx = this.dag.findIndex({ id });
      if (jobindx >= this.dag.length) {
        return [];
      }
      const current = this.dag.get(jobindx);
      const jobkey = `${this.taskId}:${current.name}`;
      const n = await this.redis.decr(jobkey);
      if (n > 0) {
        return [];
      }
      jobindx += 1;
    }
    const job = await this._getJobAtIndex(jobindx);
    if (job) {
      const nextKey = `${this.taskId}:${job.name}`;
      const jobs = [];
      if (datas && Array.isArray(datas)) {
        for (const i in datas) {
          const r = datas[i];
          const n = {
            ...job,
            id: `${job.id}:${i}`,
            data: { ...job.data, ...r },
          };
          jobs.push(n);
        }
      } else {
        jobs.push(job);
      }
      await this.redis.set(nextKey, jobs.length);
      await this.save();
      return jobs;
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
