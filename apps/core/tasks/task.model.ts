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
      jobindx = this.dag.findIndex({ id: jobId });
      if (jobindx >= this.dag.length) {
        return [];
      }
      const current = this.dag.get(jobindx);
      if (Array.isArray(current)) {
        const sub = current.find((c) => c.id === jobId);
        let w = await this.redis.setnx(
          `${this.taskId}:${sub.name}`,
          'CHECKING',
        );
        while (w != 1) {
          await new Promise((resolve) => setTimeout(() => resolve(true), 10));
          w = await this.redis.setnx(`${this.taskId}:${sub.name}`, 'CHECKING');
        }
        await this.redis.set(`${this.taskId}:${sub.name}:${sub.id}`, 1);
        const allkeys = await this.redis.keys(`${this.taskId}:${sub.name}:*`);
        for (const key of allkeys) {
          const v = await this.redis.get(key);
          if (parseInt(v) === 0) {
            await this.redis.del(`${this.taskId}:${sub.name}`);
            return [];
          }
        }
        await this.redis.del(`${this.taskId}:${sub.name}`);
      }
      jobindx += 1;
    }
    const job = await this._getJobAtIndex(jobindx);
    if (job) {
      const jobs = [];
      if (datas && Array.isArray(datas)) {
        for (const r of datas) {
          const n = {
            ...job,
            id: v4(),
            data: { ...job.data, ...r },
          };
          jobs.push(n);
          await this.redis.set(`${this.taskId}:${job.name}:${n.id}`, 0);
        }
        this.dag.set(jobindx, jobs);
      } else {
        await this.redis.set(`${this.taskId}:${job.name}:${job.id}`, 0);
        jobs.push(job);
      }
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
