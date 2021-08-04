import { InjectQueue } from '@nestjs/bull';
import { Injectable } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import Bull, { Queue } from 'bull';
import Plugin from 'plugins/core/src';

@Injectable()
export class ConsumerService {
  constructor(
    @InjectQueue('default') private queue: Queue,
    private moduleRef: ModuleRef,
  ) {
    this.queue.process('*', this.process.bind(this));
  }

  async process(job: Bull.Job): Promise<void> {
    const { name, data } = job;
    const executor = this.moduleRef.get<Plugin>(name);
    if (executor) {
      executor.execute(data);
    }
  }
}
