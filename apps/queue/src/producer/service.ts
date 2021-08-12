import { InjectQueue } from '@nestjs/bull';
import { Injectable } from '@nestjs/common';
import { Queue } from 'bull';

@Injectable()
export class ProducerService {
  constructor(@InjectQueue('default') private queue: Queue) {
    this.queue.add('Jira', {s: 1});
    this.queue.add('Jira', {s: 2});
    console.info('init')
  }

  async addJob<T>(name: string, options: T): Promise<string> {
    const job = await this.queue.add(name, options);
    return `${job.id}`;
  }
}
