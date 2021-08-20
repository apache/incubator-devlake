import { Inject, Injectable } from '@nestjs/common';
import { randomUUID } from 'crypto';
import redis from 'ioredis';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { EventsService } from '../events/events.service';

export type JobEvent = {
  jobId: string;
  taskId: string;
  results?: any;
  error?: Error;
};

@Injectable()
export class TasksService {
  constructor(
    @Inject('REDIS_TASK_CLIENT') private redis: redis.Redis,
    private events: EventsService,
  ) {
    this.events.on('job:finished', this.handleJobFinishd.bind(this));
  }

  async startTask(taskDag: DAG): Promise<string> {
    const sessionId = randomUUID();
    await this.redis.set(sessionId, JSON.stringify(taskDag));
    return sessionId;
  }

  async handleJobFinishd(job: JobEvent): Promise<void> {
    const {taskId, jobId} = job;
    const dag = JSON.parse(await this.redis.get(taskId));

  }
}
