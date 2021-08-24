import { Test, TestingModule } from '@nestjs/testing';
import { INestMicroservice } from '@nestjs/common';
import { TasksModule } from '../tasks/tasks.module';
import { TasksService } from '../tasks/tasks.services';
import { BullQueueModule } from 'apps/queue/src/bull/queue.module';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { getQueueToken } from '@nestjs/bull';
import { Queue } from 'bull';

describe('EventsModule (e2e)', () => {
  let app: INestMicroservice;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [BullQueueModule.forRoot(), TasksModule],
    }).compile();

    app = moduleFixture.createNestMicroservice(moduleFixture);
    await app.init();
    const queue = app.get<Queue>(getQueueToken('default'));
    await queue.empty();
  });

  afterAll(async () => {
    await app.close();
  });

  it('Initialized', () => {
    const service = app.get(TasksService);
    expect(service).toBeDefined();
  });

  it('Start Task', async () => {
    const service = app.get(TasksService);
    const queue = app.get<Queue>(getQueueToken('default'));
    const taskId = await service.startTask(new DAG([{ name: 'UnitTest' }]));
    expect(taskId).toBeDefined();
    const jobs = await queue.getJobs(['active', 'waiting']);
    expect(jobs).toHaveLength(1);
  }, 10000);

});
