import { Test, TestingModule } from '@nestjs/testing';
import { INestMicroservice } from '@nestjs/common';
import { QueueModule } from './../src/queue.module';
import { ConsumerService } from '../src/consumer/service';
import { ProducerService } from '../src/producer/service';
import { TasksService } from 'apps/core/tasks/tasks.services';
import { EventsService } from 'apps/core/events/events.service';

describe('QueueController (e2e)', () => {
  let app: INestMicroservice;
  let module: TestingModule;

  beforeAll(async () => {
    module = await Test.createTestingModule({
      imports: [QueueModule],
    }).compile();

    app = module.createNestMicroservice({});
    await app.init();
  });

  afterAll(async () => {
    await app.close();
  });

  it('Consumer', () => {
    const consumer = app.get(ConsumerService);
    expect(consumer).toBeDefined();
  });

  it('Producer', () => {
    const producer = app.get(ProducerService);
    expect(producer).toBeDefined();
  });

  it('TaskService', () => {
    const tasks = app.get(TasksService);
    expect(tasks).toBeDefined();
  });

  it('EventService', () => {
    const event = app.get(EventsService);
    expect(event).toBeDefined();
  });

  it('Send Task', () => {
    const event = app.get(EventsService);
    expect(event).toBeDefined();
  });
});
