import { Inject, Injectable } from '@nestjs/common';
import { ProducerService } from 'apps/queue/src/producer/service';
import { randomUUID } from 'crypto';
import redis from 'ioredis';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { EventsService } from '../events/events.service';
import Task from './task.model';

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
    private producer: ProducerService,
  ) {
    this.events.on('job:finished', this.handleJobFinishd.bind(this));
  }

  async startTask(taskDag: DAG): Promise<string> {
    const sessionId = randomUUID();
    await this.redis.set(sessionId, JSON.stringify(taskDag));
    return sessionId;
  }

  async handleJobFinishd(job: JobEvent): Promise<void> {
    const { taskId, jobId } = job;
    const dag = JSON.parse(await this.redis.get(taskId));
    const task = new Task(dag);
    const jobs = await task.next(jobId);
    for (const job of jobs) {
      await this.producer.addJob(job.name, job.data, { jobId: job.id });
    }
    await this.redis.set(taskId, task.stringify())
  }
}
