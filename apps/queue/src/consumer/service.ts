import { InjectQueue } from '@nestjs/bull';
import { Injectable } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import Bull, { Queue } from 'bull';
import IExecutable from 'plugins/core/src/executable.interface';
import { DAG } from 'plugins/core/src/dependency.resolver';

@Injectable()
export class ConsumerService {
  constructor(
    @InjectQueue('default') private queue: Queue,
    private moduleRef: ModuleRef,
  ) {
    this.queue.process('*', this.process.bind(this));
    this.queue.on('failed', this.jobFailed.bind(this));
  }

  async process(job: Bull.Job): Promise<void> {
    const { name, data } = job;
    const executor = this.moduleRef.get<IExecutable<any>>(name, {
      strict: false,
    });
    if (executor) {
      const result = await executor.execute(data);
      if (result instanceof DAG) {
        //TODO: ADD DAG IN TASK SERVICE
      }
    }
  }

  async jobFailed(job: Bull.Job, error: Error): Promise<void> {
    console.error(`${job.name}-${job.id}`, error);
  }
}
