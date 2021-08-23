import { Inject, Injectable } from '@nestjs/common';
import { ProducerService } from 'apps/queue/src/producer/service';
import { Redis } from 'ioredis';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { v4 } from 'uuid';
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
    const sessionId = v4();
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
    const { taskId, jobId, result } = job;
    if (!taskId) {
      return;
    }
    const task = new Task(taskId, this.redis);
    const jobs = await task.next(jobId, result);
    console.info(jobs)
    for (const job of jobs) {
      await this.producer.addJob(job.name, job.data, { jobId: job.id });
    }
    await task.save();
  }
}
