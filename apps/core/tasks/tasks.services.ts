import { Inject, Injectable } from '@nestjs/common';
import { ProducerService } from 'apps/queue/src/producer/service';
import { randomUUID } from 'crypto';
import { Redis } from 'ioredis';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { EventsService } from '../events/events.service';
import Task from './task.model';

export type JobEvent = {
  jobId: string;
  taskId: string;
  result?: any;
  error?: Error;
};

@Injectable()
export class TasksService {
  constructor(
    @Inject('REDIS_TASK_CLIENT') private redis: Redis,
    private events: EventsService,
    private producer: ProducerService,
  ) {
    this.events.on('job:finished', this.handleJobFinishd.bind(this));
  }

  async startTask(taskDag: DAG): Promise<string> {
    const sessionId = randomUUID();
    const task = new Task(sessionId, this.redis);
    task.init(taskDag);
    const startJobs = await task.next();
    for (const job of startJobs) {
      await this.producer.addJob(
        job.name,
        { ...job.data, taskId: sessionId },
        { jobId: job.id },
      );
    }
    await task.save();
    return sessionId;
  }

  async handleJobFinishd(job: JobEvent): Promise<void> {
    const { taskId, jobId } = job;
    const task = new Task(taskId, this.redis);
    const jobs = await task.next(jobId);
    for (const job of jobs) {
      await this.producer.addJob(job.name, job.data, { jobId: job.id });
    }
    await task.save();
  }
}
