import { InjectQueue } from '@nestjs/bull';
import { Injectable } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import Bull, { Queue } from 'bull';
import IExecutable from 'plugins/core/src/executable.interface';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { EventsService } from 'apps/core/events/events.service';
import { TasksService } from 'apps/core/tasks/tasks.services';

@Injectable()
export class ConsumerService {
  constructor(
    @InjectQueue('default') private queue: Queue,
    private moduleRef: ModuleRef,
    private eventsService: EventsService,
    private tasksService: TasksService,
  ) {
    this.queue.process('*', this.process.bind(this));
    this.queue.on('failed', this.jobFailed.bind(this));
    this.queue.on('completed', this.jobCompleted.bind(this));
  }


  async process(job: Bull.Job): Promise<void> {
    const { name, data } = job;
    const executor = this.moduleRef.get<IExecutable<any>>(name, {
      strict: false,
    });
    if (executor) {
      const result = await executor.execute(data);
      if (result instanceof DAG) {
        await this.tasksService.startTask(result);
      } else {
        this.eventsService.emit('job:completed', {
          jobId: job.id,
          taskId: job.data.taskId,
          result,
        });
      }
    }
  }

  async jobFailed(job: Bull.Job, error: Error): Promise<void> {
    console.error(`${job.name}-${job.id}`, error);
  }

  async jobCompleted(job: Bull.Job): Promise<void> {
    console.info(`${job.name}-${job.id}`);
  }
}
