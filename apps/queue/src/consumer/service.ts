import { InjectQueue } from '@nestjs/bull';
import { Injectable, OnModuleInit } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import Bull, { Queue } from 'bull';
import Plugin from '../../../../plugins/core/src/plugin';
import { plugins } from './plugins';

@Injectable()
export class ConsumerService implements OnModuleInit {
  constructor(
    @InjectQueue('default') private queue: Queue,
    private moduleRef: ModuleRef,
  ) {
    this.queue.process('*', this.process.bind(this));
    this.queue.on('failed', this.jobFailed.bind(this));
  }

  async onModuleInit(): Promise<void> {
    console.log(`The plugins(${plugins.length}) start initialize.`);
    for (const plugin of plugins) {
      const executor = this.moduleRef.get<Plugin>(plugin, {
        strict: false,
      });
      // get executor name or set to default
      const executorName = executor.name
        ? executor.name()
        : executor.constructor.name.replace('Plugin', '').toLowerCase();

      console.log(`The ${executorName} has been initialized.`);
      await executor.migrateUp(`plugin_${executorName}_`);
      console.log(`The ${executorName} has migrateUp.`);
    }
  }

  async process(job: Bull.Job): Promise<void> {
    const { name, data } = job;
    const executor = this.moduleRef.get<Plugin>(name, {
      strict: false,
    });
    if (executor) {
      executor.execute(data);
    }
  }

  async jobFailed(job: Bull.Job, error: Error): Promise<void> {
    console.error(`${job.name}-${job.id}`, error);
  }
}
