import { DAG } from "plugins/core/src/dependency.resolver";

export type Job = {
  id: string;
  name: string;
  data?: any;
}

class Task {
  constructor(private dag: DAG) {}

  async next(jobId?: string): Promise<Job[]> {
    if (jobId) {

    }
  }
}