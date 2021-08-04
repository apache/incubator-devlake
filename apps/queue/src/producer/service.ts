import { InjectQueue } from '@nestjs/bull';
import { Injectable } from '@nestjs/common';
import { Queue } from 'bull';

@Injectable()
export class ProducerService {
  constructor(@InjectQueue('default') private queue: Queue) {}

  async addJob<T>(name: string, options: T): Promise<string> {
    const job = await this.queue.add(name, options);
    return `${job.id}`;
  }
}
