import { InjectQueue } from '@nestjs/bull';
import { Injectable } from '@nestjs/common';
import { JobOptions, Queue } from 'bull';

@Injectable()
export class ProducerService {
  constructor(@InjectQueue('default') private queue: Queue) {
    this.addJob('Jira', {});
    console.info('Add Jira Task');
  }

  async addJob<T>(
    name: string,
    options: T,
    jobOption?: JobOptions,
  ): Promise<string> {
    const job = await this.queue.add(name, options, jobOption);
    return `${job.id}`;
  }
}
