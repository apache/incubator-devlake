import { randomUUID } from "crypto";
import { DAG } from "plugins/core/src/dependency.resolver";

export type Job = {
  id: string;
  name: string;
  data?: any;
}

export default class Task {
  constructor(private dag: DAG) {}

  private _getJobAtIndex(index: number): Job {
    const job = this.dag.get(index);
    if (!job) {
      return null;
    }
    if (!job.id) {
      const jobId = randomUUID();
      job.id = jobId;
    }
    return job;
  }

  async next(jobId?: string): Promise<Job[]> {
    if (jobId) {
      const currentIndex = this.dag.findIndex({ id: jobId });
      if (currentIndex < this.dag.length - 1) {
        return [this._getJobAtIndex(currentIndex + 1)];
      }
    } else {
      return [this._getJobAtIndex(0)];
    }
    return [];
  }

  stringify(): string {
    return JSON.stringify(this.dag);
  }
}
